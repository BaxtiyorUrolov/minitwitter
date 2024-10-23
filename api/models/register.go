package models

type VerifyCodeRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
