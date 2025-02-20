package utils

// utils/security.go

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/venkatvghub/code-push-server-go/config"
	"golang.org/x/crypto/bcrypt"
)

var Config = config.LoadConfig()

func Md5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandToken(length int) string {
	uuidStr := uuid.New().String()
	if length > len(uuidStr) {
		length = len(uuidStr) // Ensure length does not exceed the UUID string length
	}
	return uuidStr[:length]
}

func BoolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
