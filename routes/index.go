package routes

import (
	"telemed/controllers"
	"telemed/middleware"

	"github.com/gofiber/fiber/v2"
)

var Controller controllers.Controller

func Routes(app *fiber.App) {
	//put in oauth feature once the app has been deployed
	app.Post("/verify-email", Controller.Verifyemail)
	app.Post("/resend-email-otp", Controller.SendEmailOTP)
	app.Post("/otp", Controller.VerifyOtp)
	app.Post("/signup", Controller.Signup)
	app.Post("/login", Controller.Login)
	//dashboard , protected with jwt middleware
	app.Get("/get-doctors", middleware.JWTProtected(), Controller.FetchDoctors) //fetching the doctors so as to book an appointment
	app.Post("/book-appointment", middleware.JWTProtected(), Controller.BookAppointment)
	app.Get("/medications", middleware.JWTProtected(), Controller.FetchMedications)
	app.Get("/pharmacies", middleware.JWTProtected(), Controller.FetchPharmacies)
	app.Post("/cart/:product-id", middleware.JWTProtected(), Controller.AddToCart)
}
