package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBConfig   DatabaseConfig
	JWTSecret  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	config := &Config{
		ServerPort: getEnv("SERVER_PORT", "8888"),
		DBConfig: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "foodapp"),
		},
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-key"),
	}

	return config, nil
}

func GetJwtSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}
	return os.Getenv("JWT_SECRET")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
