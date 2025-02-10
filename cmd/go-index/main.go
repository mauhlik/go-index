package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mauhlik/go-index/config"
	"github.com/mauhlik/go-index/internal/go-index/controllers"
	"github.com/mauhlik/go-index/internal/go-index/providers"
	"github.com/mauhlik/go-index/internal/go-index/services"
	"github.com/sirupsen/logrus"
)

var (
	ErrProviderNotFound    = errors.New("provider not found")
	ErrUnknownProviderType = errors.New("unknown provider type")
)

func main() {
	configFile := "config.yml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "data",
		FieldMap:          logrus.FieldMap{},
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	logger.SetLevel(logrus.WarnLevel)
	logger.SetReportCaller(true)

	router := gin.Default()

	registerRoutes(router, cfg, logger)

	logger.Infof("Starting server on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(router *gin.Engine, cfg *config.Config, logger *logrus.Logger) {
	for _, repo := range cfg.Repositories {
		provider, err := setupProviderForRepository(cfg, repo, logger)
		if err != nil {
			log.Fatalf("Failed to setup provider for repository %s: %v", repo.Name, err)
		}

		versionService := services.NewService(provider, logger)
		versionController := controllers.NewVersionController(versionService, logger)
		group := router.Group("/api/" + repo.Name)
		{
			group.GET("/:module/:artifact/versions", versionController.GetVersions)
			group.GET("/:module/:artifact/versions/latest", versionController.GetLatestVersion)
		}
	}
}

//nolint:ireturn
func setupProviderForRepository(cfg *config.Config, repo config.RepositoryConfig,
	logger *logrus.Logger) (providers.Provider, error) {
	providerConfig, ok := cfg.Providers[repo.Provider]
	if !ok {
		return nil, fmt.Errorf("%w for repository %s: %s", ErrProviderNotFound, repo.Name, repo.Provider)
	}

	var provider providers.Provider

	var err error

	switch conf := providerConfig.(type) {
	case config.LocalProviderConfig:
		provider = providers.NewLocalProvider(conf.Path)
	case config.S3ProviderConfig:
		provider, err = providers.NewS3Provider(
			conf.Bucket, conf.Endpoint, conf.AccessKey, conf.SecretKey, conf.Region, logger,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize S3 provider for repository %s: %w", repo.Name, err)
		}
	default:
		return nil, fmt.Errorf("%w for repository %s", ErrUnknownProviderType, repo.Name)
	}

	return provider, nil
}
