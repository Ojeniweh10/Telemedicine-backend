package servers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"telemed/models"
	"telemed/responses"
	"telemed/utils"
	"time"

	"github.com/jackc/pgx/v4"
)

type UserServer struct{}

var walletServer WalletServer

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
	_, err = Db.Exec(Ctx, "UPDATE users SET phone_no = $1, gender = $2, date_of_birth = $3, password = $4 WHERE usertag = $5", data.Phone_no, data.Gender, data.Dob, password, data.Usertag)
	if err != nil {
		log.Println("failed to save user data", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]interface{}{
		"message": "verification successful",
		"usertag": data.Usertag,
	}, nil
}

func (UserServer) Login(data models.Login) (any, error) {
	var hash string
	var usertag string
	err := Db.QueryRow(Ctx, "SELECT password, usertag FROM users WHERE email = $1", data.Email).Scan(&hash, &usertag)
	if err != nil {
		log.Println(err)
		return nil, errors.New(responses.USER_NON_EXISTENT)
	}

	pwdCheck := utils.VerifyPassword(data.Password, hash)
	if !pwdCheck {
		log.Println("Invalid password for user login")
		return nil, errors.New(responses.INVALID_PASSWORD)
	}
	token, err := utils.GenerateJWT(usertag)
	if err != nil {
		log.Println("Failed to generate JWT token:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	}, nil
}

func (UserServer) GetDoctors() (any, error) {
	var doctors []models.Doctor
	rows, err := Db.Query(Ctx, "SELECT doctortag, fullname, date_of_birth, phone_number, gender, specialization, country, yrs_of_experience, price_per_session FROM doctors")
	if err != nil {
		log.Println("Failed to fetch doctors:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()
	for rows.Next() {
		var doctor models.Doctor
		if err := rows.Scan(&doctor.DoctorTag, &doctor.FullName, &doctor.Dob, &doctor.Phone_no, &doctor.Gender, &doctor.Specialization, &doctor.Country, &doctor.YearsOfExperience, &doctor.Price); err != nil {
			log.Println("Failed to scan doctor:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		doctors = append(doctors, doctor)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over doctors:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return doctors, nil
}

func (UserServer) BookAppointment(data models.BookAppointment) (interface{}, error) {
	var (
		count         int
		appointmentID int
		availability  []string
		resp          models.BookAppointmentResp
	)
	// Start a transaction
	tx, err := Db.BeginTx(Ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer tx.Rollback(Ctx)

	// Check existing appointment
	err = tx.QueryRow(Ctx,
		`SELECT COUNT(*) FROM appointments 
		 WHERE patient_tag = $1 AND status IN ('pending','confirmed')`,
		data.Usertag).Scan(&count)
	if err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	if count > 0 {
		return nil, errors.New("you already have an open appointment")
	}

	// Handle payment via wallet
	ref, err := walletServer.InitiateTransfer(tx, data.Usertag, data.Doctortag, data.Amount, "Appointment payment")
	if err != nil {
		return nil, err
	}

	// Check doctor availability
	err = tx.QueryRow(Ctx,
		`SELECT availability FROM doctors WHERE doctortag=$1 FOR UPDATE`,
		data.Doctortag).Scan(&availability)
	if err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	requestedTime := data.Scheduled_at.Format(time.RFC3339)
	timeAvailable := false
	for _, slot := range availability {
		if slot == requestedTime {
			timeAvailable = true
			break
		}
	}
	if !timeAvailable {
		return nil, errors.New("time slot not available")
	}

	// Insert appointment
	err = tx.QueryRow(Ctx,
		`INSERT INTO appointments (patient_tag, doctor_tag, scheduled_at, reason, status, payment_reference)
		 VALUES ($1,$2,$3,$4,'pending',$5) RETURNING appointment_id`,
		data.Usertag, data.Doctortag, data.Scheduled_at, data.Reason, ref).Scan(&appointmentID)
	if err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	// Update doctor availability
	_, err = tx.Exec(Ctx,
		`UPDATE doctors SET availability = array_remove(availability,$1) WHERE doctortag=$2`,
		requestedTime, data.Doctortag)
	if err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	// Doctor details
	err = tx.QueryRow(Ctx,
		`SELECT fullname, specialization, profile_pic_url 
		 FROM doctors WHERE doctortag=$1`,
		data.Doctortag).
		Scan(&resp.Fullname, &resp.Specialization, &resp.Doctor_photo_url)
	if err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	// Commit transaction
	if err := tx.Commit(Ctx); err != nil {
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	// Response
	resp.AppointmentID = fmt.Sprintf("%d", appointmentID)
	resp.Doctortag = data.Doctortag
	resp.Scheduled_at = data.Scheduled_at

	return resp, nil
}


func (UserServer) GetAppointments(usertag string) (any, error) {
	var resp []models.GetAppointmentsResp
	rows, err := Db.Query(Ctx, "SELECT appointment_id, patient_tag, doctor_tag, scheduled_at, reason, status FROM appointments WHERE patient_tag = $1 ORDER BY created_at DESC", usertag)
	if err != nil {
		log.Println("Failed to fetch appointments:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()
	for rows.Next() {
		var appointment models.GetAppointmentsResp
		if err := rows.Scan(&appointment.AppointmentID, &appointment.PatientTag, &appointment.DoctorTag, &appointment.Scheduled_at, &appointment.Reason, &appointment.Status); err != nil {
			log.Println("Failed to scan appointment:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		resp = append(resp, appointment)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over appointments:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return resp, nil
}

func (UserServer) RateDoctor(data models.RateDoctor) (any, error) {
	err := Db.QueryRow(Ctx,`INSERT INTO reviews (user_tag, doctor_tag, star_rating, review, status)
		 VALUES ($1, $2, $3, $4, 'pending')`,data.Usertag, data.Doctortag, data.Rating, data.Review)
	if err != nil {
		log.Printf("Failed to insert into reviews: %v", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return map[string]interface{}{
		"message": "review submitted successfully",
	}, nil
}

func (UserServer) GetMedications(data models.GetDataReq) (any, error) {
	var resp []models.GetMedicationsResp
	var args []any
	argIndex := 1
	offset := data.Limit*data.Page - data.Limit

	var sqlStatement string

	if data.Search == "" {
		sqlStatement = "SELECT product_id, name, milligram, price, product_image_url"
	} else {
		sqlStatement = fmt.Sprintf("SELECT product_id, name, milligram, price, product_image_url WHERE (product_id ILIKE $%d OR name ILIKE $%d OR milligram ILIKE $%d OR price ILIKE $%d)", argIndex, argIndex, argIndex, argIndex)
		args = append(args, "%"+data.Search+"%")
		argIndex++
	}

	sqlStatement += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, data.Limit, offset)

	rows, err := Db.Query(Ctx, sqlStatement, args...)
	if err != nil {
		log.Println("Failed to fetch medications:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()

	for rows.Next() {
		var res models.GetMedicationsResp
		if err := rows.Scan(&res.ProductID, &res.Name, &res.Milligram, &res.Price, &res.Product_Image_Url); err != nil {
			log.Println("Failed to scan medications:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		resp = append(resp, res)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over medications:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return resp, nil
}

func (UserServer) GetPharmacies(data models.GetDataReq) (any, error) {
	var resp []models.GetPharmaciesResp
	var args []any
	argIndex := 1
	offset := data.Limit*data.Page - data.Limit

	var sqlStatement string

	if data.Search == "" {
		sqlStatement = "SELECT pharmacy_id, name, address, country, state, about,  pharmacy_image_url"
	} else {
		sqlStatement = fmt.Sprintf("SELECT pharmacy_id, name, address, country, state, about,  pharmacy_image_url WHERE (name ILIKE $%d OR address ILIKE $%d OR country ILIKE $%d OR state ILIKE $%d OR about ILIKE $%d)", argIndex, argIndex, argIndex, argIndex, argIndex)
		args = append(args, "%"+data.Search+"%")
		argIndex++
	}

	sqlStatement += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, data.Limit, offset)

	rows, err := Db.Query(Ctx, sqlStatement, args...)
	if err != nil {
		log.Println("Failed to fetch pharmacies:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()

	for rows.Next() {
		var res models.GetPharmaciesResp
		if err := rows.Scan(&res.PharmacyID, &res.Name, &res.Address, &res.Country, &res.State, &res.About, &res.Pharmacy_Image_Url); err != nil {
			log.Println("Failed to scan into pharmacies:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		resp = append(resp, res)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over pharmacies:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return resp, nil
}

func (UserServer) AddToCart(data models.Cart) error {
	var Quantity, cartID int
	tx, err := Db.BeginTx(Ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return errors.New(responses.SOMETHING_WRONG)
	}
	defer tx.Rollback(Ctx) // Rollback on error

	err = tx.QueryRow(Ctx,
		`SELECT quantity FROM inventory WHERE product_id = $1`,
		data.ProductID).Scan(&Quantity)
	if err == sql.ErrNoRows {
		return fmt.Errorf("product not found")
	}
	if err != nil {
		log.Printf("Failed to query inventory: %v", err)
		return errors.New(responses.SOMETHING_WRONG)
	}
	if Quantity < data.Quantity {
		return fmt.Errorf("insufficient stock: only %d available", Quantity)
	}
	err = tx.QueryRow(Ctx,
		`INSERT INTO carts (usertag, product_id, quantity)
         VALUES ($1, $2, $3)
         ON CONFLICT (usertag, product_id)
         DO UPDATE SET quantity = carts.quantity + EXCLUDED.quantity
         RETURNING cart_id`,
		data.Usertag, data.ProductID, data.Quantity).Scan(&cartID)
	if err != nil {
		log.Printf("Failed to insert/update cart: %v", err)
		return errors.New(responses.SOMETHING_WRONG)
	}

	if err := tx.Commit(Ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return errors.New(responses.SOMETHING_WRONG)
	}

	return nil
}

func (UserServer) UpdateCart(data models.Cart) (any, error) {
	var resp models.UpdateCartResp
	var Quantity int
	tx, err := Db.BeginTx(Ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer tx.Rollback(Ctx)
	err = tx.QueryRow(Ctx,
		`SELECT quantity FROM inventory WHERE product_id = $1`,
		data.ProductID).Scan(&Quantity)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		log.Printf("Failed to query inventory: %v", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	if Quantity < data.Quantity {
		return nil, fmt.Errorf("insufficient stock: only %d available", Quantity)
	}
	_, err = tx.Exec(Ctx, `SELECT quantity FROM carts WHERE usertag = $1 AND product_id = $2 FOR UPDATE`, data.Usertag, data.ProductID)
	if err != nil {
		log.Println("Failed to fetch quantity from carts to lock row:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	err = tx.QueryRow(Ctx,
		`UPDATE carts SET quantity = $1 WHERE product_id = $2 AND usertag = $3
         RETURNING quantity`,
		data.Quantity, data.ProductID, data.Usertag).Scan(&resp.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		} else {
			log.Printf("Failed to update cart: %v", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
	}
	if err := tx.Commit(Ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	resp.ProductID = data.ProductID

	return resp, nil
}

func (UserServer) DeleteFromCart(data models.Cart) error {
	_, err := Db.Exec(Ctx, "DELETE FROM carts WHERE product_id = $1 AND usertag = $2 ", data.ProductID, data.Usertag)
	if err != nil {
		log.Println("Failed to delete product from cart:", err)
		return errors.New(responses.SOMETHING_WRONG)
	}
	return nil
}

func (UserServer) GetCart(Usertag string) (any, error) {
	var res []models.GetCartResp
	query := `
		SELECT c.cart_id, c.usertag, c.quantity,
		       i.product_id, i.name, i.milligram, i.price,  i.product_image_url
		FROM carts c
		JOIN inventory i ON i.product_id = c.product_id
		WHERE c.usertag = $1
	`
	rows, err := Db.Query(Ctx, query, Usertag)
	if err != nil {
		log.Println("Failed to fetch cart for user: ", Usertag, err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()

	for rows.Next() {
		var resp models.GetCartResp
		if err := rows.Scan(&resp.CartID, &resp.Usertag, &resp.Quantity, &resp.ProductID, &resp.Name, &resp.Milligram, &resp.Price, &resp.Product_Image_Url); err != nil {
			log.Println("Failed to scan cart item:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		res = append(res, resp)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over cart details:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return res, nil
}

func (UserServer) GetBillingDetails(Usertag string) (any, error) {
	var resp models.GetBillingDetailsResp
	query := `
		SELECT usertag, fullname, email, phone_no, state, delivery_address
		FROM billing_details
		WHERE usertag = $1
	`
	err := Db.QueryRow(Ctx, query, Usertag).Scan(
		&resp.Usertag,
		&resp.Fullname,
		&resp.Email,
		&resp.Phone_no,
		&resp.State,
		&resp.DeliveryAddress)
	if err != nil {
		log.Println("Failed to fetch billing details for user: ", Usertag, err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return resp, nil
}


func (UserServer) GetProfile(Usertag string) (any, error) {
	var resp models.UserProfile
	query := `SELECT usertag, firstname, lastname, email, phone_no, gender, date_of_birth, photo_url FROM users WHERE usertag = $1` 
	err := Db.QueryRow(Ctx, query, Usertag).Scan(
		&resp.Usertag,
		&resp.Firstname,
		&resp.Lastname,
		&resp.Email,
		&resp.Phone_no,
		&resp.Gender,
		&resp.Dob,
		&resp.Photo_url)
	if err != nil {
		log.Println("Failed to fetch user profile for user: ", Usertag, err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return resp, nil
}

func (UserServer) UpdateProfile(data models.UserProfile) (any, error) {
	fields := []string{}
	args := []interface{}{}
	argIdx := 1
	// Define all possible updates in a slice (value, columnName)
	updates := []struct {
		value string
		name  string
	}{
		{data.Firstname, "firstname"},
		{data.Lastname, "lastname"},
		{data.Email, "email"},
		{data.Phone_no, "phone_no"},
		{data.Gender, "gender"},
		{data.Dob, "date_of_birth"},
		{data.Photo_url, "photo_url"},
	}
	// Loop through updates and add only non-empty fields
	for _, u := range updates {
		if u.value != "" {
			fields = append(fields, fmt.Sprintf("%s = $%d", u.name, argIdx))
			args = append(args, u.value)
			argIdx++
		}
	}
	// If no fields provided, return error
	if len(fields) == 0 {
		return nil, errors.New("no fields to update")
	}
	// Add usertag as last argument for WHERE clause
	args = append(args, data.Usertag)
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE usertag = $%d",
		strings.Join(fields, ", "),
		argIdx,
	)

	// Execute query
	_, err := Db.Exec(Ctx, query, args...)
	if err != nil {
		log.Println("Failed to update user profile:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	// Return success response
	return map[string]interface{}{
		"message": "profile updated successfully",
	}, nil
}

func (UserServer) SendChangePasswordOTP(usertag string) (any, error) {
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


func (UserServer) ChangePassword(data models.ChangePasswordReq) (any, error) {
	var hash string
	err := Db.QueryRow(Ctx, "SELECT password FROM users WHERE usertag = $1", data.Usertag).Scan(&hash)
	if err != nil {
		log.Println(err)
		return nil, errors.New(responses.USER_NON_EXISTENT)
	}

	pwdCheck := utils.VerifyPassword(data.CurrentPassword, hash)
	if !pwdCheck {
		log.Println("Invalid old password for user")
		return nil, errors.New(responses.INVALID_PASSWORD)
	}
	newPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		log.Println("Unable to hash new password")
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	_, err = Db.Exec(Ctx, "UPDATE users SET password = $1 WHERE usertag = $2", newPassword, data.Usertag)
	if err != nil {
		log.Println("failed to update user password", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]interface{}{
		"message": "password updated successfully",
	}, nil
}
