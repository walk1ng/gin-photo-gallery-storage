package constant

const (
	// JWT constants
	JwtSecret         = "JWT_SECRET"
	Jwt               = "jwt"
	JwtExpMinute      = 30
	PhotoStorageAdmin = "admin"

	// Server constants
	ServerPort = "SERVER_PORT"
	PageSize   = 20

	// DB constants
	DBConnect = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	DBType    = "DB_TYPE"
	DBHost    = "DB_HOST"
	DBPort    = "DB_PORT"
	DBUser    = "DB_USER"
	DBPwd     = "DB_PWD"
	DBName    = "DB_NAME"

	// Auth constants
	CookieMaxAge = 1800
	LoginMaxAge  = 1800
	LoginUser    = "LOGIN_"

	// Redis constants
	RedisHost = "REDIS_HOST"
	RedisPort = "REDIS_PORT"

	// Azure Storage constants
	AzStorageBlobURLEndpointFormat = "https://%s.blob.core.windows.net/%s"
	AzStorageAccountName           = "AZ_STORAGE_ACCOUNT"
	AzStorageAccountKey            = "AZ_STORAGE_ACCOUNT_KEY"
	AzStorageContainerName         = "AZ_STORAGE_CONTAINER"

	// Callback constants
	PhotoURLUpdateChannel = "PHOTO_URL_UPDATE"
	PhotoDeleteChannel    = "PHOTO_DELETE"
	PhotoUpdateIDFormat   = "photo_%d"
)
