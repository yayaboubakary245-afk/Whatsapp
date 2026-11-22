package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config centralise les paramètres réseau, base de données et sécurité du serveur Go (Ndue).
type Config struct {
	HTTPPort           string
	TCPPort            string
	UDPPort            string
	DatabaseURL        string
	JWTSecret          string
	UploadDir          string
	MaxUploadBytes     int64
	AllowedEmailDomain string
}

func Load() *Config {
	_ = godotenv.Load()

	maxMB, err := strconv.ParseInt(getEnv("MAX_UPLOAD_MB", "25"), 10, 64)
	if err != nil || maxMB <= 0 {
		maxMB = 25
	}

	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		TCPPort:            getEnv("TCP_PORT", "5000"),
		UDPPort:            getEnv("UDP_PORT", "5001"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/whatsapp_campus?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "supptic_groupe12_dev_secret"),
		UploadDir:          getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadBytes:     maxMB * 1024 * 1024,
		AllowedEmailDomain: getEnv("ALLOWED_EMAIL_DOMAIN", "ensptic.cm"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
