package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/models"
	"gorm.io/gorm"
)

type AccountController struct {
	DB *gorm.DB
}

func (ctrl *AccountController) GetAccessKeys(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID

	var tokens []models.UserToken
	if err := ctrl.DB.Where("uid = ?", uid).Order("id DESC").Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch access keys"})
		return
	}

	result := make([]gin.H, len(tokens))
	for i, token := range tokens {
		result[i] = gin.H{
			"name":         "(hidden)",
			"createdTime":  token.CreatedAt.Time.UnixMilli(),
			"createdBy":    token.CreatedBy,
			"expires":      token.ExpiresAt.Time.UnixMilli(),
			"friendlyName": token.Name,
			"description":  token.Description,
		}
	}

	c.JSON(http.StatusOK, gin.H{"accessKeys": result})
}

func (ctrl *AccountController) SetupRoutes(r *gin.Engine) {
	account := r.Group("/account")
	account.Use(middleware.AuthMiddleware(ctrl.DB))
	{
		account.GET("/accessKeys", ctrl.GetAccessKeys)
	}
}
