package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig    `mapstructure:"server"`
	Database      DatabaseConfig  `mapstructure:"database"`
	Memcached     MemcachedConfig `mapstructure:"memcached"`
	Log           LogConfig       `mapstructure:"log"`
	EnableFrontend bool           `mapstructure:"enable_frontend"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // "sqlite" or "postgres"
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	DSN      string `mapstructure:"dsn"` // For SQLite, this is the file path
}

type MemcachedConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	MaxEntries int    `mapstructure:"max_entries"`
	Rolling    bool   `mapstructure:"rolling"`
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we just log if it doesn't exist
		log.Printf("Note: .env file not found or error loading: %v", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()
	
	// Set environment variable prefix (optional, for better organization)
	viper.SetEnvPrefix("")
	
	// Replace dots and dashes with underscores in env var names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	
	// Bind environment variables to config keys
	bindEnvVars()
	
	// Set default for enable_frontend
	viper.SetDefault("enable_frontend", true)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config file: %v", err)
		// Use default values
		return &Config{
			Server: ServerConfig{
				Port: 8080,
			},
			Database: DatabaseConfig{
				Type: "sqlite",
				DSN:  "midgard.db",
			},
			Memcached: MemcachedConfig{
				Host: "localhost",
				Port: 11211,
			},
			Log: LogConfig{
				Level:      "info",
				MaxEntries: 1000,
				Rolling:    true,
			},
			EnableFrontend: true, // Default to true
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &config
}

// bindEnvVars binds environment variables to viper keys
func bindEnvVars() {
	// Server config
	viper.BindEnv("server.port", "PORT")
	
	// Database config
	viper.BindEnv("database.type", "DATABASE_TYPE")
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "DATABASE_DBNAME")
	viper.BindEnv("database.dsn", "DATABASE_DSN")
	
	// Memcached config
	viper.BindEnv("memcached.host", "MEMCACHED_HOST")
	viper.BindEnv("memcached.port", "MEMCACHED_PORT")
	
	// Log config
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.max_entries", "LOG_MAX_ENTRIES")
	viper.BindEnv("log.rolling", "LOG_ROLLING")
	
	// Frontend config
	viper.BindEnv("enable_frontend", "ENABLE_FRONTEND")
	
	// Also support direct environment variable access
	// This allows environment variables to override config file values
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}
	if dbType := os.Getenv("DATABASE_TYPE"); dbType != "" {
		viper.Set("database.type", dbType)
	}
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		viper.Set("database.dsn", dsn)
	}
	if host := os.Getenv("DATABASE_HOST"); host != "" {
		viper.Set("database.host", host)
	}
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		viper.Set("database.port", port)
	}
	if user := os.Getenv("DATABASE_USER"); user != "" {
		viper.Set("database.user", user)
	}
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		viper.Set("database.password", password)
	}
	if dbname := os.Getenv("DATABASE_DBNAME"); dbname != "" {
		viper.Set("database.dbname", dbname)
	}
	if memcachedHost := os.Getenv("MEMCACHED_HOST"); memcachedHost != "" {
		viper.Set("memcached.host", memcachedHost)
	}
	if memcachedPort := os.Getenv("MEMCACHED_PORT"); memcachedPort != "" {
		viper.Set("memcached.port", memcachedPort)
	}
	if enableFrontend := os.Getenv("ENABLE_FRONTEND"); enableFrontend != "" {
		viper.Set("enable_frontend", enableFrontend == "true" || enableFrontend == "1")
	}
}
