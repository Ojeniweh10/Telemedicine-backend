package routes

import (
	"telemed/controllers"
	"telemed/middleware"
	"telemed/responses"

	"github.com/gofiber/fiber/v2"
)

var adminController controllers.AdminController

const (
	Admin   = "admin"
	God_eye = "god_eye"
)

func AdminRoutes(app *fiber.App) {
	api := app.Group("/admin")
	api.Post("/Login", roleMiddleware(Admin, God_eye), adminController.Login)
	api.Post("/otp", roleMiddleware(Admin, God_eye), adminController.VerifyOTP)
	api.Post("/forgot-password", roleMiddleware(Admin, God_eye), adminController.ForgotPassword)
	api.Post("/verify-forgot-password-otp", roleMiddleware(Admin, God_eye), adminController.VerifyPwdOTP)
	api.Post("/reset-password", roleMiddleware(Admin, God_eye), adminController.ResetPassword)
	//dashboards
	api.Get("/dashboard/summary", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchDashboardSummary)
	api.Get("/analytics", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchAnalytics)
	//appointments
	api.Get("/appointments", roleMiddleware(Admin, God_eye), adminController.FetchAppointments)
	api.Post("/appointments/:id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchAppointmentByID)
	api.Patch("/appointments/:id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateAppointmentStatus)
	api.Put("/appointments/:id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateAppointment)
	//doctors
	api.Get("/doctors", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchDoctors)
	api.Get("/doctors/:doctortag", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchDoctorByID)
	api.Delete("/doctors/:doctortag", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeleteDoctor)
	//patients
	api.Get("/patients", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchPatients)
	api.Get("/patients/:usertag", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchPatientByUsertag)
	api.Delete("/patients/:usertag", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeletePatient)
	api.Patch("/patients/:usertag", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.EditPatient)
	//pharmacy
	api.Get("/pharmacy", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchPharmacy)
	api.Get("/pharmacy/:pharmacy_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchPharmacyByID)
	api.Post("/pharmacy", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.CreatePharmacy)
	api.Delete("/pharmacy/:pharmacy_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeletePharmacy)
	api.Patch("/pharmacy/:pharmacy_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdatePharmacy)
	//hospitals
	api.Get("/hospitals", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchHospitals)
	api.Get("/hospitals/:hospital_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchHospitalByID)
	api.Post("/hospitals", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.CreateHospital)
	api.Delete("/hospitals/:hospital_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeleteHospital)
	api.Patch("/hospitals/:hospital_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateHospital)
	//inventory
	api.Get("/inventory", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchInventory)
	api.Get("/inventory/:inventory_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchInventoryByID)
	api.Post("/inventory", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.CreateInventory)
	api.Delete("/inventory/:inventory_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeleteInventory)
	api.Patch("/inventory/:inventory_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateInventory)
	//orders
	api.Get("/orders", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchOrders)
	api.Get("/orders/:order_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchOrderByID)
	api.Put("/orders/:order_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateOrder)
	//test center
	api.Get("/test-centers", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchTestCenters)
	api.Get("/test-centers/:test_center_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchTestCenterByID)
	api.Post("/test-centers", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.CreateTestCenter)
	api.Delete("/test-centers/:test_center_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeleteCenter)
	api.Patch("/test-centers/:test_center_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateTestCenter)
	//reviews
	api.Get("/reviews", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchReviews)
	api.Get("/reviews/:review_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchReviewByID)
	api.Delete("/reviews/:review_id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.DeleteReview)
	//admin profile
	api.Get("/profile", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchAdminProfile)
	api.Patch("/profile", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.UpdateAdminProfile)
}

func roleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminRole := c.Get("role")
		if !contains(allowedRoles, adminRole) {
			return responses.ErrorResponse(c, "Unauthorized access", fiber.StatusForbidden)
		}
		return c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
