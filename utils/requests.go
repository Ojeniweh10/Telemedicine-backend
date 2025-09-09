package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"telemed/config"
	"telemed/models"
	"telemed/responses"
	"time"
)


func ResolveAccountNumber(data models.PayoutAccountReq) (string, error) {
	url := config.PaystackBaseURL + "/bank/resolve?account_number=" + data.AccountNo + "&bank_code=" + data.BankCode
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating Paystack fetch bank account request:", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }
    req.Header.Set("Authorization", "Bearer "+config.PaystackSecretKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error fetching bank account:", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Println("Paystack fetch account returned non-200 status:", resp.StatusCode)
        return "", errors.New(responses.SOMETHING_WRONG)
    }
    resBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading Paystack fetch bank account response body:", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }
    var response models.FetchBankAccountResp
    if err := json.Unmarshal(resBody, &response); err != nil {
        log.Println("Error unmarshaling banks response:", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }

    if !response.Status {
        return "", errors.New("failed to fetch bank account from Paystack")
    }
	accountname := response.Data.AccountName
    return accountname, nil
}

func GenerateAccountRecipient(data models.PayoutAccountReq, accountName string) (*models.RecipientMinimal, error) {
	url := config.PaystackBaseURL + "/transferrecipient"

	reqBody := models.GenerateRecipientRequest{
		Type:          "nuban",
		Name:          accountName,
		AccountNumber: data.AccountNo,
		BankCode:      data.BankCode,
		Currency:      "NGN",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("❌ Error marshaling recipient request:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("❌ Error creating Paystack recipient request:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	req.Header.Set("Authorization", "Bearer "+config.PaystackSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Error sending Paystack recipient request:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("❌ Paystack recipient creation returned non-200 status:", resp.StatusCode)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("❌ Error reading Paystack recipient response body:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	var response models.GenerateRecipientResponse
	if err := json.Unmarshal(resBody, &response); err != nil {
		log.Println("❌ Error unmarshaling recipient response:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	if !response.Status {
		return nil, errors.New("failed to create transfer recipient")
	}

	result := &models.RecipientMinimal{
		RecipientCode: response.Data.RecipientCode,
		BankName:      response.Data.Details.BankName,
	}

	return result, nil
}

func InitiateTransferWithRetry(data models.WithdrawReq, reference string, maxRetries int) (string, error) {
    var transferCode string
    var err error

    for attempt := 1; attempt <= maxRetries; attempt++ {
        transferCode, err = InitiateTransfer(data, reference)
        if err == nil {
            return transferCode, nil
        }
        log.Printf("⚠️ Transfer attempt %d/%d failed: %v\n", attempt, maxRetries, err)
        time.Sleep(time.Duration(attempt*2) * time.Second) // Exponential backoff
    }

    return "", fmt.Errorf("transfer failed after %d retries: %w", maxRetries, err)
}


func InitiateTransfer(data models.WithdrawReq, reference string) (string, error) {
    url := config.PaystackBaseURL + "/transfer"
    amountKobo := int64(data.Amount * 100)
    reqBody := map[string]interface{}{
        "source":    "balance",
        "amount":    amountKobo,
        "reference": reference,
        "recipient": data.RecipientCode,
        "reason":    "Wallet withdrawal",
    }

    body, err := json.Marshal(reqBody)
    if err != nil {
        log.Printf("❌ Error marshalling transfer request: %v\n", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }

    // Create request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        log.Printf("❌ Error creating transfer request: %v\n", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }
    req.Header.Set("Authorization", "Bearer "+config.PaystackSecretKey)
    req.Header.Set("Content-Type", "application/json")

    // Send request
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("❌ Error sending transfer request: %v\n", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }
    defer resp.Body.Close()

    // Handle non-200 responses
    if resp.StatusCode != http.StatusOK {
        log.Printf("❌ Paystack transfer returned non-200 status: %d\n", resp.StatusCode)
        return "", errors.New(responses.SOMETHING_WRONG)
    }

    // Read response
    resBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("❌ Error reading transfer response: %v\n", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }

    // Parse into struct
    var response struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data struct {
            TransferCode string `json:"transfer_code"`
            Status       string `json:"status"`
            Reference    string `json:"reference"`
        } `json:"data"`
    }

    if err := json.Unmarshal(resBody, &response); err != nil {
        log.Printf("❌ Error unmarshaling transfer response: %v\n", err)
        return "", errors.New(responses.SOMETHING_WRONG)
    }

    if !response.Status {
        return "", fmt.Errorf("transfer failed: %s", response.Message)
    }

    // Return the transfer code (important for later reconciliation/verification)
    return response.Data.TransferCode, nil
}
