package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host    string
	Port    string
	DB      DBConfig
	JWT     JWTConfig
	Common  CommonConfig
	Storage StorageConfig
}

type SSLConfig struct {
	Cert string
	Key  string
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
	SSL      SSLConfig
	TimeZone string
}

type JWTConfig struct {
	TokenSecret string
}

type CommonConfig struct {
	AllowRegistration bool
	TryLoginTimes     int
	DiffNums          int
	TempDir           string // Renamed from DataDir and moved here
}

type StorageConfig struct {
	Type  string
	Local LocalConfig
	S3    S3Config
}

type LocalConfig struct {
	StorageDir  string
	DownloadUrl string
	Public      string
}

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
	DownloadUrl     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		// Ignore if .env file not found; fallback to os.Getenv
	}

	return Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "3000"),
		DB: DBConfig{
			Username: getEnv("RDS_USERNAME", "root"),
			Password: getEnv("RDS_PASSWORD", "password"),
			Host:     getEnv("RDS_HOST", "127.0.0.1"),
			Port:     getEnv("RDS_PORT", "3306"),
			Database: getEnv("RDS_DATABASE", "codepush"),
			TimeZone: getEnv("DB_TIMEZONE", "Asia/Calcutta"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			TokenSecret: getEnv("TOKEN_SECRET", "INSERT_RANDOM_TOKEN_KEY"),
		},
		Common: CommonConfig{
			AllowRegistration: getEnvBool("ALLOW_REGISTRATION", false),
			TryLoginTimes:     getEnvInt("TRY_LOGIN_TIMES", 4),
			DiffNums:          getEnvInt("DIFF_NUMS", 3),
			TempDir:           getEnv("TEMP_DIR", "/tmp"), // Added TempDir
		},
		Storage: StorageConfig{
			Type: getEnv("STORAGE_TYPE", "local"),
			Local: LocalConfig{
				StorageDir:  getEnv("LOCAL_STORAGE_DIR", "/tmp/codepush"),
				DownloadUrl: getEnv("LOCAL_DOWNLOAD_URL", "http://127.0.0.1:3000/download"),
				Public:      getEnv("LOCAL_PUBLIC", "/download"),
			},
			S3: S3Config{
				AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
				SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
				Region:          getEnv("AWS_REGION", "us-east-1"),
				BucketName:      getEnv("AWS_BUCKET_NAME", ""),
				DownloadUrl:     getEnv("AWS_DOWNLOAD_URL", ""),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if num, err := strconv.Atoi(value); err == nil {
			return num
		}
	}
	return defaultValue
}

func InitDB(dbConfig *DBConfig) *gorm.DB {
	dsn := "host=" + dbConfig.Host + " user=" + dbConfig.Username + " password=" + dbConfig.Password + " dbname=" + dbConfig.Database + " port=" + dbConfig.Port + " sslmode=" + dbConfig.SSLMode + " TimeZone=" + dbConfig.TimeZone
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}
