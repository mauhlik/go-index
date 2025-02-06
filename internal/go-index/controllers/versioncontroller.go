package controllers

import (
	"fmt"
	"net/http"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type VersionController struct {
	service services.VersionService
	logger  *logrus.Logger
}

func NewVersionController(service services.VersionService, logger *logrus.Logger) *VersionController {
	return &VersionController{service: service, logger: logger}
}

func (vc *VersionController) GetVersions(ctx *gin.Context) {
	moduleName := ctx.Param("module")
	artifactName := ctx.Param("artifact")

	vc.logger.Infof("Fetching versions for module: %s, artifact: %s", moduleName, artifactName)

	versions, err := vc.service.GetVersions(moduleName, artifactName)
	if err != nil {
		vc.logger.WithError(err).Errorf("Failed to get versions for %s/%s", moduleName, artifactName)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get versions: %v", err),
		})

		return
	}

	ctx.JSON(http.StatusOK, versions)
}

func (vc *VersionController) GetLatestVersion(ctx *gin.Context) {
	moduleName := ctx.Param("module")
	artifactName := ctx.Param("artifact")

	vc.logger.Infof("Fetching latest version for module: %s, artifact: %s", moduleName, artifactName)

	latestVersion, err := vc.service.GetLatestVersion(moduleName, artifactName)
	if err != nil {
		vc.logger.WithError(err).Errorf("Failed to get latest version for %s/%s", moduleName, artifactName)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get latest version: %v", err),
		})

		return
	}

	ctx.JSON(http.StatusOK, latestVersion)
}
