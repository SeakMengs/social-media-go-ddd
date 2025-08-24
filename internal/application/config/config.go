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

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisCacheConfig struct {
	Host     string
	Port     string
	Password string
	// Redis database number
	DB int
}

type Config struct {
	DB      DBConfig
	AppPort string
	Redis   RedisCacheConfig
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found, reading from environment: %v", err)
	}

	dbConfig := DBConfig{
		Driver:   viper.GetString("DB_DRIVER"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Name:     viper.GetString("DB_NAME"),
	}

	appPort := viper.GetString("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	redisConfig := RedisCacheConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	}

	return &Config{
		DB:      dbConfig,
		AppPort: appPort,
		Redis:   redisConfig,
	}
}

func (d *DBConfig) DriverValid() error {
	drivers := []string{DB_DRIVER_PG, DB_DRIVER_MYSQL}
	if !slices.Contains(drivers, d.Driver) {
		return fmt.Errorf("invalid database driver: %s", d.Driver)
	}
	return nil
}

// returns a database connection string based on the driver
func (d *DBConfig) BuildDSN() string {
	switch d.Driver {
	case DB_DRIVER_PG:
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			d.User, d.Password, d.Host, d.Port, d.Name)
	case DB_DRIVER_MYSQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			d.User, d.Password, d.Host, d.Port, d.Name)
	default:
		log.Fatalf("Unsupported DB driver: %s", d.Driver)
		return ""
	}
}

func (r *RedisCacheConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
