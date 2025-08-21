package routes

import (
	"telemed/controllers"

	"github.com/gofiber/fiber/v2"
)

var Controller controllers.Controller

func Routes(app *fiber.App) {
	app.Post("/verify-email", Controller.Verifyemail)
	app.Post("/resend-email-otp", Controller.SendEmailOTP)
	app.Post("/otp", Controller.VerifyOtp)
	app.Post("/signup", Controller.Signup)
}
