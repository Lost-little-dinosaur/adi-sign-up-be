package signUp

type AddSignUpRequest struct {
	TeamName     string   `json:"teamName"`
	IsHDU        bool     `json:"isHDU"`
	School       string   `json:"school"`
	MemberArr    []Member `json:"memberArr"`
	CaptchaId    string   `json:"captchaId"`
	CaptchaValue string   `json:"captchaValue"`
}

type Member struct {
	Phone    string `json:"phone"`
	QQ       string `json:"qq"`
	Name     string `json:"name"`
	IDNumber string `json:"idNumber"`
	//BankCardNumber string `json:"bankCardNumber"`
	//BankName       string `json:"bankName"`
	HDUID string `json:"hduId"`
}

type SignUpForm struct {
	TeamName  string `json:"teamName"`
	IsHDU     bool   `json:"isHDU"`
	School    string `json:"school"`
	MemberArr []Member
}
type GetSignUpResponse struct {
	Number        int `json:"number"`
	SignUpFormArr []SignUpForm
}
