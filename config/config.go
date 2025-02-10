package config

type LocalProviderConfig struct {
	Type string `json:"type" yaml:"type"`
	Path string `json:"path" yaml:"path"` // Path is required for local provider
}

type S3ProviderConfig struct {
	Type      string `json:"type" yaml:"type"`
	Bucket    string `json:"bucket" yaml:"bucket"`       // Bucket is required for S3 provider
	Endpoint  string `json:"endpoint" yaml:"endpoint"`   // Endpoint is required for S3 provider
	AccessKey string `json:"accessKey" yaml:"accessKey"` // AccessKey is required for S3 provider
	SecretKey string `json:"secretKey" yaml:"secretKey"` // SecretKey is required for S3 provider
	Region    string `json:"region" yaml:"region"`       // Region is required for S3 provider
}

type RepositoryConfig struct {
	Name     string `json:"name" yaml:"name"`
	Provider string `json:"provider" yaml:"provider"`
}

type Config struct {
	Port         string                 `json:"port" yaml:"port"`
	Repositories []RepositoryConfig     `json:"repositories" yaml:"repositories"`
	Providers    map[string]interface{} `json:"providers" yaml:"providers"`
}
