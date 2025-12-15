package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Server     ServerConfig     `mapstructure:"server"`
	Services   ServicesConfig   `mapstructure:"services"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Cron       CronConfig       `mapstructure:"cron"`
}

type MonitoringConfig struct {
	Port           string         `mapstructure:"port"`
	UpdatePeriod   time.Duration  `mapstructure:"update_period"`
	Enabled        bool           `mapstructure:"enabled"`
	UploadDir      string         `mapstructure:"upload_dir"`
	Password       string         `mapstructure:"password"`
	Title          string         `mapstructure:"title"`
	Subtitle       string         `mapstructure:"subtitle"`
	MaxPhotoSizeMB int            `mapstructure:"max_photo_size_mb"`
	MinIO          MinIOConfig    `mapstructure:"minio"`
	External       ExternalConfig `mapstructure:"external"`
	ObfuscateAPI   bool           `mapstructure:"obfuscate_api"`
}

type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
}

type ExternalConfig struct {
	Services []ExternalService `mapstructure:"services"`
}

type ExternalService struct {
	Name string `mapstructure:"name"`
	URL  string `mapstructure:"url"`
}

type CronConfig struct {
	Enabled bool              `mapstructure:"enabled"`
	Jobs    map[string]string `mapstructure:"jobs"`
}

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Debug        bool   `mapstructure:"debug"`
	Env          string `mapstructure:"env"`
	BannerPath   string `mapstructure:"banner_path"`
	StartupDelay int    `mapstructure:"startup_delay"` // seconds to show TUI boot screen (0 to skip)
	QuietStartup bool   `mapstructure:"quiet_startup"` // suppress console logs at startup (TUI only)
	EnableTUI    bool   `mapstructure:"enable_tui"`    // enable fancy TUI mode (false = traditional console)
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// ServicesConfig is a dynamic map of service names to their enabled status.
// Add new services directly in config.yaml without modifying this type.
// Example: services:
//
//	service_a: true
//	service_b: false
type ServicesConfig map[string]bool

// IsEnabled checks if a service is enabled. Returns true by default if not specified.
func (s ServicesConfig) IsEnabled(serviceName string) bool {
	if enabled, exists := s[serviceName]; exists {
		return enabled
	}
	return true // Default to enabled if not specified
}

type AuthConfig struct {
	Type   string `mapstructure:"type"` // e.g., "jwt", "apikey", "none"
	Secret string `mapstructure:"secret"`
}

type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type KafkaConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

type PostgresConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath("./config")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("app.name", "Go-Echo-App")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.banner_path", "banner.txt")
	viper.SetDefault("app.startup_delay", 15)   // 15 seconds default
	viper.SetDefault("app.quiet_startup", true) // clean console by default
	viper.SetDefault("app.enable_tui", true)    // TUI enabled by default
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("auth.type", "none")
	// Services config uses a dynamic map - no hardcoded defaults needed
	// Services default to enabled if not specified (see ServicesConfig.IsEnabled)

	viper.SetDefault("redis.enabled", false)
	viper.SetDefault("kafka.enabled", false)
	viper.SetDefault("postgres.enabled", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found; ignore error if desired or return
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
