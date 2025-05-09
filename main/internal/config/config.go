package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type FirebaseConfig struct {
	CredentialsFile string `mapstructure:"credentials_file"`
}

type StripeConfig struct {
	SecretKey     string `mapstructure:"secret_key"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Firebase FirebaseConfig `mapstructure:"firebase"`
	Stripe   StripeConfig   `mapstructure:"stripeclient"`
}

var Cfg Config

// Init reads config.yml into Cfg
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Error parsing config into struct: %v", err)
	}
	fmt.Printf("Loaded config: %+v\n", Cfg)
}
