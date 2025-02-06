package main

import (
	"log"
	"os"
	"time"

	"github.com/MaUhlik-cen56998/go-index/config"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/controllers"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	for _, repo := range cfg.Repositories {
		providerConfig, ok := cfg.Providers[repo.Provider]
		if !ok {
			log.Fatalf("Provider %s not found for repository %s", repo.Provider, repo.Name)
		}

		var provider providers.Provider
		switch conf := providerConfig.(type) {
		case config.LocalProviderConfig:
			provider = providers.NewLocalProvider(conf.Path)
		case config.S3ProviderConfig:
			provider, err = providers.NewS3Provider(
				conf.Bucket, conf.Endpoint, conf.AccessKey, conf.SecretKey, conf.Region, logger,
			)
			if err != nil {
				log.Fatalf("Failed to initialize S3 provider for repository %s: %v", repo.Name, err)
			}
		default:
			log.Fatalf("Unknown provider type for repository %s", repo.Name)
		}

		versionService := services.NewService(provider, logger)
		versionController := controllers.NewVersionController(versionService, logger)
		group := router.Group("/api/" + repo.Name)
		{
			group.GET("/:module/:artifact/versions", versionController.GetVersions)
			group.GET("/:module/:artifact/versions/latest", versionController.GetLatestVersion)
		}
	}

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
