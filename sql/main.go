package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/venkatvghub/code-push-server-go/config"
	"github.com/venkatvghub/code-push-server-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "db",
		Short: "Database management CLI",
	}

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database schema",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.LoadConfig()

			dsn := fmt.Sprintf(
				"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
				cfg.DB.Host,
				cfg.DB.Port,
				cfg.DB.Username,
				cfg.DB.Database,
				cfg.DB.SSLMode,
				cfg.DB.Password,
			)

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal("Failed to connect to database:", err)
			}

			// Drop existing tables
			if err := db.Migrator().DropTable(
				&models.App{},
				&models.Deployment{},
				&models.Package{},
				&models.DeploymentVersion{},
				&models.PackageDiff{},
				&models.UserToken{},
				&models.User{},
				&models.LogReportDeploy{},
				&models.LogReportDownload{},
			); err != nil {
				log.Fatal("Failed to drop tables:", err)
			}

			// Migrate the schema
			if err := db.AutoMigrate(
				&models.App{},
				&models.Deployment{},
				&models.Package{},
				&models.DeploymentVersion{},
				&models.PackageDiff{},
				&models.UserToken{},
				&models.User{},
				&models.LogReportDeploy{},
				&models.LogReportDownload{},
			); err != nil {
				log.Fatal("Failed to migrate database:", err)
			}

			fmt.Println("Database migrated successfully.")
		},
	}

	var seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with initial data",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.LoadConfig()

			dsn := fmt.Sprintf(
				"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
				cfg.DB.Host,
				cfg.DB.Port,
				cfg.DB.Username,
				cfg.DB.Database,
				cfg.DB.SSLMode,
				cfg.DB.Password,
			)

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal("Failed to connect to database:", err)
			}

			// Add seed data
			if err := seedDatabase(db); err != nil {
				log.Fatal("Failed to seed database:", err)
			}

			fmt.Println("Database seeded successfully.")
		},
	}

	rootCmd.AddCommand(migrateCmd, seedCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func seedDatabase(db *gorm.DB) error {
	// Add your seed data here
	db.Create(&models.App{
		Name: "Test App",
	})
	return nil
}
