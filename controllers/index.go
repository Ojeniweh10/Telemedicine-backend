package controllers

import (
	"telemed/models"
	"telemed/responses"
	"telemed/servers"

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
