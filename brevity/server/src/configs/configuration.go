package configs

import "time"

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Email      EmailConfig      `mapstructure:"email"`
	Cloudinary CloudinaryConfig `mapstructure:"cloudinary"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	CORS       CORSConfig       `mapstructure:"cors"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	Storage    StorageConfig    `mapstructure:"storage"`
}

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	Environment  string `mapstructure:"environment"`
	Debug        bool   `mapstructure:"debug"`
	UploadDir    string `mapstructure:"upload_dir"`
	BaseURL      string `mapstructure:"base_url"`
	AnonURLLimit int    `mapstructure:"anon_url_limit"`
	AuthURLLimit int    `mapstructure:"auth_url_limit"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	SQLite SQLiteConfig `mapstructure:"sqlite"`
}

type SQLiteConfig struct {
	Path        string `mapstructure:"path"`
	BusyTimeout int    `mapstructure:"busy_timeout"`
	ForeignKeys bool   `mapstructure:"foreign_keys"`
	JournalMode string `mapstructure:"journal_mode"`
	CacheSize   int    `mapstructure:"cache_size"`
}

type JWTConfig struct {
	AccessTokenSecret  string        `mapstructure:"access_token_secret"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenSecret string        `mapstructure:"refresh_token_secret"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
	ResetTokenSecret   string        `mapstructure:"reset_token_secret"`
	Issuer             string        `mapstructure:"issuer"`
	SecureCookie       bool          `mapstructure:"secure_cookie"`
}

type EmailConfig struct {
	Provider string     `mapstructure:"provider"`
	SMTP     SMTPConfig `mapstructure:"smtp"`
}

type SMTPConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	FromEmail string `mapstructure:"from_email"`
	FromName  string `mapstructure:"from_name"`
	UseTLS    bool   `mapstructure:"use_tls"`
}

type CloudinaryConfig struct {
	CloudName string `mapstructure:"cloud_name"`
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`
	Folder    string `mapstructure:"folder"`
}

type StorageConfig struct {
	MaxAvatarSize int64  `mapstructure:"max_avatar_size"`
	UploadDir     string `mapstructure:"upload_dir"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	FilePath string `mapstructure:"file_path"`
}

type CORSConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	MaxAge       string   `mapstructure:"max_age"`
}

type RateLimitConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Requests int    `mapstructure:"requests"`
	Window   string `mapstructure:"window"`
}
