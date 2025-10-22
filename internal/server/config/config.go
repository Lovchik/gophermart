package config

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var appConfig Config

type Config struct {
	Address              string
	DatabaseDNS          string
	PrivateKey           string
	PublicKey            string
	AccuralSystemAddress string
}

func GetConfig() Config {
	return appConfig
}

func InitConfig() {
	config := Config{}
	getEnv("DATABASE_URI", "d", "postgres://postgres:root@localhost:5001/gofermart?sslmode=disable", "file storage path ", &config.DatabaseDNS)
	getEnv("RUN_ADDRESS", "a", ":8080", "Server address", &config.Address)
	getEnv("PRIVATE_KEY", "p", "Ci0tLS0tQkVHSU4gRUMgUFJJVkFURSBLRVktLS0tLQpNSGNDQVFFRUlBaDVxQTNybXFRUXV1MHZiS1YvK3pvdXoveS9JeTJwTHBJY1dVU3lJbVN3b0FvR0NDcUdTTTQ5CkF3RUhvVVFEUWdBRVlENTRWL3ZwKzU0UDlEWGFyWXF4NE1QY20rSEtSSVF6TmFzWVNvUlFIUS82UzZQczh0cE0KY1QrS3ZJSUM4Vy9lOWswVzdDbTcyTTFQOWpVN1NMZi92Zz09Ci0tLS0tRU5EIEVDIFBSSVZBVEUgS0VZLS0tLS0K", "Private Key", &config.PrivateKey)
	getEnv("PUBLIC_KEY", "b", "Ci0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tCk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRVlENTRWL3ZwKzU0UDlEWGFyWXF4NE1QY20rSEsKUklRek5hc1lTb1JRSFEvNlM2UHM4dHBNY1QrS3ZJSUM4Vy9lOWswVzdDbTcyTTFQOWpVN1NMZi92Zz09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=", "Public Key", &config.PublicKey)
	getEnv("ACCRUAL_SYSTEM_ADDRESS", "r", "localhost:8080", "address", &config.AccuralSystemAddress)
	flag.Parse()
	appConfig = config
	log.Info("Server config : ", config)

}

func getEnv(envName, flagName, defaultValue, usage string, config *string) {
	flag.StringVar(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		log.Info("Using environment variable "+envName, "- value "+value)
		*config = value
	}
}

func getEnvInt(envName string, flagName string, defaultValue int64, usage string, config *int64) {
	flag.Int64Var(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			*config = parsed
		}
	}
}

func getEnvBool(envName string, flagName string, defaultValue bool, usage string, config *bool) {
	flag.BoolVar(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			*config = parsed
		}
	}
}
