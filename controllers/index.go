package controllers

import (
	"errors"
	"log"
	"strconv"
	"telemed/models"
	"telemed/responses"
	"telemed/servers"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct{}

var UserServer servers.UserServer

func (Controller) Verifyemail(c *fiber.Ctx) error {
	var body models.Verifyemail
	if err := c.BodyParser(&body); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if body.Firstname == "" || body.Lastname == "" || body.Email == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.Verifyemail(body)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}

func (Controller) SendEmailOTP(c *fiber.Ctx) error {
	usertag := c.Get("usertag")
	if usertag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.SendEmailOTP(usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)

}

func (Controller) VerifyOtp(c *fiber.Ctx) error {
	var payload models.OTPVerify
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.OTP == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.VerifyOTP(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_VERIFIED, res, 200)
}

func (Controller) Signup(c *fiber.Ctx) error {
	var body models.Signup
	if err := c.BodyParser(&body); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	body.Usertag = c.Get("usertag")
	if body.Usertag == "" || body.Phone_no == "" || body.Gender == "" || body.Dob == "" || body.Password == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.Signup(body)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 200)
}

func (Controller) Login(c *fiber.Ctx) error {
	var payload models.Login
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.Email == "" || payload.Password == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.Login(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.LOGIN_SUCCESSFUL, res, 200)
}

func (Controller) FetchDoctors(c *fiber.Ctx) error {
	res, err := UserServer.GetDoctors()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (Controller) BookAppointment(c *fiber.Ctx) error {
	var data models.BookAppointment
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)

	if data.Doctortag == "" || data.Reason == "" || data.Amount <= 0 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if data.Scheduled_at.Before(time.Now()) {
		return errors.New("appointment date must be in the future")
	}
	res, err := UserServer.BookAppointment(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 200)
}

func (Controller) FetchAppointment(c *fiber.Ctx) error {
	usertag := c.Locals("usertag").(string)
	res, err := UserServer.GetAppointments(usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (Controller) RateDoctor(c *fiber.Ctx) error {
	var data models.RateDoctor
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)

	if data.Doctortag == "" || data.Rating < 1 || data.Rating > 5 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}

	res, err := UserServer.RateDoctor(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 200)
}

func (Controller) FetchMedications(c *fiber.Ctx) error {
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

	data.Status = c.Query("search")

	res, err := UserServer.GetMedications(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 200)
}

func (Controller) FetchPharmacies(c *fiber.Ctx) error {
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
	res, err := UserServer.GetPharmacies(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, res, 200)
}

func (Controller) AddToCart(c *fiber.Ctx) error {
	var data models.Cart
	var err error
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	ProductID := c.Params("product_id")
	data.Usertag = c.Locals("usertag").(string)
	if ProductID == "" || data.Quantity <= 0 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	data.ProductID, err = strconv.Atoi(ProductID)
	if err != nil {
		log.Println("failed to convert string to int for product id", err)
		return responses.ErrorResponse(c, responses.SOMETHING_WRONG, 400)
	}
	err = UserServer.AddToCart(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_CREATED, nil, 200)
}

func (Controller) UpdateCart(c *fiber.Ctx) error {
	var data models.Cart
	var err error
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	ProductID := c.Params("product_id")
	data.Usertag = c.Locals("usertag").(string)
	if ProductID == "" || data.Quantity <= 0 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	data.ProductID, err = strconv.Atoi(ProductID)
	if err != nil {
		log.Println("failed to convert string to int for product id", err)
		return responses.ErrorResponse(c, responses.SOMETHING_WRONG, 400)
	}
	res, err := UserServer.UpdateCart(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)

}
func (Controller) DeleteFromCart(c *fiber.Ctx) error {
	var data models.Cart
	var err error
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	ProductID := c.Params("product_id")
	data.Usertag = c.Locals("usertag").(string)
	if ProductID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	data.ProductID, err = strconv.Atoi(ProductID)
	if err != nil {
		log.Println("failed to convert string to int for product id", err)
		return responses.ErrorResponse(c, responses.SOMETHING_WRONG, 400)
	}
	err = UserServer.DeleteFromCart(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_DELETED, nil, 200)
}

func (Controller) FetchCart(c *fiber.Ctx) error {
	Usertag := c.Locals("usertag").(string)
	res, err := UserServer.GetCart(Usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)

}

func (Controller) FetchBillingDetails(c *fiber.Ctx) error {
	Usertag := c.Locals("usertag").(string)
	res, err := UserServer.GetBillingDetails(Usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (Controller) FetchProfile(c *fiber.Ctx) error {
	Usertag := c.Locals("usertag").(string)
	res, err := UserServer.GetProfile(Usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (Controller) UpdateProfile(c *fiber.Ctx) error {
	var data models.UserProfile
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)
	res, err := UserServer.UpdateProfile(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)	
}
func (Controller) SendChangePasswordOTP(c *fiber.Ctx) error {
	usertag := c.Locals("usertag").(string)
	if usertag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.SendChangePasswordOTP(usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}

func (Controller) ChangePassword(c *fiber.Ctx) error {
	var data models.ChangePasswordReq
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)
	if data.CurrentPassword == "" || data.NewPassword == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := UserServer.ChangePassword(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}