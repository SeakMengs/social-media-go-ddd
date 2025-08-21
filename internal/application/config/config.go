package config

import (
	"fmt"
	"log"
	"slices"

	"github.com/spf13/viper"
)

const (
	DB_DRIVER_PG    = "postgres"
	DB_DRIVER_MYSQL = "mysql"
)

type Config struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found, reading from environment: %v", err)
	}

	config := &Config{
		DBDriver:   viper.GetString("DB_DRIVER"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		AppPort:    viper.GetString("PORT"),
	}

	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	return config
}

func (c *Config) DBDriverValid() error {
	drivers := []string{DB_DRIVER_PG, DB_DRIVER_MYSQL}
	if !slices.Contains(drivers, c.DBDriver) {
		return fmt.Errorf("invalid database driver: %s", c.DBDriver)
	}
	return nil
}

// returns a database connection string based on the driver
func (c *Config) BuildDSN() string {
	switch c.DBDriver {
	case DB_DRIVER_PG:
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	case DB_DRIVER_MYSQL:
		// add parseTime to fix unsupported Scan, storing driver.Value type []uint8 into type *time.Time
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	default:
		log.Fatalf("Unsupported DB driver: %s", c.DBDriver)
		return ""
	}
}
