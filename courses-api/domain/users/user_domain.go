package users

// User representa el modelo de usuario que viene de la API de users
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	IsAdmin   bool   `json:"is_admin"`
}
