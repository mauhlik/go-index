package services

import (
	"fmt"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
	"github.com/blang/semver"
)

type VersionService interface {
	GetVersions(moduleName, artifactName string) ([]string, error)
	GetLatestVersion(moduleName, artifactName string) (string, error)
}

type versionService struct {
	provider providers.Provider
}

func NewService(provider providers.Provider) VersionService {
	return &versionService{provider: provider}
}

func (s *versionService) GetVersions(moduleName, artifactName string) ([]string, error) {
	versions, err := s.provider.GetVersions(moduleName, artifactName)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions from provider: %w", err)
	}

	return versions, nil
}

func (s *versionService) GetLatestVersion(moduleName, artifactName string) (string, error) {
	versions, err := s.provider.GetVersions(moduleName, artifactName)
	if err != nil {
		return "", fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		return "", nil
	}

	semVersions := make([]semver.Version, len(versions))

	for i, version := range versions {
		semVersion, err := semver.Parse(version)
		if err != nil {
			return "", fmt.Errorf("failed to parse version: %w", err)
		}
		semVersions[i] = semVersion
	}

	semver.Sort(semVersions)

	return semVersions[len(semVersions)-1].String(), nil
}
