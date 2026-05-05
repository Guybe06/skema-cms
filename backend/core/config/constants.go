package config

const (
	EnvAppPort = "APP_PORT"
	EnvAppEnv  = "APP_ENV"

	EnvCmsDbHost = "CMS_DB_HOST"
	EnvCmsDbPort = "CMS_DB_PORT"
	EnvCmsDbName = "CMS_DB_NAME"
	EnvCmsDbUser = "CMS_DB_USER"
	EnvCmsDbPass = "CMS_DB_PASSWORD"

	EnvRedisURL = "REDIS_URL"

	EnvJwtSecret = "JWT_SECRET"
	EnvJwtExpiry = "JWT_EXPIRATION"

	EnvEncryptionKey = "ENCRYPTION_KEY"

	EnvMailerFrom   = "MAILER_FROM"
	EnvResendAPIKey = "RESEND_API_KEY"

	EnvStorageDriver     = "STORAGE_DRIVER"
	EnvStorageLocalPath  = "STORAGE_LOCAL_PATH"
	EnvStorageS3Endpoint = "STORAGE_S3_ENDPOINT"
	EnvStorageS3Bucket   = "STORAGE_S3_BUCKET"
	EnvStorageS3Key      = "STORAGE_S3_KEY"
	EnvStorageS3Secret   = "STORAGE_S3_SECRET"

	EnvBackendURL  = "BACKEND_URL"
	EnvFrontendURL = "FRONTEND_URL"

	EnvCORSOrigins       = "CORS_ORIGINS"
	EnvRequestTimeoutSec = "REQUEST_TIMEOUT_SECONDS"
	EnvMaxBodySizeMB     = "MAX_BODY_SIZE_MB"

	DefaultPort              = "3000"
	DefaultEnv               = "development"
	DefaultCmsDbPort         = "5432"
	DefaultJwtExpiry         = "7d"
	DefaultStorageDriver     = "local"
	DefaultStoragePath       = "./uploads"
	DefaultBackendURL        = "http://localhost:3000"
	DefaultFrontendURL       = "http://localhost:3001"
	DefaultRequestTimeoutSec = "30"
	DefaultMaxBodySizeMB     = "10"
)
