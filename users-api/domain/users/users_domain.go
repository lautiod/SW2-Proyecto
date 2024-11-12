package users

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsAdmin bool   `json:"isadmin"`
	Email   string `json:"email"`
}

type Login_Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
