package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/services"
	"gorm.io/gorm"
)

type IndexController struct {
	DB        *gorm.DB
	ClientSvc *services.ClientService
}

func (ctrl *IndexController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "CodePushServer"})
}

func (ctrl *IndexController) Tokens(c *gin.Context) {
	c.HTML(http.StatusOK, "tokens.html", gin.H{"title": "Obtain token"})
}

func (ctrl *IndexController) Authenticated(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}

func (ctrl *IndexController) UpdateCheck(c *gin.Context) {
	deploymentKey := c.Query("deploymentKey")
	appVersion := c.Query("appVersion")
	label := c.Query("label")
	packageHash := c.Query("packageHash")
	clientUniqueID := c.Query("clientUniqueId")

	updateInfo, err := ctrl.ClientSvc.UpdateCheck(deploymentKey, appVersion, label, packageHash, clientUniqueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updateInfo": updateInfo})
}

func (ctrl *IndexController) ReportStatusDownload(c *gin.Context) {
	var input struct {
		ClientUniqueID string `json:"clientUniqueId" binding:"required"`
		Label          string `json:"label" binding:"required"`
		DeploymentKey  string `json:"deploymentKey" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctrl.ClientSvc.ReportStatusDownload(input.DeploymentKey, input.Label, input.ClientUniqueID); err != nil {
		// Log error but return OK as per original behavior
	}
	c.JSON(http.StatusOK, "OK")
}

func (ctrl *IndexController) ReportStatusDeploy(c *gin.Context) {
	var input struct {
		ClientUniqueID string `json:"clientUniqueId" binding:"required"`
		Label          string `json:"label" binding:"required"`
		DeploymentKey  string `json:"deploymentKey" binding:"required"`
		Status         int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctrl.ClientSvc.ReportStatusDeploy(input.DeploymentKey, input.Label, input.ClientUniqueID, input.Status); err != nil {
		// Log error but return OK as per original behavior
	}
	c.JSON(http.StatusOK, "OK")
}

func (ctrl *IndexController) SetupRoutes(r *gin.Engine) {
	r.GET("/", ctrl.Index)
	r.GET("/tokens", ctrl.Tokens)
	r.GET("/authenticated", middleware.AuthMiddleware(ctrl.DB), ctrl.Authenticated)
	r.GET("/updateCheck", ctrl.UpdateCheck)
	r.POST("/reportStatus/download", ctrl.ReportStatusDownload)
	r.POST("/reportStatus/deploy", ctrl.ReportStatusDeploy)
}
