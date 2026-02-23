package config

import "os"

type Config struct {
	ServerPort string
	DBPath     string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "app.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
