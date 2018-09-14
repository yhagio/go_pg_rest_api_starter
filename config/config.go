package config

import (
	"fmt"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
}

func getPostgresConfig() PostgresConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func getMailgunConfig() MailgunConfig {
	return MailgunConfig{
		APIKey:       os.Getenv("MAILGUN_API_KEY"),
		PublicAPIKey: os.Getenv("MAILGUN_PUBLIC_KEY"),
		Domain:       os.Getenv("MAILGUN_DOMAIN"),
	}
}

type Config struct {
	Env        string         `env:"APP_ENV"`
	Pepper     string         `env:"PEPPER"`
	HMACKey    string         `env:"HMAC_KEY"`
	Database   PostgresConfig `json:"database"`
	Mailgun    MailgunConfig  `json:"mailgun"`
	SigningKey string         `env:"signing_key"`
}

type MailgunConfig struct {
	APIKey       string `env:"MAILGUN_API_KEY"`
	PublicAPIKey string `env:"MAILGUN_PUBLIC_KEY"`
	Domain       string `env:"MAILGUN_DOMAIN"`
}

func (c Config) IsProd() bool {
	return c.Env == "production"
}

func GetConfig() Config {
	return Config{
		Env:        os.Getenv("APP_ENV"),
		Pepper:     os.Getenv("PEPPER"),
		HMACKey:    os.Getenv("HMAC_KEY"),
		Database:   getPostgresConfig(),
		Mailgun:    getMailgunConfig(),
		SigningKey: os.Getenv("JWT_SIGN_KEY"),
	}
}
