package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/config"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/services"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

type AppsController struct {
	DB      *gorm.DB
	AppSvc  *services.AppService
	AcctSvc *services.AccountService
}

func (ctrl *AppsController) AddApp(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID

	var input struct {
		Name                         string `json:"name" binding:"required"`
		OS                           string `json:"os" binding:"required"`
		Platform                     string `json:"platform" binding:"required"`
		ManuallyProvisionDeployments bool   `json:"manuallyProvisionDeployments"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := ctrl.AppSvc.AddApp(uid, input.Name, input.OS, input.Platform)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"app": gin.H{
			"name": input.Name,
			"collaborators": gin.H{
				user.(models.User).Email: gin.H{"permission": "Owner"},
			},
		},
	})
}

func (ctrl *AppsController) DeleteApp(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))

	collaborator, err := ctrl.AcctSvc.OwnerCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Delete(&models.App{}, collaborator.AppID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete app"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ctrl *AppsController) RenameApp(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collaborator, err := ctrl.AcctSvc.OwnerCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	if existingApp, _ := ctrl.AppSvc.FindAppByName(uid, input.Name); existingApp != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": input.Name + " exists"})
		return
	}

	if err := ctrl.DB.Model(&models.App{}).Where("id = ?", collaborator.AppID).Update("name", input.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename app"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ctrl *AppsController) ListCollaborators(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))

	collaborator, err := ctrl.AcctSvc.CollaboratorCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var collaborators []models.Collaborator
	if err := ctrl.DB.Where("app_id = ?", collaborator.AppID).Find(&collaborators).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collaborators"})
		return
	}

	result := make(map[string]gin.H)
	for _, col := range collaborators {
		var userModel models.User
		if err := ctrl.DB.Where("id = ?", col.UID).First(&userModel).Error; err == nil {
			result[userModel.Email] = gin.H{
				"permission":       col.Roles,
				"isCurrentAccount": col.UID == uid,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"collaborators": result})
}

func (ctrl *AppsController) AddCollaborator(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))
	email := strings.TrimSpace(c.Param("email"))

	collaborator, err := ctrl.AcctSvc.OwnerCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	targetUser, err := ctrl.AcctSvc.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&models.Collaborator{
		AppID: collaborator.AppID,
		UID:   targetUser.ID,
		Roles: "Collaborator",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add collaborator"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ctrl *AppsController) AddDeployment(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collaborator, err := ctrl.AcctSvc.CollaboratorCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	deployment, err := ctrl.AppSvc.AddDeployment(collaborator.AppID, input.Name)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deployment": gin.H{
			"key":  deployment.DeploymentKey,
			"name": deployment.Name,
		},
	})
}

func (ctrl *AppsController) ReleasePackage(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))
	deploymentName := strings.TrimSpace(c.Param("deploymentName"))

	collaborator, err := ctrl.AcctSvc.CollaboratorCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	deployment, err := ctrl.AppSvc.FindDeploymentByName(collaborator.AppID, deploymentName)
	if err != nil || deployment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if !strings.HasSuffix(file.Filename, ".zip") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type: must be zip"})
		return
	}

	storage := utils.NewStorage()
	cfg := config.LoadConfig()
	tempFilePath := cfg.Common.TempDir + "/" + file.Filename
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file temporarily"})
		return
	}

	key := utils.RandToken(10) + "_" + file.Filename
	if err := storage.UploadFile(tempFilePath, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to storage"})
		os.Remove(tempFilePath)
		return
	}

	defer os.Remove(tempFilePath)

	_, err = ctrl.AppSvc.ReleasePackage(collaborator.AppID, deployment.ID, tempFilePath, c.PostForm("description"), uid, c.PostForm("isMandatory") == "true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "succeed"})
}

func (ctrl *AppsController) PromotePackage(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))

	var input struct {
		SourceDeploymentName string `json:"sourceDeploymentName" binding:"required"`
		DestDeploymentName   string `json:"destDeploymentName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	sourceDeploymentName := strings.TrimSpace(input.SourceDeploymentName)
	destDeploymentName := strings.TrimSpace(input.DestDeploymentName)

	collaborator, err := ctrl.AcctSvc.CollaboratorCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	sourceDeployment, err := ctrl.AppSvc.FindDeploymentByName(collaborator.AppID, sourceDeploymentName)
	if err != nil || sourceDeployment == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": sourceDeploymentName + " does not exist"})
		return
	}

	destDeployment, err := ctrl.AppSvc.FindDeploymentByName(collaborator.AppID, destDeploymentName)
	if err != nil || destDeployment == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": destDeploymentName + " does not exist"})
		return
	}

	var sourcePkg models.Package
	if err := ctrl.DB.Where("id = ?", sourceDeployment.LastDeploymentVersionID).First(&sourcePkg).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source package not found"})
		return
	}

	newPkg := sourcePkg
	newPkg.ID = 0
	newPkg.DeploymentID = destDeployment.ID
	newPkg.ReleaseMethod = "Promote"
	newPkg.ReleasedBy = uid
	if err := ctrl.DB.Create(&newPkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote package"})
		return
	}

	destDeployment.LastDeploymentVersionID = newPkg.ID
	destDeployment.LabelID++
	if err := ctrl.DB.Save(destDeployment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deployment"})
		return
	}

	cfg := config.LoadConfig()
	go ctrl.AppSvc.CreateDiffPackagesByLastNums(collaborator.AppID, &newPkg, cfg.Common.DiffNums)

	c.JSON(http.StatusOK, gin.H{"package": newPkg})
}

