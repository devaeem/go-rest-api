package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	CORSAllowedOrigins string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	corsAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	if port == "" {
		port = "8080"
	}
	return &Config{Port: port, CORSAllowedOrigins: corsAllowedOrigins}
}
