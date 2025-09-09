package servers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"telemed/config"
	"telemed/models"
	"telemed/responses"
	"telemed/utils"
	"time"

	"github.com/jackc/pgx/v4"
)

type WalletServer struct{}

func (WalletServer) GetBalances(Wisetag string) (any , error) {
	var res models.WalletResp
	var balance float64
	query := `SELECT balance FROM wallets WHERE usertag=$1`
	err := Db.QueryRow(Ctx,query, Wisetag).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errors.New("wallet not found")
		}else {
			log.Println("Error fetching wallet balance:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		
	}
	res.Balance = balance
	res.Wisetag = Wisetag
	return res, nil
}

func (WalletServer) TopUp(data models.WalletTopUp) (any, error) {
	//fetch user email from users table
	var email string
	query := `SELECT email FROM users WHERE usertag=$1`
	err := Db.QueryRow(Ctx,query, data.Usertag).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}else {
			log.Println("Error fetching user email:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
	}
	//check wallet status
	var walletStatus string
	query = `SELECT wallet_status FROM wallets WHERE usertag=$1`
	err = Db.QueryRow(Ctx,query, data.Usertag).Scan(&walletStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wallet not found")
		}else {
			log.Println("Error fetching wallet status:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		
	}
	if walletStatus != "active" {
		return nil, errors.New("wallet is not active")
	}
	//converting amount to kobo
	paystackAmount := int(data.Amount * 100)
	reference := fmt.Sprintf("wallet_topup_%s_%d", data.Usertag, time.Now().Unix())
	//insert transaction record into wallet_transactions table with status pending
	query = `INSERT INTO wallet_transactions (usertag, amount,transaction_type, transaction_reference, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = Db.Exec(Ctx,query, data.Usertag, data.Amount, "credit", reference, "pending", time.Now())
	if err != nil {
		log.Println("Error inserting wallet transaction:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	payload := map[string]interface{}{
		"email":        email,
		"amount":       paystackAmount,
		"reference":    reference,
		"metadata":     map[string]string{"usertag": data.Usertag},
		"callback_url": config.PaystackCallbackURL,
		"channels":     []string{"card", "bank_transfer"},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling Paystack payload:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	url := config.PaystackBaseURL + "/transaction/initialize"
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		log.Println("Error creating Paystack request:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	req.Header.Set("Authorization", "Bearer "+config.PaystackSecretKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error initializing Paystack transaction:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
        log.Println("Paystack returned non-200 status:", resp.StatusCode)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Paystack response body:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	var response models.WalletTopUpResp
	json.Unmarshal(resBody, &response)
	if !response.Status {
        log.Println("Paystack initialization failed:", response.Message)
        return nil, errors.New(response.Message)
    }
	_, err = Db.Exec(Ctx,`UPDATE wallet_transactions SET access_code=$1, paystack_reference = $2 WHERE transaction_reference=$3`, response.Data.AccessCode, response.Data.Reference, reference)
	if err != nil {
		log.Println("Error updating wallet transaction with access code:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	Authorization_url := response.Data.AuthorizationURL
	return Authorization_url, nil
}

func (WalletServer) VerifyPayment(reference string) (bool, error) {
	url := config.PaystackBaseURL + "/transaction/verify/" + reference
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating Paystack verify request:", err)
        return false, errors.New(responses.SOMETHING_WRONG)
    }
    req.Header.Set("Authorization", "Bearer "+config.PaystackSecretKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error verifying Paystack transaction:", err)
        return false, errors.New(responses.SOMETHING_WRONG)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Println("Paystack verify returned non-200 status:", resp.StatusCode)
        return false, errors.New(responses.SOMETHING_WRONG)
    }

    resBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading Paystack verify response body:", err)
        return false, errors.New(responses.SOMETHING_WRONG)
    }

    var response models.VerifyTransactionResponse
    if err := json.Unmarshal(resBody, &response); err != nil {
        log.Println("Error unmarshaling verify response:", err)
        return false, errors.New(responses.SOMETHING_WRONG)
    }

    if !response.Status || response.Data.Status != "success" {
        return false, errors.New("transaction failed")
    }

    return true, nil
}


func (WalletServer) VerifyWebhook(eventData map[string]interface{}, signature string) bool {
    payload, _ := json.Marshal(eventData)
    hash := hmac.New(sha512.New, []byte(config.PaystackSecretKey))
    hash.Write(payload)
    expectedSignature := hex.EncodeToString(hash.Sum(nil))
    return signature == expectedSignature
}


func (WalletServer) HandleWebhook(eventData map[string]interface{}) error {
    event := eventData["event"].(string)

    switch event {
    case "charge.success":
        return handleChargeSuccess(eventData["data"])
    case "charge.dispute.create":
        return handleDisputeCreate(eventData["data"])
    case "charge.dispute.remind":
        return handleDisputeRemind(eventData["data"])
    case "charge.dispute.resolve":
        return handleDisputeResolve(eventData["data"])
    case "transfer.success":
        return handleTransferSuccess(eventData["data"])
    case "transfer.failed":
        return handleTransferFailed(eventData["data"])
    case "transfer.reversed":
        return handleTransferReversed(eventData["data"])
    default:
        log.Println("Unhandled webhook event:", event)
    }

    return nil
}


func handleChargeSuccess(data interface{}) error {
    d := data.(map[string]interface{})
    reference := d["reference"].(string)
    amount := d["amount"].(float64) / 100 // Paystack sends kobo
   // Begin transaction
    tx, err := Db.BeginTx(Ctx, pgx.TxOptions{})
    if err != nil {
        return err
    }
    defer tx.Rollback(Ctx)
    var usertag string
    err = tx.QueryRow(Ctx, `SELECT usertag FROM wallet_transactions WHERE transaction_reference=$1`, reference).Scan(&usertag)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("Transaction not found:", reference)
            return fmt.Errorf("transaction not found: %s", reference)
        }
        log.Println("Error querying transaction:", err)
        return err
    }
    _, err = tx.Exec(Ctx, `UPDATE wallet_transactions SET status='success' WHERE transaction_reference=$1`, reference)
    if err != nil {
        log.Println("Error updating transaction status:", err)
        return err
    }
    _, err = tx.Exec(Ctx, `UPDATE wallets SET balance = balance + $1 WHERE usertag=$2`, amount, usertag)
    if err != nil {
        log.Println("Error crediting wallet:", err)
        return err
    }
    if err = tx.Commit(Ctx); err != nil {
        log.Println("Error committing transaction:", err)
        return err
    }
    log.Println("Wallet funded for user:", usertag, "amount:", amount)
    return nil
}

func handleDisputeCreate(data interface{}) error {
    d := data.(map[string]interface{})
    reference := d["transaction"].(map[string]interface{})["reference"].(string)
    _, err := Db.Exec(Ctx, `UPDATE wallet_transactions SET status='disputed' WHERE transaction_reference=$1`, reference)
    if err != nil {
        return err
    }
    log.Println("Dispute created for transaction:", reference)
    return nil
}

func handleDisputeRemind(data interface{}) error {
    d := data.(map[string]interface{})
    disputeId := d["id"]
    log.Println("Reminder: respond to dispute ID:", disputeId)
    return nil
}


func handleDisputeResolve(data interface{}) error {
    d := data.(map[string]interface{})
    status := d["status"].(string) // "won" or "lost"
    reference := d["transaction"].(map[string]interface{})["reference"].(string)

    if status == "lost" {
        // Mark transaction as reversed and debit wallet
        var amount float64
        var usertag string
		tx, err := Db.BeginTx(Ctx, pgx.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback(Ctx)
        err = tx.QueryRow(Ctx, `SELECT amount, usertag FROM wallet_transactions WHERE transaction_reference=$1`, reference).Scan(&amount, &usertag)
        if err != nil {
            return err
        }
        _, err = tx.Exec(Ctx, `UPDATE wallet_transactions SET status='reversed' WHERE transaction_reference=$1`, reference)
        if err != nil {
            return err
        }
        _, err = tx.Exec(Ctx, `UPDATE wallets SET balance = balance - $1 WHERE usertag=$2`, amount, usertag)
        if err != nil {
            return err
        }
		if err = tx.Commit(Ctx); err != nil {
        log.Println("Error committing transaction:", err)
        return err
    	}
        log.Println("Dispute lost. Funds reversed for user:", usertag)
    } else {
        // Dispute won → keep transaction successful
        _, err := Db.Exec(Ctx, `UPDATE wallet_transactions SET status='success' WHERE transaction_reference=$1`, reference)
        if err != nil {
            return err
        }

        log.Println("Dispute won. Transaction remains successful:", reference)
    }
    return nil
}

func handleTransferSuccess(data interface{}) error {
    d := data.(map[string]interface{})
    reference := d["reference"].(string)

    var amount float64
    var usertag string
    err := Db.QueryRow(Ctx, `SELECT amount, usertag FROM wallet_transactions WHERE transaction_reference=$1`, reference).
        Scan(&amount, &usertag)
    if err != nil {
        log.Println("Transfer success webhook: transaction not found:", reference)
        return err
    }

    tx, err := Db.Begin(Ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(Ctx)

    _, err = tx.Exec(Ctx,
        `UPDATE wallet_transactions SET status='success' WHERE transaction_reference=$1`, reference)
    if err != nil {
        return err
    }
    _, err = tx.Exec(Ctx,
        `UPDATE wallets SET pending_balance = pending_balance - $1 WHERE usertag=$2`, amount, usertag)
    if err != nil {
        return err
    }

    if err = tx.Commit(Ctx); err != nil {
        return err
    }

    log.Println("✅ Transfer successful. User:", usertag, "amount finalized:", amount)
    return nil
}


func handleTransferFailed(data interface{}) error {
    d := data.(map[string]interface{})
    reference := d["reference"].(string)

    tx, err := Db.Begin(Ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(Ctx)

    var amount float64
    var usertag string
    err = tx.QueryRow(Ctx,
        `SELECT amount, usertag FROM wallet_transactions WHERE transaction_reference=$1`,
        reference).Scan(&amount, &usertag)
    if err != nil {
        return err
    }

    _, err = tx.Exec(Ctx,
        `UPDATE wallet_transactions SET status='failed' WHERE transaction_reference=$1`,
        reference)
    if err != nil {
        return err
    }
    // Refund + clean pending balance
    _, err = tx.Exec(Ctx,
        `UPDATE wallets SET balance = balance + $1, pending_balance = pending_balance - $1 
         WHERE usertag=$2`, amount, usertag)
    if err != nil {
        return err
    }

    if err = tx.Commit(Ctx); err != nil {
        log.Println("Error committing transaction:", err)
        return err
    }

    log.Println("❌ Transfer failed. Refunded user:", usertag, "amount:", amount)
    return nil
}


func handleTransferReversed(data interface{}) error {
    d := data.(map[string]interface{})
    reference := d["reference"].(string)
	tx, err := Db.BeginTx(Ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(Ctx)
    var amount float64
    var usertag, currentStatus string
    err = tx.QueryRow(Ctx, `SELECT status, amount, usertag FROM wallet_transactions WHERE transaction_reference=$1`, reference).Scan(&currentStatus, &amount, &usertag)
    if err != nil {
        return err
    }
    if currentStatus == "reversed" {
        log.Println("Transfer already reversed, skipping:", reference)
        return nil
    }
    _, err = tx.Exec(Ctx, `UPDATE wallet_transactions SET status='reversed' WHERE transaction_reference=$1`, reference)
    if err != nil {
        return err
    }
    // Refund wallet
    _, err = tx.Exec(Ctx, `UPDATE wallets SET balance = balance + $1 WHERE usertag=$2`, amount, usertag)
    if err != nil {
        return err
    }
	if err = tx.Commit(Ctx); err != nil {
        log.Println("Error committing transaction:", err)
        return err
    }
    log.Println("Transfer reversed. Refunded user:", usertag, "amount:", amount)
    return nil
}


func (WalletServer) GetBanks() (any, error) {
    url := config.PaystackBaseURL + "/bank?country=nigeria"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating Paystack get banks request:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }
    req.Header.Set("Authorization", "Bearer "+config.PaystackSecretKey)

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error fetching Paystack banks:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Println("Paystack get banks returned non-200 status:", resp.StatusCode)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    resBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading Paystack get banks response body:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    var response models.GetBanksResponse
    if err := json.Unmarshal(resBody, &response); err != nil {
        log.Println("Error unmarshaling banks response:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    if !response.Status {
        return nil, errors.New("failed to fetch banks from Paystack")
    }

    return response.Data, nil
}

func (WalletServer) CreatePayoutAccount(data models.PayoutAccountReq) (any, error) {
    //check how many payout accounts the user has
    var count int
    err := Db.QueryRow(Ctx, `SELECT COUNT(*) FROM payout_accounts WHERE usertag=$1 and is_active = true`, data.Usertag).Scan(&count)
    if err != nil {
        log.Println("Error counting payout accounts:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }
    if count >= 3 {
        return nil, errors.New("maximum of 3 payout accounts allowed")
    }
    account_name, err := utils.ResolveAccountNumber(data)
    if err != nil {
        return nil, err
    }
    resp, err := utils.GenerateAccountRecipient(data, account_name)
    if err != nil {
        return nil, err
    }
    _, err = Db.Exec(Ctx,
        `INSERT INTO payout_accounts (usertag, account_name, account_number, bank_code, recipient_code, bank_name, is_active, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        data.Usertag, account_name, data.AccountNo, data.BankCode,
        resp.RecipientCode, resp.BankName, true, time.Now(),
    )
    if err != nil {
        log.Println("Error inserting payout account:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }
    res := map[string]string{
        "account_name":   account_name,
        "account_number": data.AccountNo,
        "bank_name":      resp.BankName,
    }
    return res, nil
}


func (WalletServer) FetchPayoutAccounts(usertag string) (any, error) {
    rows, err := Db.Query(Ctx, `SELECT account_name, account_number, bank_code, bank_name, recipient_code FROM payout_accounts WHERE usertag=$1 AND is_active=true`, usertag)
    if err != nil {
        log.Println("Error fetching payout accounts:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }
    defer rows.Close()

    var accounts []models.PayoutAccountResp
    for rows.Next() {
        var acc models.PayoutAccountResp
        err := rows.Scan(&acc.AccountName, &acc.AccountNo, &acc.BankCode, &acc.BankName, &acc.RecipientCode)
        if err != nil {
            log.Println("Error scanning payout account row:", err)
            return nil, errors.New(responses.SOMETHING_WRONG)
        }
        accounts = append(accounts, acc)
    }

    if err = rows.Err(); err != nil {
        log.Println("Row iteration error:", err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    return accounts, nil
}

func (WalletServer) Withdraw(data models.WithdrawReq) (any, error) {
    var hash string
    err := Db.QueryRow(Ctx, "SELECT transaction_pin FROM users WHERE usertag = $1", data.Usertag).Scan(&hash)
    if err != nil {
        log.Println(err)
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    // Verify pin
    if !utils.VerifyPassword(data.Transaction_pin, hash) {
        log.Println("Invalid transaction pin for user:", data.Usertag)
        return nil, errors.New(responses.INVALID_PIN)
    }

    // Check wallet balances
    var balance, pendingBalance float64
    err = Db.QueryRow(Ctx, `SELECT balance, pending_balance FROM wallets WHERE usertag=$1`, data.Usertag).
        Scan(&balance, &pendingBalance)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("wallet not found")
        }
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    if balance < data.Amount {
        return nil, errors.New("insufficient wallet balance")
    }

    reference := fmt.Sprintf("wallet_withdrawal_%s_%d", data.Usertag, time.Now().Unix())

    // Reserve funds: move to pending balance
    tx, err := Db.Begin(Ctx)
    if err != nil {
        return nil, errors.New(responses.SOMETHING_WRONG)
    }
    defer tx.Rollback(Ctx)

    _, err = tx.Exec(Ctx,
        `UPDATE wallets 
         SET balance = balance - $1, pending_balance = pending_balance + $1 
         WHERE usertag=$2`, data.Amount, data.Usertag)
    if err != nil {
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    // Insert transaction record with "initiated"
    _, err = tx.Exec(Ctx,
        `INSERT INTO wallet_transactions (usertag, amount, transaction_type, transaction_reference, status, created_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        data.Usertag, data.Amount, "debit", reference, "initiated", time.Now())
    if err != nil {
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    if err = tx.Commit(Ctx); err != nil {
        return nil, errors.New(responses.SOMETHING_WRONG)
    }

    // Call Paystack with retry logic
    transferCode, err := utils.InitiateTransferWithRetry(data, reference, 3)
    if err != nil {
        log.Println("Error initiating transfer after retries:", err)
        // Optional: mark transaction as failed here
        return nil, errors.New("could not initiate transfer, please try again later")
    }

    // Update transaction with transfer code
    _, err = Db.Exec(Ctx,
        `UPDATE wallet_transactions SET transfer_code=$1, status='pending' WHERE transaction_reference=$2`,
        transferCode, reference)
    if err != nil {
        log.Println("Error updating wallet transaction with transfer code:", err)
    }

    return map[string]string{
        "message": "Withdrawal initiated, funds reserved in pending balance.",
        "reference": reference,
    }, nil
}
