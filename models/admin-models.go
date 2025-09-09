package models

import (
	"time"

	"gorm.io/datatypes"
)

type Adminlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Usertag string `json:"usertag"`
}

type OTPVerify struct {
	OTP     string `json:"otp"`
	Usertag string `json:"usertag"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type VerifyPwdOTP struct {
	OTP   string `json:"otp"`
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

type AnalyticsReq struct {
	Metric string `json:"metric"`
	Month  string `json:"month"`
	Year   string `json:"year"`
}

type AnalyticsResp struct {
	Metric          string  `json:"metric"`
	Month           string  `json:"month"`
	Year            string  `json:"year"`
	Total_amount    float64 `json:"total_amount"`
	Payment_count   int     `json:"payment_count"`
	Average_payment float64 `json:"average_payment"`
}

type Appointment struct {
	ID           int       `json:"id"`
	UserTag      string    `json:"usertag"`
	DoctorTag    string    `json:"doctortag"`
	Scheduled_at time.Time `json:"appointment_date"`
	Reason       string    `json:"reason"`
	Status       string    `json:"status"`
	Fileurl      string    `json:"fileurl"`
	Created_at   time.Time `json:"created_at"`
}

type AppointmentID struct {
	ID string `json:"id"`
}

type Userdata struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
	Dob       string `json:"dob"`
	Phone_no  string `json:"phone_no"`
}

type Doctordata struct {
	Fullname string  `json:"fullname"`
	Price    float64 `json:"price"`
}

type AppointmentIDResp struct {
	UserTag         string    `json:"usertag"`
	DoctorTag       string    `json:"doctortag"`
	Scheduled_At    string    `json:"appointment_date"`
	Reason          string    `json:"reason"`
	File_URL        string    `json:"fileurl"`
	Status          string    `json:"status"`
	Created_At      time.Time `json:"created_at"`
	First_name      string    `json:"firstname"`
	Last_name       string    `json:"lastname"`
	Phone_No        string    `json:"phone_no"`
	Gender          string    `json:"gender"`
	Dob             string    `json:"dob"`
	Doctor_Fullname string    `json:"doctor_fullname"`
	Price           float64   `json:"price"`
}

type Doctorreq struct {
	DoctorTag string
}

type Doctor struct {
	DoctorTag           string         `json:"doctortag"`
	FullName            string         `json:"fullname"`
	Dob                 string         `json:"date_of_birth"`
	Phone_no            string         `json:"phone_number"`
	Gender              string         `json:"gender"`
	Specialization      string         `json:"specialization"`
	Country             string         `json:"country"`
	City                string         `json:"city"`
	YearsOfExperience   int            `json:"yrs_of_experience"`
	Price               float64        `json:"price_per_session"`
	About               string         `json:"about"`
	Availability        datatypes.JSON `json:"availability"` // or []string if unmarshalled
	ProfilePicURL       string         `json:"profile_pic_url"`
	HospitalAffiliation string         `json:"hospital_affiliation"` // from hospital.name
}

type UpdateAppointmentStatus struct {
	Status         string `json:"status"`
	Appointment_id string `json:"appointment_id"`
}

type RescheduleAppointmentReq struct {
	Appointment_id string `json:"appointment_id"`
	NewScheduledAt string `json:"new_scheduled_at"`
}

type Patient struct {
	UserTag   string `json:"usertag"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone_no  string `json:"phone_no"`
	Gender    string `json:"gender"`
	Dob       string `json:"dob"`
}

type PatientIdReq struct {
	Usertag string
}

type PatientIdResp struct {
	UserTag          string         `json:"usertag"`
	Name             string         `json:"name"`
	Phone_No         string         `json:"phone_no"`
	Gender           string         `json:"gender"`
	Dob              datatypes.Date `json:"dob"`
	Reason           string         `json:"reason"`
	Attending_Doctor string         `json:"attending_doctor"`
	File_URL         string         `json:"file_url"`
	Status           string         `json:"status"`
}

type Pharmacy struct {
	PharmacyID   string `json:"pharmacy_id"`
	PharmacyName string `json:"pharmacy_name"`
	Address      string `json:"address"`
	Country      string `json:"country"`
	State        string `json:"state"`
	About        string `json:"about"`
	Picture_url  string `json:"picture_url"`
}

type Hospital struct {
	HospitalID   string `json:"hospital_id"`
	HospitalName string `json:"hospital_name"`
	Address      string `json:"address"`
	Country      string `json:"country"`
	State        string `json:"state"`
	About        string `json:"about"`
	Picture_url  string `json:"picture_url"`
}

type Inventory struct {
	ProductID         string  `json:"product_id"`
	ProductName       string  `json:"product_name"`
	Milligrams        string  `json:"milligrams"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	Product_image_url string  `json:"product_image_url"`
}

type Orders struct {
	OrderID  string `json:"order_id"`
	UserTag  string `json:"usertag"`
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

type TestCentre struct {
	CentreID      string         `json:"centre_id"`
	CentreName    string         `json:"centre_name"`
	Address       string         `json:"address"`
	Country       string         `json:"country"`
	State         string         `json:"state"`
	DailyCapacity int            `json:"daily_capacity"`
	About         string         `json:"about"`
	Availability  datatypes.JSON `json:"availability"`
	TestType      string         `json:"test_type"`
	Price         float64        `json:"price"`
}

type Reviews struct {
	ReviewID  string `json:"review_id"`
	UserTag   string `json:"usertag"`
	DoctorTag string `json:"doctortag"`
	Review    string `json:"review"`
	Rating    int    `json:"rating"`
	Status    string `json:"status"`
}

type Getreviews struct {
	Status string `json:"status"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type AdminProfile struct {
	AdminTag      string `json:"admintag"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	ProfilePicURL string `json:"profile_pic"`
}

type GetDataReq struct {
	Status string
	Page   int
	Limit  int
	Search string
}
