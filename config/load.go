package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	ErrUnsupportedFileExtension = errors.New("unsupported file extension")
	ErrInvalidProviderFormat    = errors.New("invalid provider configuration format")
	ErrProviderTypeRequired     = errors.New("provider type is required")
	ErrUnsupportedProviderType  = errors.New("unsupported provider type")
)

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config

	ext := filepath.Ext(filename)

	switch ext {
	case ".json":
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
	case ".yaml", ".yml":
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&config)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileExtension, ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Parse provider configurations
	for key, value := range config.Providers {
		providerMap, isMap := value.(map[interface{}]interface{})
		if !isMap {
			return nil, fmt.Errorf("%w: %s", ErrInvalidProviderFormat, key)
		}

		providerType, isString := providerMap["type"].(string)
		if !isString {
			return nil, fmt.Errorf("%w: %s", ErrProviderTypeRequired, key)
		}

		var providerConfig interface{}

		switch providerType {
		case "local":
			var localConfig LocalProviderConfig
			err = mapToStruct(providerMap, &localConfig, ext)

			if err != nil {
				return nil, err
			}

			providerConfig = localConfig
		case "s3":
			var s3Config S3ProviderConfig
			err = mapToStruct(providerMap, &s3Config, ext)

			if err != nil {
				return nil, err
			}

			providerConfig = s3Config
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedProviderType, providerType)
		}

		config.Providers[key] = providerConfig
	}

	return &config, nil
}

func mapToStruct(providerMap map[interface{}]interface{}, targetStruct interface{}, ext string) error {
	var (
		data []byte
		err  error
	)

	switch ext {
	case ".json":
		stringMap := make(map[string]interface{})

		for key, value := range providerMap {
			strKey, ok := key.(string)
			if !ok {
				return fmt.Errorf("%w: %v", ErrInvalidProviderFormat, key)
			}

			stringMap[strKey] = value
		}

		data, err = json.Marshal(stringMap)

		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		err = json.Unmarshal(data, targetStruct)

		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	case ".yaml", ".yml":
		data, err = yaml.Marshal(providerMap)

		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}

		err = yaml.Unmarshal(data, targetStruct)

		if err != nil {
			return fmt.Errorf("failed to unmarshal YAML: %w", err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFileExtension, ext)
	}

	return nil
}
