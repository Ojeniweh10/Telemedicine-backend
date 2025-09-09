package controllers

import (
	"strconv"
	"telemed/models"
	"telemed/responses"
	"telemed/servers"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct{}

var adminServer servers.AdminServer

func (AdminController) Login(c *fiber.Ctx) error {
	var payload models.Adminlogin
	//parse data from request
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	//vet if data exists
	if payload.Email == "" || payload.Password == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	//pass data to servers
	res, err := adminServer.Login(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}

func (AdminController) VerifyOTP(c *fiber.Ctx) error {
	var payload models.OTPVerify
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.OTP == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.VerifyOTP(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_VERIFIED, res, 200)
}

func (AdminController) ForgotPassword(c *fiber.Ctx) error {
	var payload models.ForgotPassword
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.Email == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.ForgotPassword(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}

func (AdminController) VerifyPwdOTP(c *fiber.Ctx) error {
	var payload models.VerifyPwdOTP
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.OTP == "" || payload.Email == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.VerifyPwdOTP(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_VERIFIED, res, 200)
}

func (AdminController) ResetPassword(c *fiber.Ctx) error {
	var payload models.ResetPassword
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.Email == "" || payload.NewPassword == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.ResetPassword(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.PASSWORD_RESET_SUCCESS, res, 200)
}
func (AdminController) FetchDashboardSummary(c *fiber.Ctx) error {
	res, err := adminServer.GetDashboardSummary()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}
func (AdminController) FetchAnalytics(c *fiber.Ctx) error {
	var data models.AnalyticsReq
	data.Metric = c.Query("metric")
	data.Month = c.Query("month")
	data.Year = c.Query("year")

	if data.Metric == "" || data.Month == "" || data.Year == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}

	res, err := adminServer.GetAnalytics(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)

}

func (AdminController) FetchAppointments(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")

	res, err := adminServer.GetAppointments(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchAppointmentByID(c *fiber.Ctx) error {
	var payload models.AppointmentID
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.ID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetAppointmentByID(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchDoctorByID(c *fiber.Ctx) error {
	var payload models.Doctorreq
	payload.DoctorTag = c.Params("id")
	if payload.DoctorTag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetDoctorByID(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) UpdateAppointmentStatus(c *fiber.Ctx) error {
	var payload models.UpdateAppointmentStatus
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.Appointment_id = c.Params("id")
	if payload.Appointment_id == "" || payload.Status == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateAppointmentStatus(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (a *AdminController) UpdateAppointment(c *fiber.Ctx) error {
	var payload models.RescheduleAppointmentReq
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}

	payload.Appointment_id = c.Params("id")
	if payload.Appointment_id == "" || payload.NewScheduledAt == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}

	res, err := adminServer.RescheduleAppointment(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}

	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchDoctors(c *fiber.Ctx) error {
	res, err := adminServer.GetDoctors()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) DeleteDoctor(c *fiber.Ctx) error {
	var payload models.Doctorreq
	payload.DoctorTag = c.Params("id")
	if payload.DoctorTag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeleteDoctor(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}

func (AdminController) FetchPatients(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")

	res, err := adminServer.GetPatients(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchPatientByUsertag(c *fiber.Ctx) error {
	var payload models.PatientIdReq
	payload.Usertag = c.Params("usertag")

	res, err := adminServer.GetPatientByUsertag(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) DeletePatient(c *fiber.Ctx) error {
	var payload models.PatientIdReq
	payload.Usertag = c.Params("usertag")
	if payload.Usertag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeletePatient(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}

func (AdminController) EditPatient(c *fiber.Ctx) error {
	var payload models.Patient
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.UserTag = c.Params("usertag")
	if payload.UserTag == "" || payload.Firstname == "" || payload.Lastname == "" || payload.Phone_no == "" || payload.Dob == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.EditPatient(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchPharmacy(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")
	res, err := adminServer.GetPharmacy(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) CreatePharmacy(c *fiber.Ctx) error {
	var payload models.Pharmacy
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.PharmacyName == "" || payload.Address == "" || payload.Country == "" || payload.State == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.CreatePharmacy(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 201)
}

func (AdminController) DeletePharmacy(c *fiber.Ctx) error {
	pharmacyID := c.Params("pharmacy_id")
	if pharmacyID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeletePharmacy(pharmacyID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}

func (AdminController) FetchPharmacyByID(c *fiber.Ctx) error {
	pharmacyID := c.Params("pharmacy_id")
	if pharmacyID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetPharmacyByID(pharmacyID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) UpdatePharmacy(c *fiber.Ctx) error {
	var payload models.Pharmacy
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.PharmacyID = c.Params("pharmacy_id")
	if payload.PharmacyID == "" || payload.PharmacyName == "" || payload.Address == "" || payload.Country == "" || payload.State == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdatePharmacy(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchHospitals(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")
	res, err := adminServer.GetHospitals(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) CreateHospital(c *fiber.Ctx) error {
	var payload models.Hospital
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.HospitalName == "" || payload.Address == "" || payload.Country == "" || payload.State == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.CreateHospital(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 201)
}

func (AdminController) FetchHospitalByID(c *fiber.Ctx) error {
	hospitalID := c.Params("hospital_id")
	if hospitalID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetHospitalByID(hospitalID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) DeleteHospital(c *fiber.Ctx) error {
	hospitalID := c.Params("hospital_id")
	if hospitalID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeleteHospital(hospitalID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}
func (AdminController) UpdateHospital(c *fiber.Ctx) error {
	var payload models.Hospital
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.HospitalID = c.Params("hospital_id")
	if payload.HospitalID == "" || payload.HospitalName == "" || payload.Address == "" || payload.Country == "" || payload.State == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateHospital(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchInventory(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")
	res, err := adminServer.GetInventory(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchInventoryByID(c *fiber.Ctx) error {
	inventoryID := c.Params("inventory_id")
	if inventoryID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetInventoryByID(inventoryID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) CreateInventory(c *fiber.Ctx) error {
	var payload models.Inventory
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.ProductName == "" || payload.Milligrams == "" || payload.Price == 0 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.CreateInventory(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 201)
}

func (AdminController) DeleteInventory(c *fiber.Ctx) error {
	inventoryID := c.Params("inventory_id")
	if inventoryID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeleteInventory(inventoryID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}

func (AdminController) UpdateInventory(c *fiber.Ctx) error {
	var payload models.Inventory
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.ProductID = c.Params("inventory_id")
	if payload.ProductID == "" || payload.ProductName == "" || payload.Milligrams == "" || payload.Price == 0 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateInventory(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchOrders(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")
	res, err := adminServer.GetOrders(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchOrderByID(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	if orderID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetOrderByID(orderID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) UpdateOrder(c *fiber.Ctx) error {
	var payload models.Orders
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.OrderID = c.Params("order_id")
	if payload.OrderID == "" || payload.Status == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateOrder(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchTestCenters(c *fiber.Ctx) error {
	var data models.GetDataReq
	if c.Query("page") != "" {
		data.Page, _ = strconv.Atoi(c.Query("page"))
	} else {
		data.Page = 1
	}
	if c.Query("limit") != "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		data.Limit = min(limit, 100)
	} else {
		data.Limit = 100
	}

	data.Status = c.Query("status")
	data.Search = c.Query("search")
	res, err := adminServer.GetTestCenters(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchTestCenterByID(c *fiber.Ctx) error {
	testCenterID := c.Params("test_center_id")
	if testCenterID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetTestCenterByID(testCenterID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) CreateTestCenter(c *fiber.Ctx) error {
	var payload models.TestCentre
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.CentreName == "" || payload.Address == "" || payload.Country == "" || payload.State == "" || payload.TestType == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.CreateTestCenter(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 201)
}

func (AdminController) DeleteCenter(c *fiber.Ctx) error {
	testCenterID := c.Params("test_center_id")
	if testCenterID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeleteTestCenter(testCenterID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}
func (AdminController) UpdateTestCenter(c *fiber.Ctx) error {
	var payload models.TestCentre
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.CentreID = c.Params("test_center_id")
	if payload.CentreID == "" || payload.CentreName == "" || payload.Address == "" || payload.Country == "" || payload.State == "" || payload.TestType == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateTestCenter(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchReviews(c *fiber.Ctx) error {
	var payload models.Getreviews
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.Status == "" || payload.Status != "approved" && payload.Status != "pending" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetReviews(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchReviewByID(c *fiber.Ctx) error {
	reviewID := c.Params("review_id")
	if reviewID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetReviewByID(reviewID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) DeleteReview(c *fiber.Ctx) error {
	reviewID := c.Params("review_id")
	if reviewID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeleteReview(reviewID)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}

func (AdminController) FetchAdminProfile(c *fiber.Ctx) error {
	AdminTag := c.Get("admintag")
	if AdminTag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetAdminProfile(AdminTag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) UpdateAdminProfile(c *fiber.Ctx) error {
	var payload models.AdminProfile
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.AdminTag = c.Get("admintag")
	if payload.AdminTag == "" || payload.Firstname == "" || payload.Lastname == "" || payload.Email == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateAdminProfile(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}
