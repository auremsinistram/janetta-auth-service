package account

type SignUpReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type SignUpConfirmReq struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
}

type SignInReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type SignInConfirmReq struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
}

type LogoutReq struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type PasswordForgotReq struct {
	Phone string `json:"phone"`
}

type PasswordResetReq struct {
	Code     string `json:"code"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type PasswordChangeReq struct {
	Password string `json:"password"`
}

type PasswordChangeConfirmReq struct {
	Code string `json:"code"`
}
