package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
		return nil, errors.New("unsupported file extension")
	}

	if err != nil {
		return nil, err
	}

	// Parse provider configurations
	for key, value := range config.Providers {
		providerMap, ok := value.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("invalid provider configuration format")
		}

		providerType, ok := providerMap["type"].(string)
		if !ok {
			return nil, errors.New("provider type is required")
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
			return nil, errors.New("unknown provider type")
		}

		config.Providers[key] = providerConfig
	}

	return &config, nil
}

func mapToStruct(m map[interface{}]interface{}, s interface{}, ext string) error {
	var data []byte
	var err error

	switch ext {
	case ".json":
		stringMap := make(map[string]interface{})
		for k, v := range m {
			strKey, ok := k.(string)
			if !ok {
				return errors.New("map key is not a string")
			}
			stringMap[strKey] = v
		}
		data, err = json.Marshal(stringMap)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, s)
	case ".yaml", ".yml":
		data, err = yaml.Marshal(m)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(data, s)
	default:
		return errors.New("unsupported file extension")
	}

	return err
}
