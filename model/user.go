package model

type UserInfo struct {
	Location string `json:"location"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}
