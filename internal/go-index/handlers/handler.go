package handlers

import (
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/services"
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	service services.VersionService
}

func NewHandler(service services.VersionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetVersions(c *gofr.Context) (interface{}, error) {
	moduleName := c.PathParam("module")
	artifactName := c.PathParam("artifact")

	versions, err := h.service.GetVersions(moduleName, artifactName)
	if err != nil {
		c.Logger.Errorf("Failed to get versions for %s/%s: %v", moduleName, artifactName, err)
		return nil, err
	}
	c.Logger.Infof("Versions for %s/%s: %v", moduleName, artifactName, versions)
	return versions, nil
}

func (h *Handler) GetLatestVersion(c *gofr.Context) (interface{}, error) {
	moduleName := c.PathParam("module")
	artifactName := c.PathParam("artifact")

	latestVersion, err := h.service.GetLatestVersion(moduleName, artifactName)
	if err != nil {
		c.Logger.Errorf("Failed to get latest version for %s/%s: %v", moduleName, artifactName, err)
		return nil, err
	}

	c.Logger.Infof("Latest version for %s/%s: %s", moduleName, artifactName, latestVersion)
	return latestVersion, nil
}
