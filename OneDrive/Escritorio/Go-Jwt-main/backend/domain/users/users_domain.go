package users

type Signup_Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isadmin"`
}

type Login_Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsAdmin bool   `json:"isadmin"`
	Email   string `json:"email"`
}

type User struct {
	UserID   uint   `json:"userID"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}
