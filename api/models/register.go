package models

type VerifyCodeRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
