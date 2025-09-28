package config

import (
	"github.com/spf13/viper"
	"log"
)

// Init initializes the application configuration by reading from application.yaml.
func Init() error {
	viper.SetConfigName("application") // name of config file (without extension)
	viper.SetConfigType("yaml")      // or viper.SetConfigType("YAML")
	viper.AddConfigPath("resource")  // path to look for the config file in

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found; using default settings")
		} else {
			return err
		}
	}
	return nil
}

// GetString retrieves a string value from the configuration.
func GetString(key string, defaultValue ...string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}