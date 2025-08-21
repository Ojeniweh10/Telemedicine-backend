package models

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
