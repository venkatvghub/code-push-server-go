// main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/venkatvghub/code-push-server-go/config"
	"github.com/venkatvghub/code-push-server-go/middleware"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/routes"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

var db *gorm.DB

/*func SetupRoutes(r *gin.Engine) {
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

	authCtrl.SetupRoutes(r)
	indexCtrl.SetupRoutes(r)
	usersCtrl.SetupRoutes(r)
	accessKeysCtrl.SetupRoutes(r)
	accountCtrl.SetupRoutes(r)
	appsCtrl.SetupRoutes(r)
	indexV1Ctrl.SetupRoutes(r)

}*/

func setupStaticRoutes(r *gin.Engine) {
	if utils.Config.Storage.Type == "local" {
		r.Static("/download", utils.Config.Storage.Local.StorageDir)
	}
	// Static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")
}

func main() {
	// Initialize configuration
	cfg := config.LoadConfig()

	// Database connection
	db = config.InitDB(&cfg.DB)

	// Auto migrate models
	err := db.AutoMigrate(
		&models.App{}, &models.Collaborator{}, &models.Deployment{}, &models.DeploymentHistory{},
		&models.DeploymentVersion{}, &models.Package{}, &models.PackageDiff{}, &models.PackageMetrics{},
		&models.UserToken{}, &models.User{}, &models.Version{}, &models.LogReportDeploy{}, &models.LogReportDownload{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())

	// Static files
	setupStaticRoutes(r)
	routes.SetupRoutes(r, db)

	// Start server
	err = r.Run(cfg.Host + ":" + cfg.Port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
