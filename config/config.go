package config

import (
	"fmt"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
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
	Port       int            `json:"port"`
	Env        string         `json:"env"`
	Pepper     string         `json:"pepper"`
	HMACKey    string         `json:"hmac_key"`
	Database   PostgresConfig `json:"database"`
	Mailgun    MailgunConfig  `json:"mailgun"`
	SigningKey string         `json:"signing_key"`
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

type OAuthConfig struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	AuthURL  string `json:"auth_url"`
	TokenURL string `json:"token_url"`
}

func (c Config) IsProd() bool {
	return c.Env == "production"
}

func GetConfig() Config {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}

	return Config{
		Port:       port,
		Env:        os.Getenv("APP_ENV"),
		Pepper:     os.Getenv("PEPPER"),
		HMACKey:    os.Getenv("HMAC_KEY"),
		Database:   getPostgresConfig(),
		Mailgun:    getMailgunConfig(),
		SigningKey: os.Getenv("JWT_SIGN_KEY"),
	}
}
