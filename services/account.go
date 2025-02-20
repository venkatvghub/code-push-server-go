package services

import (
	"errors"
	"time"

	"github.com/venkatvghub/code-push-server-go/models"
	"github.com/venkatvghub/code-push-server-go/utils"
	"gorm.io/gorm"
)

type AccountService struct {
	DB *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{DB: db}
}

func (s *AccountService) CollaboratorCan(uid uint64, appName string) (*models.Collaborator, error) {
	var collaborator models.Collaborator
	err := s.DB.Joins("JOIN apps ON apps.id = collaborators.app_id").
		Where("collaborators.uid = ? AND apps.name = ?", uid, appName).
		First(&collaborator).Error
	if err != nil {
		return nil, errors.New("App " + appName + " not exists or permission denied")
	}
	return &collaborator, nil
}

func (s *AccountService) OwnerCan(uid uint64, appName string) (*models.Collaborator, error) {
	collaborator, err := s.CollaboratorCan(uid, appName)
	if err != nil {
		return nil, err
	}
	if collaborator.Roles != "Owner" {
		return nil, errors.New("permission denied, you are not owner")
	}
	return collaborator, nil
}

func (s *AccountService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New(email + " does not exist")
	}
	return &user, nil
}

func (s *AccountService) GetAllAccessKeysByUID(uid uint64) ([]models.UserToken, error) {
	var tokens []models.UserToken
	if err := s.DB.Where("uid = ?", uid).Order("id DESC").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *AccountService) IsExistAccessKeyName(uid uint64, friendlyName string) (bool, error) {
	var token models.UserToken
	err := s.DB.Where("uid = ? AND name = ?", uid, friendlyName).First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AccountService) CreateAccessKey(uid uint64, newAccessKey, friendlyName, createdBy, description string, ttl int64) (*models.UserToken, error) {
	token := models.UserToken{
		UID:         uid,
		Name:        friendlyName,
		Tokens:      newAccessKey,
		CreatedBy:   createdBy,
		Description: description,
		IsSession:   0,
		ExpiresAt:   gorm.DeletedAt{Time: time.Now().Add(time.Duration(ttl) * time.Millisecond)},
	}
	if err := s.DB.Create(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *AccountService) Login(account, password string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ? OR username = ?", account, account).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !utils.VerifyPassword(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}
	return &user, nil
}

// TODO: Note: Registration-related methods (sendRegisterCode, checkRegisterCode, register) are omitted since Redis is removed
// and we don't need email verification for this English-only version per your requirements.
