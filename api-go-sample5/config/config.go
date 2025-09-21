package config

import (
	"os"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	Grpc     GrpcConfig
	Nats     NatsConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type GrpcConfig struct {
	Port string
}

type NatsConfig struct {
	URL string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "api-go-sample5"),
		},
		Grpc: GrpcConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
		Nats: NatsConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
