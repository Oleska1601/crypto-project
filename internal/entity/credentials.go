package entity

type Credentials struct {
	Login    string `json:"login" example:"user"`
	Password string `json:"password" example:"pass"`
}
