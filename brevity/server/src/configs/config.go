package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var BrevityApp Config

func LoadConfig(configPath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetEnvPrefix("")
	v.AllowEmptyEnv(true)

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// envrionments variables
	envVariables(v)

	if err := v.Unmarshal(&BrevityApp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if BrevityApp.App.Environment == "development" {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed:", e.Name)
			// envrionments variables
			if err := v.Unmarshal(&BrevityApp); err != nil {
				log.Printf("Error reloading config: %v", err)
			} else {
				log.Println("Config reloaded successfully")
			}
		})
	}

	return &BrevityApp, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "Brevity")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.read_timeout", 10*time.Second)
	v.SetDefault("server.write_timeout", 10*time.Second)
	v.SetDefault("server.shutdown_timeout", 15*time.Second)

	v.SetDefault("database.sqlite.path", "./data/brevity.db")
	v.SetDefault("database.sqlite.busy_timeout", 5000)
	v.SetDefault("database.sqlite.foreign_keys", true)
	v.SetDefault("database.sqlite.journal_mode", "WAL")
	v.SetDefault("database.sqlite.cache_size", -2000)

	v.SetDefault("jwt.access_token_expiry", "15m")
	v.SetDefault("jwt.refresh_token_expiry", "168h")
	v.SetDefault("jwt.reset_token_secret", "default_reset_secret_change_in_production")
	v.SetDefault("jwt.issuer", "brevity-service")
	v.SetDefault("jwt.secure_cookie", false)

	v.SetDefault("logger.level", "debug")
	v.SetDefault("logger.format", "console")
	v.SetDefault("logger.file_path", "./logs/brevity.log")

	v.SetDefault("cors.enabled", true)
	v.SetDefault("cors.allow_origins", []string{"*"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.max_age", "12h")

	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.requests", 100)
	v.SetDefault("rate_limit.window", "1m")
}

func envVariables(v *viper.Viper) {
	keys := []string{
		"app.environment",
		"app.debug",
		"app.upload_dir",
		"app.base_url",
		"server.host",
		"server.port",
		"server.read_timeout",
		"server.write_timeout",
		"server.shutdown_timeout",

		"database.sqlite.path",
		"database.sqlite.busy_timeout",
		"database.sqlite.foreign_keys",
		"database.sqlite.journal_mode",
		"database.sqlite.cache_size",

		"jwt.access_token_secret",
		"jwt.access_token_expiry",
		"jwt.refresh_token_secret",
		"jwt.refresh_token_expiry",
		"jwt.reset_token_secret",
		"jwt.issuer",
		"jwt.secure_cookie",

		"email.provider",
		"email.smtp.host",
		"email.smtp.port",
		"email.smtp.username",
		"email.smtp.password",
		"email.smtp.from_email",
		"email.smtp.from_name",
		"email.smtp.use_tls",

		"cloudinary.cloud_name",
		"cloudinary.api_key",
		"cloudinary.api_secret",
		"cloudinary.folder",

		"storage.max_avatar_size",
		"storage.upload_dir",

		"logger.level",
		"logger.format",
		"logger.file_path",

		"cors.enabled",
		"cors.allow_origins",
		"cors.allow_methods",
		"cors.max_age",

		"rate_limit.enabled",
		"rate_limit.requests",
		"rate_limit.window",
	}

	for _, key := range keys {
		if value := v.GetString(key); value != "" {
			expnd := os.ExpandEnv(value)
			v.Set(key, expnd)
		}
	}
}

func GetConfigPath() string {
	path := []string{
		"configs/app.yaml",
		"../configs/app.yaml",
		filepath.Join("src", "configs", "app.yaml"),
		"./app.yaml",
	}

	for _, path := range path {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "configs/app.yaml"
}
