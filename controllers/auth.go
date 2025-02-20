package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input struct {
		Account  string `form:"account" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "message": "Invalid input"})
		return
	}

	var user models.User
	if err := ctrl.DB.Where("email = ? OR username = ?", input.Account, input.Account).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ERROR", "message": "Invalid email or password"})
		return
	}

	if !utils.VerifyPassword(input.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{"status": "ERROR", "message": "Invalid email or password"})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":       user.ID,
		"hash":      utils.Md5(user.AckCode),
		"expiredIn": 7200,
	}).SignedString([]byte(utils.Config.JWT.TokenSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "results": gin.H{"tokens": token}})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "message": "Invalid input"})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := ctrl.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ERROR", "message": input.Email + " already registered"})
		return
	}

	// Create new user
	user := models.User{
		Email:     input.Email,
		Username:  input.Email, // Use email as username for simplicity
		Password:  utils.HashPassword(input.Password),
		Identical: utils.RandToken(9),
		AckCode:   utils.RandToken(5),
		CreatedAt: gorm.DeletedAt{Time: time.Now()},
	}
	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": "Failed to register user"})
		return
	}

	// Immediately return success and redirect to login
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Registration successful, please log in"})
}

func (ctrl *AuthController) SetupRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title":        "CodePushServer",
				"email":        c.Query("email"),
				"showRegister": utils.Config.Common.AllowRegistration,
			})
		})
		auth.GET("/password", func(c *gin.Context) {
			c.HTML(http.StatusOK, "password.html", gin.H{"title": "CodePushServer"})
		})
		auth.GET("/register", func(c *gin.Context) {
			if !utils.Config.Common.AllowRegistration {
				c.Redirect(http.StatusFound, "/auth/login")
				return
			}
			c.HTML(http.StatusOK, "register.html", gin.H{
				"title": "CodePushServer",
				"email": c.Query("email"),
			})
		})
		auth.POST("/login", ctrl.Login)
		auth.POST("/logout", ctrl.Logout)
		auth.POST("/register", ctrl.Register)
	}
}
