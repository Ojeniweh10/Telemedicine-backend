package responses

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type squadResponse struct {
	Response_code         int    `json:"response_code"`
	Transaction_reference string `json:"transaction_reference"`
	Response_description  string `json:"response_description"`
}

func ErrorResponse(c *fiber.Ctx, message string, statusCode int) error {
	res := Response{
		Success: false,
		Message: message,
	}
	return c.Status(statusCode).JSON(res)
}

func SuccessResponse(c *fiber.Ctx, message string, data any, statusCode int) error {
	res := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(res)
}

func SquadResponse(c *fiber.Ctx, reference, description string, statusCode int) error {
	res := squadResponse{
		Response_code:         statusCode,
		Transaction_reference: reference,
		Response_description:  description,
	}
	return c.Status(statusCode).JSON(res)
}

const (
	UNAUTHORIZED_ACCESS    = "unauthorized access"
	BAD_DATA               = "invalid data"
	INCOMPLETE_DATA        = "incomplete data"
	LOGIN_SUCCESSFUL       = "login successful"
	OTP_SENT               = "otp has been sent to your email"
	ACCOUNT_NON_EXISTENT   = "admin account does not exist"
	USER_NON_EXISTENT      = "user does not exist"
	INVALID_PASSWORD       = "invalid password"
	ACCOUNT_CREATED        = "account created successfully"
	OTP_VERIFIED           = "otp verified successfully"
	SOMETHING_WRONG        = "something went wrong, please try again later"
	PASSWORD_RESET_SUCCESS = "password reset successfully"
	DATA_FETCHED           = "data fetched successfully"
	DATA_UPDATED           = "data updated successfully"
	DATA_CREATED           = "data created successfully"
	DATA_DELETED           = "data deleted successfully"
	BALANCE_FETCHED	       = "wallet balance fetched successfully"
	WALLET_TOPUP_SUCCESS   = "wallet top-up successful"
	PAYMENT_VERIFIED       = "payment verified successfully"
	INVALID_SIGNATURE	   = "invalid signature"
	DATA_PROCESSED		   = "data processed successfully"
	BANKS_FETCHED          = "banks fetched successfully"
	PAYOUT_ACCOUNT_CREATED = "payout account created successfully"
	WITHDRAWAL_INITIATED     = "withdrawal initiated successfully"
	INVALID_PIN			= "invalid transaction pin"
	TRANSFER_FAILED ="transfer failed"
)
