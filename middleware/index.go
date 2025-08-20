package middleware

import (
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
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// set claims in context for handlers to use
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("usertag", claims["usertag"])
			c.Locals("role", claims["role"])
		}

		return c.Next()
	}
}
