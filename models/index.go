package models

import (
	"encoding/json"
	"time"
)

type Verifyemail struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

type Verifyemailresp struct {
	Usertag string `json:"usertag"`
}

type Signup struct {
	Usertag  string `json:"usertag"`
	Phone_no string `json:"phone_no"`
	Gender   string `json:"gender"`
	Dob      string `json:"Dob"`
	Password string `json:"password"`
}

type SignupResp struct {
	Usertag   string `json:"usertag"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BookAppointment struct {
	Usertag      string    `json:"usertag"`
	Doctortag    string    `json:"doctortag"`
	Scheduled_at time.Time `json:"appointment_date"`
	Reason       string    `json:"reason"`
}

type BookAppointmentResp struct {
	AppointmentID    string    `json:"appointment_id"`
	Doctortag        string    `json:"doctortag"`
	Fullname         string    `json:"fullname"`
	Specialization   string    `json:"specialization"`
	Doctor_photo_url string    `json:"Doctor_photo_url"`
	Scheduled_at     time.Time `json:"appointment_date"`
}

type GetMedicationsResp struct {
	ProductID         string  `json:"product_id"`
	Name              string  `json:"name"`
	Milligram         string  `json:"milligram"`
	Price             float64 `json:"price"`
	Product_Image_Url string  `json:"Pharmacy_image_url"`
}

type GetPharmaciesResp struct {
	PharmacyID         string `json:"pharmacy_id"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	Country            string `json:"country"`
	State              string `json:"state"`
	About              string `json:"about"`
	Pharmacy_Image_Url string `json:"Pharmacy_image_url"`
}

type Cart struct {
	Usertag   string
	ProductID int
	Quantity  int `json:"quantity"`
}

type UpdateCartResp struct {
	Quantity  int `json:"quantity"`
	ProductID int `json:"product_id"`
}

type GetCartResp struct {
	Usertag           string  `json:"usertag"`
	CartID            int     `json:"cart_id"`
	ProductID         int     `json:"product_id"`
	Quantity          int     `json:"quantity"`
	Name              string  `json:"name"`
	Milligram         string  `json:"milligram"`
	Price             float64 `json:"price"`
	Product_Image_Url string  `json:"Pharmacy_image_url"`
}

type GetBillingDetailsResp struct {
	Usertag         string `json:"usertag"`
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Phone_no        string `json:"phone_no"`
	State           string `json:"state"`
	DeliveryAddress string `json:"delivery_address"`
}

type WalletResp struct {
	Wisetag  string  `json:"wisetag"`
	Balance float64 `json:"balance"`
}

