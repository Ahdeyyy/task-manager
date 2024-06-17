package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Session  string
	Sessions []UserSession
}

type UserSession struct {
	SID    string `json:"sid"`
	IP     string `json:"ip"`
	Login  string `json:"login"`
	Expiry string `json:"expiry"`
	UA     string `json:"ua"`
}
