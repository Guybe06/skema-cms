package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type StorageConfig struct {
	Driver     string
	LocalPath  string
	S3Endpoint string
	S3Bucket   string
	S3Key      string
	S3Secret   string
}

type MailerConfig struct {
	From      string
	ResendKey string
}

type SecurityConfig struct {
	CORSOrigins      string
	RequestTimeoutSec int
	MaxBodySizeMB    int64
}

type Config struct {
	Port          string
	Env           string
	Database      DatabaseConfig
	RedisURL      string
	JwtSecret     string
	JwtExpiry     string
	EncryptionKey string
	Storage       StorageConfig
	Mailer        MailerConfig
	Security      SecurityConfig
	BackendURL    string
	FrontendURL   string
}

/*
 * Load charge la configuration depuis les variables d'environnement.
 *
 * Attend  : un fichier .env optionnel dans le répertoire courant.
 * Retourne: un pointeur vers Config rempli avec toutes les valeurs.
 */

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port: getEnv(EnvAppPort, DefaultPort),
		Env:  getEnv(EnvAppEnv, DefaultEnv),
		Database: DatabaseConfig{
			Host:     getEnv(EnvCmsDbHost, "localhost"),
			Port:     getEnv(EnvCmsDbPort, DefaultCmsDbPort),
			Name:     getEnv(EnvCmsDbName, "skema"),
			User:     getEnv(EnvCmsDbUser, "postgres"),
			Password: getEnv(EnvCmsDbPass, ""),
		},
		RedisURL:      getEnv(EnvRedisURL, ""),
		JwtSecret:     getEnv(EnvJwtSecret, ""),
		JwtExpiry:     getEnv(EnvJwtExpiry, DefaultJwtExpiry),
		EncryptionKey: getEnv(EnvEncryptionKey, ""),
		Storage: StorageConfig{
			Driver:     getEnv(EnvStorageDriver, DefaultStorageDriver),
			LocalPath:  getEnv(EnvStorageLocalPath, DefaultStoragePath),
			S3Endpoint: getEnv(EnvStorageS3Endpoint, ""),
			S3Bucket:   getEnv(EnvStorageS3Bucket, ""),
			S3Key:      getEnv(EnvStorageS3Key, ""),
			S3Secret:   getEnv(EnvStorageS3Secret, ""),
		},
		Mailer: MailerConfig{
			From:      getEnv(EnvMailerFrom, ""),
			ResendKey: getEnv(EnvResendAPIKey, ""),
		},
		Security: SecurityConfig{
			CORSOrigins:       getEnv(EnvCORSOrigins, ""),
			RequestTimeoutSec: getEnvInt(EnvRequestTimeoutSec, DefaultRequestTimeoutSec),
			MaxBodySizeMB:     int64(getEnvInt(EnvMaxBodySizeMB, DefaultMaxBodySizeMB)),
		},
		BackendURL:  getEnv(EnvBackendURL, DefaultBackendURL),
		FrontendURL: getEnv(EnvFrontendURL, DefaultFrontendURL),
	}
}

func getEnvInt(key, fallback string) int {
	raw := getEnv(key, fallback)
	n, err := strconv.Atoi(raw)
	if err != nil {
		n, _ = strconv.Atoi(fallback)
	}
	return n
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
