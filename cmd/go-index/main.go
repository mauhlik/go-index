package main

import (
	"log"
	"os"

	"github.com/MaUhlik-cen56998/go-index/config"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/handlers"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/services"
	"gofr.dev/pkg/gofr"
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

	app := gofr.New()

	for _, repo := range cfg.Repositories {
		providerConfig, ok := cfg.Providers[repo.Provider]
		if !ok {
			log.Fatalf("Provider %s not found for repository %s", repo.Provider, repo.Name)
		}

		var provider providers.Provider
		switch config := providerConfig.(type) {
		case config.LocalProviderConfig:
			provider = providers.NewLocalProvider(config.Path)
		case config.S3ProviderConfig:
			provider, err = providers.NewS3Provider(config.Bucket, config.Endpoint, config.AccessKey, config.SecretKey, config.Region, app.Logger())
			if err != nil {
				log.Fatalf("Failed to initialize S3 provider for repository %s: %v", repo.Name, err)
			}
		default:
			log.Fatalf("Unknown provider type for repository %s", repo.Name)
		}

		service := services.NewService(provider)
		handler := handlers.NewHandler(service)

		app.GET("/api/"+repo.Name+"/{module}/{artifact}/versions", handler.GetVersions)
		app.GET("/api/"+repo.Name+"/{module}/{artifact}/latest", handler.GetLatestVersion)
	}

	app.Run()
}
