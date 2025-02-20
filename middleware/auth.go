// middleware/auth.go
package middleware

import (
	"strings"
	"time"

	"github.com/venkatvghub/code-push-server-go/config"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Claims struct {
	UID  uint64 `json:"uid"`
	Hash string `json:"hash"`
	jwt.StandardClaims
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || (parts[0] != "Bearer" && parts[0] != "Basic") {
			c.JSON(401, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token := parts[1]
		var user models.User
		cfg := config.LoadConfig()

		if parts[0] == "Bearer" && len(token) > 64 { // Assuming long tokens are access tokens
			claims := &Claims{}
			tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWT.TokenSecret), nil
			})
			if err != nil || !tkn.Valid {
				c.JSON(401, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			if err := db.Where("id = ?", claims.UID).First(&user).Error; err != nil {
				c.JSON(401, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			if utils.Md5(user.AckCode) != claims.Hash {
				c.JSON(401, gin.H{"error": "Invalid token hash"})
				c.Abort()
				return
			}
		} else { // Auth token or Basic auth
			var tokenModel models.UserToken
			if err := db.Where("tokens = ? AND expires_at > ?", token, time.Now()).First(&tokenModel).Error; err != nil {
				c.JSON(401, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			if err := db.Where("id = ?", tokenModel.UID).First(&user).Error; err != nil {
				c.JSON(401, gin.H{"error": "User not found"})
				c.Abort()
				return
			}
		}

		c.Set("user", user)
		c.Next()
	}
}
