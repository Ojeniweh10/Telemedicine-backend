package models

import "time"

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

type AddCart struct {
	Usertag   string
	ProductID int
	Quantity  int
}
