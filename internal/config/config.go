package config

import "os"

type Config struct {
	DatabaseURL  string
	Port         string
	KafkaBrokers string
	GithubToken  string
}

func Load() Config {
	return Config{
		DatabaseURL:  env("DATABASE_URL_GENERAL", ""),
		Port:         env("PORT", "8085"),
		KafkaBrokers: env("KAFKA_BROKERS", "kafka:9092"),
		GithubToken:  env("GITHUB_TOKEN", ""),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