type WalletTopUp struct {
	Usertag string  `json:"usertag"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WalletTopUpResp struct {
    Status  bool   `json:"status"`
    Message string `json:"message"`
    Data    struct {
        AuthorizationURL string `json:"authorization_url"`
        AccessCode       string `json:"access_code"`
        Reference        string `json:"reference"`
    } `json:"data"`
}


type Paystackreq struct {
	Email  string `json:"email"`
	Amount int    `json:"amount"`
}

type VerifyTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64   `json:"id"`
		Domain          string  `json:"domain"`
		Status          string  `json:"status"`
		Reference       string  `json:"reference"`
		ReceiptNumber   *string `json:"receipt_number"`
		Amount          int64   `json:"amount"`
		Message         *string `json:"message"`
		GatewayResponse string  `json:"gateway_response"`
		PaidAt          string  `json:"paid_at"`
		CreatedAt       string  `json:"created_at"`
		Channel         string  `json:"channel"`
		Currency        string  `json:"currency"`
		IPAddress       string  `json:"ip_address"`
		Metadata        string  `json:"metadata"`
		Log             struct {
			StartTime  int64  `json:"start_time"`
			TimeSpent  int    `json:"time_spent"`
			Attempts   int    `json:"attempts"`
			Errors     int    `json:"errors"`
			Success    bool   `json:"success"`
			Mobile     bool   `json:"mobile"`
			Input      []any  `json:"input"`
			History    []struct {
				Type    string `json:"type"`
				Message string `json:"message"`
				Time    int    `json:"time"`
			} `json:"history"`
		} `json:"log"`
		Fees      int64       `json:"fees"`
		FeesSplit interface{} `json:"fees_split"`
		Authorization struct {
			AuthorizationCode string  `json:"authorization_code"`
			Bin               string  `json:"bin"`
			Last4             string  `json:"last4"`
			ExpMonth          string  `json:"exp_month"`
			ExpYear           string  `json:"exp_year"`
			Channel           string  `json:"channel"`
			CardType          string  `json:"card_type"`
			Bank              string  `json:"bank"`
			CountryCode       string  `json:"country_code"`
			Brand             string  `json:"brand"`
			Reusable          bool    `json:"reusable"`
			Signature         string  `json:"signature"`
			AccountName       *string `json:"account_name"`
		} `json:"authorization"`
		Customer struct {
			ID                       int64       `json:"id"`
			FirstName                *string     `json:"first_name"`
			LastName                 *string     `json:"last_name"`
			Email                    string      `json:"email"`
			CustomerCode             string      `json:"customer_code"`
			Phone                    *string     `json:"phone"`
			Metadata                 interface{} `json:"metadata"`
			RiskAction               string      `json:"risk_action"`
			InternationalFormatPhone *string     `json:"international_format_phone"`
		} `json:"customer"`
		Plan             interface{} `json:"plan"`
		Split            interface{} `json:"split"`
		OrderID          interface{} `json:"order_id"`
		PaidAtAlt        string      `json:"paidAt"`
		CreatedAtAlt     string      `json:"createdAt"`
		RequestedAmount  int64       `json:"requested_amount"`
		POSTransaction   interface{} `json:"pos_transaction_data"`
		Source           interface{} `json:"source"`
		FeesBreakdown    interface{} `json:"fees_breakdown"`
		Connect          interface{} `json:"connect"`
		TransactionDate  string      `json:"transaction_date"`
		PlanObject       interface{} `json:"plan_object"`
		Subaccount       interface{} `json:"subaccount"`
	} `json:"data"`
}


type PaystackWebhook struct {
    Event string `json:"event"`
    Data  json.RawMessage `json:"data"`
}


type Bank struct {
    Name          string  `json:"name"`
    Slug          string  `json:"slug"`
    Code          string  `json:"code"`
    Longcode      string  `json:"longcode"`
    Gateway       *string `json:"gateway"`
    PayWithBank   bool    `json:"pay_with_bank"`
    Active        bool    `json:"active"`
    IsDeleted     bool    `json:"is_deleted"`
    Country       string  `json:"country"`
    Currency      string  `json:"currency"`
    Type          string  `json:"type"`
    ID            int     `json:"id"`
    CreatedAt     string  `json:"createdAt"`
    UpdatedAt     string  `json:"updatedAt"`
}

type GetBanksResponse struct {
    Status  bool   `json:"status"`
    Message string `json:"message"`
    Data    []Bank `json:"data"`
}


type PayoutAccountReq struct {
	Usertag     string `json:"usertag"`
	BankCode    string `json:"bank_code"`
	AccountNo   string `json:"account_no"`
}

type BankAccountResp struct {
	AccountName string `json:"account_name"`
	AccountNo   string `json:"account_no"`
	BankCode    string `json:"bank_code"`
}

type FetchBankAccountResp struct {
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    BankAccountResp `json:"data"`
}

type GenerateRecipientRequest struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

type GenerateRecipientResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		RecipientCode string `json:"recipient_code"`
		Details       struct {
			BankName string `json:"bank_name"`
		} `json:"details"`
	} `json:"data"`
}

type RecipientMinimal struct {
	RecipientCode string `json:"recipient_code"`
	BankName      string `json:"bank_name"`
}

type PayoutAccountResp struct {
	Usertag     string `json:"usertag"`
	BankCode    string `json:"bank_code"`
	BankName    string `json:"bank_name"`
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
	RecipientCode string `json:"account_id"`
}

type WithdrawReq struct {
	Usertag string  `json:"usertag"`
	RecipientCode string `json:"recipient_code"`
	Amount  float64 `json:"amount"`
	Transaction_pin  string  `json:"transaction_pin"`
}