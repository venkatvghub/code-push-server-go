package services

import (
	"errors"

	"github.com/venkatvghub/code-push-server-go/models"
	"gorm.io/gorm"
)

type ClientService struct {
	DB *gorm.DB
}

func NewClientService(db *gorm.DB) *ClientService {
	return &ClientService{DB: db}
}

func (s *ClientService) UpdateCheck(deploymentKey, appVersion, label, packageHash, clientUniqueID string) (map[string]interface{}, error) {
	var deployment models.Deployment
	if err := s.DB.Where("deployment_key = ?", deploymentKey).First(&deployment).Error; err != nil {
		return nil, errors.New("invalid deployment key")
	}

	var pkg models.Package
	if err := s.DB.Where("deployment_id = ? AND label = ?", deployment.ID, label).First(&pkg).Error; err != nil {
		if label == "" && packageHash == "" {
			err = s.DB.Where("id = ?", deployment.LastDeploymentVersionID).First(&pkg).Error
		}
		if err != nil {
			return map[string]interface{}{
				"isAvailable": false,
			}, nil
		}
	}

	if pkg.PackageHash == packageHash {
		return map[string]interface{}{
			"isAvailable": false,
		}, nil
	}

	return map[string]interface{}{
		"isAvailable": true,
		"downloadUrl": pkg.BlobURL,
		"description": pkg.Description,
		"label":       pkg.Label,
		"packageHash": pkg.PackageHash,
		"packageSize": pkg.Size,
		"isMandatory": pkg.IsMandatory == 1,
		"appVersion":  appVersion,
		"packageId":   pkg.ID,
		"rollout":     pkg.Rollout,
		"isDisabled":  pkg.IsDisabled == 1,
	}, nil
}

func (s *ClientService) ReportStatusDownload(deploymentKey, label, clientUniqueID string) error {
	var deployment models.Deployment
	if err := s.DB.Where("deployment_key = ?", deploymentKey).First(&deployment).Error; err != nil {
		return errors.New("invalid deployment key")
	}

	var pkg models.Package
	if err := s.DB.Where("deployment_id = ? AND label = ?", deployment.ID, label).First(&pkg).Error; err != nil {
		return errors.New("invalid label")
	}

	log := models.LogReportDownload{
		PackageID:      pkg.ID,
		ClientUniqueID: clientUniqueID,
	}
	return s.DB.Create(&log).Error
}

func (s *ClientService) ReportStatusDeploy(deploymentKey, label, clientUniqueID string, status int) error {
	var deployment models.Deployment
	if err := s.DB.Where("deployment_key = ?", deploymentKey).First(&deployment).Error; err != nil {
		return errors.New("invalid deployment key")
	}

	var pkg models.Package
	if err := s.DB.Where("deployment_id = ? AND label = ?", deployment.ID, label).First(&pkg).Error; err != nil {
		return errors.New("invalid label")
	}

	log := models.LogReportDeploy{
		Status:         uint8(status),
		PackageID:      pkg.ID,
		ClientUniqueID: clientUniqueID,
	}
	return s.DB.Create(&log).Error
}
