package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host          string
	Port          string
	Public        string
	SqliteDB      string
	Migrate       string
	NginxConfig   string
	NginxTemplate string
}

var Envs = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		Host:          getEnv("HOST", "192.168.0.102"),
		Port:          getEnv("PORT", "8080"),
		Public:        getEnv("PUBLIC_IP", "192.168.0.102"),
		SqliteDB:      getEnv("SQLITE_DB_PATH", "files/sqlite/envbox.sqlite3"),
		Migrate:       getEnv("MIGRATE", "false"),
		NginxTemplate: getEnv("NGINX_CONF_TMPL_PATH", "template/nginxcfg/nginx.conf.tmpl"),
		NginxConfig:   getEnv("NGINX_CONF_FILE_PATH", "files/nginxcfg/default.conf"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
