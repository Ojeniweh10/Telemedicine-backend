package middleware

import (
	"fmt"
	"log"
	"strings"
	"telemed/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing or invalid Authorization header",
			})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := config.JwtSecret
		if secret == "" {
			log.Println("No JWT secret key found in config")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Something went wrong, please try again later",
			})
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Token validation error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// Validate and set usertag in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if usertag, ok := claims["usertag"].(string); ok && usertag != "" {
				c.Locals("usertag", usertag)
			} else {
				log.Println("Usertag missing or invalid in token claims")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Unauthorized: Please log in again",
				})
			}
		}
		return c.Next()
	}
}
