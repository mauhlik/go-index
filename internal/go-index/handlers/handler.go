package handlers

import (
	"fmt"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/services"
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	service services.VersionService
}

func NewHandler(service services.VersionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetVersions(ctx *gofr.Context) (interface{}, error) {
	moduleName := ctx.PathParam("module")
	artifactName := ctx.PathParam("artifact")

	versions, err := h.service.GetVersions(moduleName, artifactName)
	if err != nil {
		ctx.Logger.Errorf("Failed to get versions for %s/%s: %v", moduleName, artifactName, err)

		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	ctx.Logger.Infof("Versions for %s/%s: %v", moduleName, artifactName, versions)

	return versions, nil
}

func (h *Handler) GetLatestVersion(ctx *gofr.Context) (interface{}, error) {
	moduleName := ctx.PathParam("module")
	artifactName := ctx.PathParam("artifact")

	latestVersion, err := h.service.GetLatestVersion(moduleName, artifactName)
	if err != nil {
		ctx.Logger.Errorf("Failed to get latest version for %s/%s: %v", moduleName, artifactName, err)

		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	ctx.Logger.Infof("Latest version for %s/%s: %s", moduleName, artifactName, latestVersion)

	return latestVersion, nil
}
