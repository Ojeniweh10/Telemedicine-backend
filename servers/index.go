package servers

import (
	"errors"
	"log"
	"telemed/models"
	"telemed/responses"
	"telemed/utils"
	"time"
)

type UserServer struct{}

func (UserServer) Verifyemail(data models.Verifyemail) (any, error) {
	var resp models.Verifyemailresp
	resp.Usertag = utils.GenerateUsertag(data.Firstname)
	_, err := Db.Exec(Ctx, "Insert into users(usertag, firstname, lastname, email) VALUES($1, $2, $3, $4)", resp.Usertag, data.Firstname, data.Lastname, data.Email)
	if err != nil {
		log.Println("failed to save user data", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Println("Failed to generate OTP:", err)
		return nil, errors.New("failed to generate OTP")
	}
	_, err = Db.Exec(Ctx, "UPDATE users SET otp = $1, otp_expiry = NOW()+ INTERVAL '5 minutes' WHERE usertag = $2", otp, resp.Usertag)
	if err != nil {
		log.Println("failed to save OTP", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	err = utils.SendEmailOTP(data.Email, otp)
	if err != nil {
		log.Println("Failed to send OTP email:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return resp, nil
}

func (UserServer) SendEmailOTP(usertag string) (any, error) {
	var Email string
	if err := Db.QueryRow(Ctx, "SELECT email FROM users WHERE usertag = $1", usertag).Scan(&Email); err != nil {
		return nil, errors.New("email not found")
	}
	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Println("Failed to generate OTP:", err)
		return nil, errors.New("failed to generate OTP")
	}
	_, err = Db.Exec(Ctx, "UPDATE users SET otp = $1, otp_expiry = NOW()+ INTERVAL '5 minutes' WHERE usertag = $2", otp, usertag)
	if err != nil {
		log.Println("failed to save OTP", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	err = utils.SendEmailOTP(Email, otp)
	if err != nil {
		log.Println("Failed to send OTP email:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return usertag, nil
}

func (UserServer) VerifyOTP(data models.OTPVerify) (any, error) {
	var dbOtp string
	var otpExpiryTime time.Time
	err := Db.QueryRow(Ctx, "SELECT otp, otp_expiry, role FROM users WHERE usertag = $1", data.Usertag).Scan(&dbOtp, &otpExpiryTime)
	if err != nil {
		log.Println(err)
		return nil, errors.New("invalid email or OTP")
	}

	if data.OTP != dbOtp {
		log.Println("Invalid OTP for admin login")
		return nil, errors.New("invalid OTP")
	}

	if time.Now().After(otpExpiryTime) {
		log.Println("OTP has expired")
		return nil, errors.New("OTP has expired")
	}
	_, err = Db.Exec(Ctx, `UPDATE users SET otp = NULL, otp_expiry = NULL WHERE usertag = $1`, data.Usertag)
	if err != nil {
		log.Println("Failed to clear OTP:", err)
	}

	return map[string]interface{}{
		"message": "verification successful",
	}, nil
}

func (UserServer) Signup(data models.Signup) (any, error) {
	password, err := utils.HashPassword(data.Password)
	if err != nil {
		log.Println("Unable to hash password")
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	_, err = Db.Exec(Ctx, "Insert into users(phone_no, gender, dob, password) VALUES($1, $2, $3, $4)", data.Usertag, data.Phone_no, data.Gender, data.Dob, password)
	if err != nil {
		log.Println("failed to save user data", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]interface{}{
		"message": "verification successful",
		"usertag": data.Usertag,
	}, nil
}
