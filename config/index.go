package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load("dev.env")

type DbConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

func Db() DbConfig {
	return DbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

var GatewaySecret = os.Getenv("GATEWAY_SECRET")
var AppPassword = os.Getenv("APP_PASSWORD")
var AppEmail = os.Getenv("APP_EMAIL")
var JwtSecret = os.Getenv("JWT_SECRET")
var PaystackSecretKey = os.Getenv("PAYSTACK_SECRET_KEY")
var PaystackCallbackURL = os.Getenv("PAYSTACK_CALLBACK_URL")
var PaystackBaseURL = os.Getenv("PAYSTACK_BASE_URL")
