package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	CORSAllowedOrigins []string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiresIn       time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port: getEnv("PORT", "8080"),

		CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpiresIn:       getDuration("JWT_EXPIRES_IN", 24*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getDuration อ่านค่าแบบ "24h", "30m" ถ้าเขียนผิดจะหยุดโปรแกรม ดีกว่าใช้ค่าผิดเงียบๆ
func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid %s=%q: use a value like 24h or 30m", key, value)
	}
	return duration
}
