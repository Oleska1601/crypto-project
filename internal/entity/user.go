package entity

type User struct {
	ID           int    `json:"id,omitempty"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash,omitempty"`
	Salt         string `json:"salt,omitempty"`
	Secret       string `json:"secret,omitempty"`
}
