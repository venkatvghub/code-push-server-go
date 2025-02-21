package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/controllers"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/services"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

func setupAuthRoutes(r *gin.Engine, ctrl *controllers.AuthController) {
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

func setupIndexRoutes(r *gin.Engine, ctrl *controllers.IndexController) {
	r.GET("/", ctrl.Index)
	r.GET("/tokens", ctrl.Tokens)
	r.GET("/authenticated", middleware.AuthMiddleware(ctrl.DB), ctrl.Authenticated)
	r.GET("/updateCheck", ctrl.UpdateCheck)
	r.POST("/reportStatus/download", ctrl.ReportStatusDownload)
	r.POST("/reportStatus/deploy", ctrl.ReportStatusDeploy)
}

func setupUsersRoutes(r *gin.Engine, ctrl *controllers.UsersController) {
	users := r.Group("/users")
	{
		users.PATCH("/password", middleware.AuthMiddleware(ctrl.DB), ctrl.ChangePassword)
		// Add other user routes (register, exists, etc.) as needed
	}
}

func setupAccessKeysRoutes(r *gin.Engine, ctrl *controllers.AccessKeysController) {
	accessKeys := r.Group("/accessKeys")
	accessKeys.Use(middleware.AuthMiddleware(ctrl.DB))
	{
		accessKeys.POST("", ctrl.CreateAccessKey)
	}
}

func setupAccountRoutes(r *gin.Engine, ctrl *controllers.AccountController) {
	account := r.Group("/account")
	account.Use(middleware.AuthMiddleware(ctrl.DB))
	{
		account.GET("/accessKeys", ctrl.GetAccessKeys)
	}
}

func setupAppsRoutes(r *gin.Engine, ctrl *controllers.AppsController) {
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

func setupIndexV1Routes(r *gin.Engine, ctrl *controllers.IndexV1Controller) {
	v1 := r.Group("/v0.1/public/codepush")
	{
		v1.GET("/update_check", ctrl.UpdateCheck)
		v1.POST("/report_status/download", ctrl.ReportStatusDownload)
		v1.POST("/report_status/deploy", ctrl.ReportStatusDeploy)
	}
}
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authCtrl := controllers.AuthController{DB: db}
	indexCtrl := controllers.IndexController{DB: db, ClientSvc: services.NewClientService(db)}
	usersCtrl := controllers.UsersController{DB: db}
	accessKeysCtrl := controllers.AccessKeysController{DB: db}
	accountCtrl := controllers.AccountController{DB: db}
	appsCtrl := controllers.AppsController{
		DB:      db,
		AppSvc:  services.NewAppService(db),
		AcctSvc: services.NewAccountService(db),
	}
	indexV1Ctrl := controllers.IndexV1Controller{DB: db, ClientSvc: services.NewClientService(db)}

	//authCtrl.SetupRoutes(r)
	setupAuthRoutes(r, &authCtrl)
	//indexCtrl.SetupRoutes(r)
	setupIndexRoutes(r, &indexCtrl)
	//usersCtrl.SetupRoutes(r)
	setupUsersRoutes(r, &usersCtrl)
	//accessKeysCtrl.SetupRoutes(r)
	setupAccessKeysRoutes(r, &accessKeysCtrl)
	//accountCtrl.SetupRoutes(r)
	setupAccountRoutes(r, &accountCtrl)
	//appsCtrl.SetupRoutes(r)
	setupAppsRoutes(r, &appsCtrl)
	//indexV1Ctrl.SetupRoutes(r)
	setupIndexV1Routes(r, &indexV1Ctrl)
}
