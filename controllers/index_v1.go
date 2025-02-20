package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/services"
	"gorm.io/gorm"
)

type IndexV1Controller struct {
	DB        *gorm.DB
	ClientSvc *services.ClientService
}

func (ctrl *IndexV1Controller) UpdateCheck(c *gin.Context) {
	deploymentKey := c.Query("deployment_key")
	appVersion := c.Query("app_version")
	label := c.Query("label")
	packageHash := c.Query("package_hash")
	clientUniqueID := c.Query("client_unique_id")

	updateInfo, err := ctrl.ClientSvc.UpdateCheck(deploymentKey, appVersion, label, packageHash, clientUniqueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"update_info": gin.H{
			"download_url":              updateInfo["downloadUrl"],
			"description":               updateInfo["description"],
			"is_available":              updateInfo["isAvailable"],
			"is_disabled":               updateInfo["isDisabled"],
			"target_binary_range":       updateInfo["appVersion"],
			"label":                     updateInfo["label"],
			"package_hash":              updateInfo["packageHash"],
			"package_size":              updateInfo["packageSize"],
			"should_run_binary_version": false,
			"update_app_version":        false,
			"is_mandatory":              updateInfo["isMandatory"],
		},
	})
}

func (ctrl *IndexV1Controller) ReportStatusDownload(c *gin.Context) {
	var input struct {
		ClientUniqueID string `json:"client_unique_id" binding:"required"`
		Label          string `json:"label" binding:"required"`
		DeploymentKey  string `json:"deployment_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctrl.ClientSvc.ReportStatusDownload(input.DeploymentKey, input.Label, input.ClientUniqueID); err != nil {
		// Log error but return OK
	}
	c.JSON(http.StatusOK, "OK")
}

func (ctrl *IndexV1Controller) ReportStatusDeploy(c *gin.Context) {
	var input struct {
		ClientUniqueID string `json:"client_unique_id" binding:"required"`
		Label          string `json:"label" binding:"required"`
		DeploymentKey  string `json:"deployment_key" binding:"required"`
		Status         int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctrl.ClientSvc.ReportStatusDeploy(input.DeploymentKey, input.Label, input.ClientUniqueID, input.Status); err != nil {
		// Log error but return OK
	}
	c.JSON(http.StatusOK, "OK")
}

func (ctrl *IndexV1Controller) SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/v0.1/public/codepush")
	{
		v1.GET("/update_check", ctrl.UpdateCheck)
		v1.POST("/report_status/download", ctrl.ReportStatusDownload)
		v1.POST("/report_status/deploy", ctrl.ReportStatusDeploy)
	}
}
