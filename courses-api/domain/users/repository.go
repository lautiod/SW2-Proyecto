package users

import "context"

// ExternalRepository define los métodos para interactuar con la API de users
type ExternalRepository interface {
	GetUserByID(ctx context.Context, id string) (User, error)
	ValidateAdminUser(ctx context.Context, userID string) (bool, error)
	// Puedes agregar más métodos según necesites
}
