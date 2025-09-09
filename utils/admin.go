package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/smtp"
	"telemed/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	check := true

	if err != nil {
		check = false
	}
	return check
}

func GenerateOTP() (string, error) {
	const digits = "0123456789"
	var length = 6
	otp := make([]byte, 6)
	_, err := rand.Read(otp)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		otp[i] = digits[otp[i]%byte(len(digits))]
	}
	return string(otp), nil
}

func SendEmailOTP(Email, otp string) error {
	// Gmail SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s  and it will expire in 10 mins", otp)
	message := []byte("Subject: " + subject + "\r\n" +
		"To: " + Email + "\r\n" +
		"From: " + senderEmail + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{Email}, message)
	return err
}

func GenerateJWT(usertag string) (string, error) {
	secret := config.JwtSecret
	if secret == "" {
		return "", errors.New("no secret key found")
	}

	claims := jwt.MapClaims{
		"usertag": usertag,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
