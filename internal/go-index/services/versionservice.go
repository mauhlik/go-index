package services

import (
	"fmt"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
)

type VersionService interface {
	GetVersions(moduleName, artifactName string) ([]string, error)
	GetLatestVersion(moduleName, artifactName string) (string, error)
}

type VersionServiceImpl struct {
	provider providers.Provider
	logger   *logrus.Logger
}

func NewService(provider providers.Provider, logger *logrus.Logger) *VersionServiceImpl {
	return &VersionServiceImpl{provider: provider, logger: logger}
}

func (vs *VersionServiceImpl) GetVersions(moduleName, artifactName string) ([]string, error) {
	vs.logger.Infof("Fetching versions for module: %s, artifact: %s", moduleName, artifactName)
	versions, err := vs.provider.GetVersions(moduleName, artifactName)

	if err != nil {
		vs.logger.WithError(err).Errorf("Failed to get versions for %s/%s", moduleName, artifactName)

		return nil, fmt.Errorf("failed to get versions from provider: %w", err)
	}

	return versions, nil
}

func (vs *VersionServiceImpl) GetLatestVersion(moduleName, artifactName string) (string, error) {
	vs.logger.Infof("Fetching latest version for module: %s, artifact: %s", moduleName, artifactName)
	versions, err := vs.provider.GetVersions(moduleName, artifactName)

	if err != nil {
		vs.logger.WithError(err).Errorf("Failed to get versions for %s/%s", moduleName, artifactName)

		return "", fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		vs.logger.Infof("No versions found for %s/%s", moduleName, artifactName)

		return "", nil
	}

	semVersions := make([]semver.Version, len(versions))

	for index, version := range versions {
		semVersion, err := semver.Parse(version)
		if err != nil {
			vs.logger.WithError(err).Errorf("Failed to parse version: %s", version)

			return "", fmt.Errorf("failed to parse version: %w", err)
		}

		semVersions[index] = semVersion
	}

	semver.Sort(semVersions)

	return semVersions[len(semVersions)-1].String(), nil
}
