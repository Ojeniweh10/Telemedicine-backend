package servers

import (
	"errors"
	"log"
	"telemed/models"
	"telemed/responses"
	"telemed/utils"
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
