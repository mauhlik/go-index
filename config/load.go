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

	if err := decodeConfig(file, ext, &config); err != nil {
		return nil, err
	}

	if err := parseProviderConfigs(&config, ext); err != nil {
		return nil, err
	}

	return &config, nil
}

func decodeConfig(file *os.File, ext string, config *Config) error {
	switch ext {
	case ".json":
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return fmt.Errorf("failed to decode JSON config file: %w", err)
		}
	case ".yaml", ".yml":
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return fmt.Errorf("failed to decode YAML config file: %w", err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFileExtension, ext)
	}

	return nil
}

func parseProviderConfigs(config *Config, ext string) error {
	for key, value := range config.Providers {
		providerMap, err := getProviderMap(value)
		if err != nil {
			return fmt.Errorf("%w: %s", err, key)
		}

		providerType, ok := providerMap["type"].(string)
		if !ok {
			return fmt.Errorf("%w: %s", ErrProviderTypeRequired, key)
		}

		providerConfig, err := getProviderConfig(providerType, providerMap, ext)
		if err != nil {
			return err
		}

		config.Providers[key] = providerConfig
	}

	return nil
}

func getProviderMap(value interface{}) (map[string]interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case map[interface{}]interface{}:
		return convertMap(v), nil
	default:
		return nil, ErrInvalidProviderFormat
	}
}

func getProviderConfig(providerType string, providerMap map[string]interface{}, ext string) (interface{}, error) {
	switch providerType {
	case "local":
		var localConfig LocalProviderConfig
		if err := mapToStruct(providerMap, &localConfig, ext); err != nil {
			return nil, err
		}

		return localConfig, nil
	case "s3":
		var s3Config S3ProviderConfig
		if err := mapToStruct(providerMap, &s3Config, ext); err != nil {
			return nil, err
		}

		return s3Config, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProviderType, providerType)
	}
}

func convertMap(input map[interface{}]interface{}) map[string]interface{} {
	output := make(map[string]interface{})

	for key, value := range input {
		strKey, ok := key.(string)
		if !ok {
			continue
		}

		output[strKey] = value
	}

	return output
}

func mapToStruct(providerMap map[string]interface{}, targetStruct interface{}, ext string) error {
	var (
		data []byte
		err  error
	)

	switch ext {
	case ".json":
		data, err = json.Marshal(providerMap)
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
