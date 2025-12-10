package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Services ServicesConfig `mapstructure:"services"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

type AppConfig struct {
	Name       string `mapstructure:"name"`
	Debug      bool   `mapstructure:"debug"`
	Env        string `mapstructure:"env"`
	BannerPath string `mapstructure:"banner_path"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type ServicesConfig struct {
	EnableServiceA bool `mapstructure:"enable_service_a"`
	EnableServiceB bool `mapstructure:"enable_service_b"`
	EnableServiceC bool `mapstructure:"enable_service_c"`
}

type AuthConfig struct {
	Type   string `mapstructure:"type"` // e.g., "jwt", "apikey", "none"
	Secret string `mapstructure:"secret"`
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
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("auth.type", "none")
	viper.SetDefault("services.enable_service_a", true)
	viper.SetDefault("services.enable_service_b", true)
	viper.SetDefault("services.enable_service_c", true)

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
