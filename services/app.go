package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/venkatvghub/code-push-server-go/config"
	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

type AppService struct {
	DB *gorm.DB
}

func NewAppService(db *gorm.DB) *AppService {
	return &AppService{DB: db}
}

func (s *AppService) CreateDiffPackagesByLastNums(appID uint, pkg *models.Package, diffNums int) error {
	var packages []models.Package
	cfg := config.LoadConfig()
	if err := s.DB.Where("deployment_id = ?", pkg.DeploymentID).
		Order("id DESC").
		Limit(diffNums).
		Find(&packages).Error; err != nil {
		log.Printf("Failed to fetch previous packages: %v", err)
		return err
	}

	storage := utils.NewStorage()
	newPkgPath := pkg.BlobURL[len(storage.GetFileURL("")):] // Remove base URL prefix

	for _, oldPkg := range packages {
		if oldPkg.ID == pkg.ID || oldPkg.PackageHash == pkg.PackageHash {
			continue
		}

		oldPkgPath := oldPkg.BlobURL[len(storage.GetFileURL("")):]
		diffFileName := fmt.Sprintf("%d_%s_%s_diff.zip", pkg.ID, pkg.PackageHash[:8], oldPkg.PackageHash[:8])
		tempDiffPath := filepath.Join(cfg.Common.TempDir, diffFileName)

		err := s.createSimpleDiffZip(newPkgPath, oldPkgPath, tempDiffPath)
		if err != nil {
			log.Printf("Failed to create diff for package %d against %d: %v", pkg.ID, oldPkg.ID, err)
			continue
		}

		if err := storage.UploadFile(tempDiffPath, diffFileName); err != nil {
			log.Printf("Failed to upload diff file: %v", err)
			os.Remove(tempDiffPath)
			continue
		}
		os.Remove(tempDiffPath)

		diffInfo, err := os.Stat(tempDiffPath)
		if err != nil {
			log.Printf("Failed to stat diff file: %v", err)
			continue
		}

		diffURL := storage.GetFileURL(diffFileName)
		diff := models.PackageDiff{
			PackageID:              pkg.ID,
			DiffAgainstPackageHash: oldPkg.PackageHash,
			DiffBlobURL:            diffURL,
			DiffSize:               uint(diffInfo.Size()),
		}
		if err := s.DB.Create(&diff).Error; err != nil {
			log.Printf("Failed to save diff to database: %v", err)
			continue
		}
	}
	return nil
}

func (s *AppService) createSimpleDiffZip(newPath, oldPath, diffPath string) error {
	diffFile, err := os.Create(diffPath)
	if err != nil {
		return err
	}
	defer diffFile.Close()

	writer := zip.NewWriter(diffFile)
	defer writer.Close()

	newFile, err := os.Open(newPath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	newWriter, err := writer.Create("new.zip")
	if err != nil {
		return err
	}
	if _, err := io.Copy(newWriter, newFile); err != nil {
		return err
	}

	oldFile, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer oldFile.Close()
	oldWriter, err := writer.Create("old.zip")
	if err != nil {
		return err
	}
	if _, err := io.Copy(oldWriter, oldFile); err != nil {
		return err
	}

	return nil
}

func (s *AppService) AddApp(uid uint64, name, os, platform string) (*models.App, error) {
	var existingApp models.App
	if err := s.DB.Where("uid = ? AND name = ?", uid, name).First(&existingApp).Error; err == nil {
		return nil, errors.New(name + " exists")
	}

	osVal := map[string]uint8{"ios": 1, "android": 2, "windows": 3}[strings.ToLower(os)]
	platformVal := map[string]uint8{"react-native": 1, "cordova": 2}[strings.ToLower(platform)]
	if osVal == 0 || platformVal == 0 {
		return nil, errors.New("invalid OS or Platform")
	}

	app := models.App{
		Name:     name,
		UID:      uid,
		OS:       osVal,
		Platform: platformVal,
	}
	if err := s.DB.Create(&app).Error; err != nil {
		return nil, err
	}

	collaborator := models.Collaborator{
		AppID: app.ID,
		UID:   uid,
		Roles: "Owner",
	}
	if err := s.DB.Create(&collaborator).Error; err != nil {
		return nil, err
	}

	return &app, nil
}

func (s *AppService) FindAppByName(uid uint64, name string) (*models.App, error) {
	var app models.App
	if err := s.DB.Where("uid = ? AND name = ?", uid, name).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (s *AppService) FindDeploymentByName(appID uint, name string) (*models.Deployment, error) {
	var deployment models.Deployment
	if err := s.DB.Where("app_id = ? AND name = ?", appID, name).First(&deployment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &deployment, nil
}

func (s *AppService) AddDeployment(appID uint, name string) (*models.Deployment, error) {
	var existingDeployment models.Deployment
	if err := s.DB.Where("app_id = ? AND name = ?", appID, name).First(&existingDeployment).Error; err == nil {
		return nil, errors.New("deployment already exists")
	}

	deployment := models.Deployment{
		AppID:         appID,
		Name:          name,
		DeploymentKey: utils.RandToken(40),
	}
	if err := s.DB.Create(&deployment).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (s *AppService) ReleasePackage(appID, deploymentID uint, filePath, description string, uid uint64, isMandatory bool) (*models.Package, error) {
	var deployment models.Deployment
	if err := s.DB.First(&deployment, deploymentID).Error; err != nil {
		return nil, errors.New("deployment not found")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	storage := utils.NewStorage()
	key := utils.RandToken(10) + "_" + filepath.Base(filePath)
	if err := storage.UploadFile(filePath, key); err != nil {
		return nil, err
	}

	label := "v" + fmt.Sprintf("%d", deployment.LabelID+1)
	pkg := models.Package{
		DeploymentID:  deploymentID,
		Description:   description,
		PackageHash:   utils.Md5(filePath),
		BlobURL:       storage.GetFileURL(key),
		Size:          uint(fileInfo.Size()),
		ReleaseMethod: "Upload",
		Label:         label,
		ReleasedBy:    uid,
		IsMandatory:   utils.BoolToUint8(isMandatory),
		Rollout:       100,
	}
	if err := s.DB.Create(&pkg).Error; err != nil {
		return nil, err
	}

	deployment.LabelID++
	deployment.LastDeploymentVersionID = pkg.ID
	if err := s.DB.Save(&deployment).Error; err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := s.CreateDiffPackagesByLastNums(appID, &pkg, utils.Config.Common.DiffNums); err != nil {
			log.Printf("Failed to create diff packages for package %d: %v", pkg.ID, err)
		}
	}()

	return &pkg, nil
}
