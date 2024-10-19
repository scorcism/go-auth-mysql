package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	JWT_SECRET   string
	JWT_AUTH_EXP int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:   getEnv("PUBLIC_HOST", "127.0.0.1"),
		Port:         getEnv("PORT", "8080"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", "admin"),
		DBAddress:    fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:       getEnv("DB_NAME", "secrets_management"),
		JWT_SECRET:   getEnv("JWT_SECRET", "apple"),
		JWT_AUTH_EXP: getEnvAsInt("JWT_AUTH_EXP", 3600*24*2),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
