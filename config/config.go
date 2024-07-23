package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port string
	SqliteDB string
}

var Envs = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port: getEnv("PORT", "8080"),
		SqliteDB: getEnv("SQLITE_DB_PATH=", "sqlite/envbox.sqlite3"),
		// SqliteDB: getEnv("SQLITE_DB_PATH=", ":memory:"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}