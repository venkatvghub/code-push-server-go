package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

type AccessKeysController struct {
	DB *gorm.DB
}

func (ctrl *AccessKeysController) CreateAccessKey(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID

	var input struct {
		CreatedBy    string `json:"createdBy" binding:"required"`
		FriendlyName string `json:"friendlyName" binding:"required"`
		TTL          int64  `json:"ttl" binding:"required"`
		Description  string `json:"description"`
		IsSession    bool   `json:"isSession"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if friendlyName already exists
	var existingToken models.UserToken
	if err := ctrl.DB.Where("uid = ? AND name = ?", uid, input.FriendlyName).First(&existingToken).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Friendly name already exists"})
		return
	}

	newAccessKey := utils.RandToken(40)
	log.Printf("newAccessKey: %s\n", newAccessKey)
	token := models.UserToken{
		UID:         uid,
		Name:        input.FriendlyName,
		Tokens:      newAccessKey,
		CreatedBy:   input.CreatedBy,
		Description: input.Description,
		IsSession:   utils.BoolToUint8(input.IsSession),
		ExpiresAt:   gorm.DeletedAt{Time: time.Now().Add(time.Duration(input.TTL) * time.Millisecond)},
	}

	if err := ctrl.DB.Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessKey": gin.H{"name": newAccessKey}})
}

func (ctrl *AccessKeysController) SetupRoutes(r *gin.Engine) {
	accessKeys := r.Group("/accessKeys")
	accessKeys.Use(middleware.AuthMiddleware(ctrl.DB))
	{
		accessKeys.POST("", ctrl.CreateAccessKey)
	}
}
