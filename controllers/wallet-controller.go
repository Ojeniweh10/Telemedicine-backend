package controllers

import (
	"telemed/models"
	"telemed/responses"
	"telemed/servers"

	"github.com/gofiber/fiber/v2"
)

type WalletController struct{}

var walletServer servers.WalletServer

func (WalletController) FetchBalance(c *fiber.Ctx) error {
	Wisetag := c.Locals("usertag").(string)

	if Wisetag == ""{
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := walletServer.GetBalances(Wisetag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.BALANCE_FETCHED, res, 200)
}

func (WalletController) TopUp(c *fiber.Ctx) error {
	var data models.WalletTopUp
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)

	if data.Usertag == "" || data.Amount <= 0 {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := walletServer.TopUp(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.WALLET_TOPUP_SUCCESS, res, 200)
}

func (WalletController) PaymentCallback(c *fiber.Ctx) error {
	reference := c.Query("reference")
    if reference == "" {
        return c.Status(400).Render("error", fiber.Map{
            "message": "Reference not found in query parameter",
            "redirect": "/wallet/topup",
        })
    }
	res, err := walletServer.VerifyPayment(reference)
	if err != nil {
		return c.Status(400).Render("error", fiber.Map{
            "message": err.Error(),
            "redirect": "/wallet/topup",
        })
	}
	return c.Render("success", fiber.Map{
        "message": "Payment successful! Redirecting to wallet...",
        "redirect": "/wallet",
        "status":  res,
    })
}

func (WalletController) PaystackWebhook(c *fiber.Ctx) error {
    var eventData map[string]interface{}
    if err := c.BodyParser(&eventData); err != nil {
        return responses.ErrorResponse(c, responses.BAD_DATA, 400)
    }

    // Verify signature
    signature := c.Get("x-paystack-signature")
    if !walletServer.VerifyWebhook(eventData, signature) {
        return responses.ErrorResponse(c, responses.INVALID_SIGNATURE, 400)
    }

    // Pass to server for handling
    if err := walletServer.HandleWebhook(eventData); err != nil {
        return responses.ErrorResponse(c, err.Error(), 400)
    }

    return responses.SuccessResponse(c, responses.DATA_PROCESSED, nil, 200)
}

func (WalletController) FetchBanks(c *fiber.Ctx) error {
	res, err := walletServer.GetBanks()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.BANKS_FETCHED, res, 200)
}

func (WalletController) CreatePayoutAccount(c *fiber.Ctx) error {
	var data models.PayoutAccountReq
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)

	if data.Usertag == "" || data.AccountNo == "" || data.BankCode == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := walletServer.CreatePayoutAccount(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.PAYOUT_ACCOUNT_CREATED, res, 200)
}

func (WalletController) FetchPayoutAccounts(c *fiber.Ctx) error {
	usertag := c.Locals("usertag").(string)

	if usertag == ""{
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := walletServer.FetchPayoutAccounts(usertag)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (WalletController) Withdraw(c *fiber.Ctx) error {
	var data models.WithdrawReq
	if err := c.BodyParser(&data); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	data.Usertag = c.Locals("usertag").(string)

	if data.Usertag == "" || data.RecipientCode == "" || data.Amount <= 0 || data.Transaction_pin == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := walletServer.Withdraw(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.WITHDRAWAL_INITIATED, res, 200)
}