func (ctrl *AppsController) RollbackPackage(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID
	appName := strings.TrimSpace(c.Param("appName"))
	deploymentName := strings.TrimSpace(c.Param("deploymentName"))
	label := strings.TrimSpace(c.Param("label"))

	collaborator, err := ctrl.AcctSvc.CollaboratorCan(uid, appName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	deployment, err := ctrl.AppSvc.FindDeploymentByName(collaborator.AppID, deploymentName)
	if err != nil || deployment == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Deployment not found"})
		return
	}

	var pkg models.Package
	if label != "" {
		if err := ctrl.DB.Where("deployment_id = ? AND label = ?", deployment.ID, label).First(&pkg).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
			return
		}
	} else {
		var history []models.DeploymentHistory
		if err := ctrl.DB.Where("deployment_id = ?", deployment.ID).Order("id DESC").Limit(2).Find(&history).Error; err != nil || len(history) < 2 {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "No previous package to rollback to"})
			return
		}
		if err := ctrl.DB.Where("id = ?", history[1].PackageID).First(&pkg).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Previous package not found"})
			return
		}
	}

	newPkg := pkg
	newPkg.ID = 0
	newPkg.ReleaseMethod = "Rollback"
	newPkg.ReleasedBy = uid
	newPkg.Label = "v" + strconv.Itoa(int(deployment.LabelID+1))
	if err := ctrl.DB.Create(&newPkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback package"})
		return
	}

	deployment.LastDeploymentVersionID = newPkg.ID
	deployment.LabelID++
	if err := ctrl.DB.Save(deployment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deployment"})
		return
	}

	if err := ctrl.DB.Create(&models.DeploymentHistory{
		DeploymentID: deployment.ID,
		PackageID:    newPkg.ID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log history"})
		return
	}

	cfg := config.LoadConfig()
	go ctrl.AppSvc.CreateDiffPackagesByLastNums(collaborator.AppID, &newPkg, cfg.Common.DiffNums)

	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func (ctrl *AppsController) SetupRoutes(r *gin.Engine) {
	apps := r.Group("/apps")
	apps.Use(middleware.AuthMiddleware(ctrl.DB))
	{
		apps.POST("", ctrl.AddApp)
		apps.DELETE("/:appName", ctrl.DeleteApp)
		apps.PATCH("/:appName", ctrl.RenameApp)
		apps.GET("/:appName/collaborators", ctrl.ListCollaborators)
		apps.POST("/:appName/collaborators/:email", ctrl.AddCollaborator)
		apps.POST("/:appName/deployments", ctrl.AddDeployment)
		apps.POST("/:appName/deployments/:deploymentName/release", ctrl.ReleasePackage)
		apps.POST("/:appName/deployments/promote", ctrl.PromotePackage) // Changed route
		apps.POST("/:appName/deployments/:deploymentName/rollback", ctrl.RollbackPackage)
		apps.POST("/:appName/deployments/:deploymentName/rollback/:label", ctrl.RollbackPackage)
	}
}
