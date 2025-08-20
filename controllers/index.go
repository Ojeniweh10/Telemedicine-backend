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
