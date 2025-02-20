// controllers/users.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

type UsersController struct {
	DB *gorm.DB
}

func (ctrl *UsersController) ChangePassword(c *gin.Context) {
	user, _ := c.Get("user")
	uid := user.(models.User).ID

	var input struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "message": "Invalid input"})
		return
	}

	var userModel models.User
	if err := ctrl.DB.First(&userModel, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERROR", "message": "User not found"})
		return
	}

	if !utils.VerifyPassword(input.OldPassword, userModel.Password) {
		c.JSON(http.StatusOK, gin.H{"status": "ERROR", "message": "Incorrect old password"})
		return
	}

	userModel.Password = utils.HashPassword(input.NewPassword)
	userModel.AckCode = utils.RandToken(5)
	if err := ctrl.DB.Save(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (ctrl *UsersController) SetupRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.PATCH("/password", middleware.AuthMiddleware(ctrl.DB), ctrl.ChangePassword)
		// Add other user routes (register, exists, etc.) as needed
	}
}
