package routes

import (
	"telemed/controllers"
	"telemed/middleware"

	"github.com/gofiber/fiber/v2"
)

var Controller controllers.Controller
var WalletController controllers.WalletController

func Routes(app *fiber.App) {
	//onboarding feature, put in oauth feature once the app has been deployed
	app.Post("/verify-email", Controller.Verifyemail)
	app.Post("/resend-email-otp", Controller.SendEmailOTP)
	app.Post("/otp", Controller.VerifyOtp)
	app.Post("/signup", Controller.Signup)
	app.Post("/login", Controller.Login)
	//dashboard , protected with jwt middleware
	app.Get("/get-doctors", middleware.JWTProtected(), Controller.FetchDoctors) //fetching the doctors so as to book an appointment
	app.Post("/book-appointment", middleware.JWTProtected(), Controller.BookAppointment)
	app.Get("/appointments", middleware.JWTProtected(), Controller.FetchAppointment)
	//next endpoint after fetching appointments is to get on a video call to start the consultation
	app.Post("rate-doctor", middleware.JWTProtected(), Controller.RateDoctor)
	app.Get("/medications", middleware.JWTProtected(), Controller.FetchMedications)
	app.Get("/pharmacies", middleware.JWTProtected(), Controller.FetchPharmacies)
	//cart functionality
	app.Post("/cart/:product-id", middleware.JWTProtected(), Controller.AddToCart)
	app.Patch("/cart/:product-id", middleware.JWTProtected(), Controller.UpdateCart)
	app.Delete("/cart/:product-id", middleware.JWTProtected(), Controller.DeleteFromCart)
	app.Get("/cart", middleware.JWTProtected(), Controller.FetchCart)
	app.Get("/billing-details", middleware.JWTProtected(), Controller.FetchBillingDetails) //user clicks checkout button
	//wallet system (crucial for users to be able to pay for services and medications and top up or withdraw from their balance)
	app.Get("/wallet", middleware.JWTProtected(), WalletController.FetchBalance)
	app.Get("/wallet/banks", middleware.JWTProtected(), WalletController.FetchBanks)
	app.Post("/wallet/create-account", middleware.JWTProtected(), WalletController.CreatePayoutAccount)
	app.Post("/wallet/top-up", middleware.JWTProtected(), WalletController.TopUp)
	app.Post("/wallet/withdraw", middleware.JWTProtected(), WalletController.Withdraw)
	app.Get("/wallet/accounts", middleware.JWTProtected(), WalletController.FetchPayoutAccounts)
	app.Get("/payment/callback", WalletController.PaymentCallback) //paystack will redirect to this endpoint after payment
	app.Post("/paystack/webhook", WalletController.PaystackWebhook)
	//profile management
	app.Get("/profile", middleware.JWTProtected(), Controller.FetchProfile)
	app.Patch("/profile", middleware.JWTProtected(), Controller.UpdateProfile)
	app.Post("/update-password", middleware.JWTProtected(), Controller.SendChangePasswordOTP)
	app.Post("change-pwd/otp", middleware.JWTProtected(), Controller.ChangePassword)
}
