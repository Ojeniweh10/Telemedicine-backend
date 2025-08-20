package models

type Verifyemail struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

type Verifyemailresp struct {
	Usertag string `json:"usertag"`
}
