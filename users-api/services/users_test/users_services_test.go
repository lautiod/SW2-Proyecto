package users_test

import (
	"testing"
	domain "users-api/domain/users"
	repositories "users-api/repositories/users_test"
	"users-api/services/users"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	// Inicializar los mocks y el servicio para cada test
	mainRepo := repositories.NewMock()
	cacheRepo := repositories.NewMock()
	memcachedRepo := repositories.NewMock()
	userService := users.NewService(mainRepo, cacheRepo, memcachedRepo)

	t.Run("Create User", func(t *testing.T) {
		// Preparar datos de prueba
		newUser := domain.User{
			Email:    "test@test.com",
			Password: "password123",
			IsAdmin:  false,
		}

		// Ejecutar el método a probar
		id, err := userService.Create(newUser)

		// Verificar resultados
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("Login Success", func(t *testing.T) {
		// Preparar datos de prueba
		loginRequest := domain.Login_Request{
			Email:    "test@test.com",
			Password: "password123",
		}

		// Ejecutar el método a probar
		response, token, err := userService.Login(loginRequest)

		// Verificar resultados
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, loginRequest.Email, response.Email)
	})

	t.Run("Get User By ID", func(t *testing.T) {
		// Primero crear un usuario
		newUser := domain.User{
			Email:    "test@test.com",
			Password: "password123",
			IsAdmin:  false,
		}
		id, _ := userService.Create(newUser)

		// Obtener el usuario creado
		user, err := userService.GetByID(id)

		// Verificar resultados
		assert.NoError(t, err)
		assert.Equal(t, newUser.Email, user.Email)
		assert.Equal(t, id, user.ID)
	})

	t.Run("Get All Users", func(t *testing.T) {
		// Crear algunos usuarios de prueba
		users := []domain.User{
			{Email: "test1@test.com", Password: "pass1", IsAdmin: false},
			{Email: "test2@test.com", Password: "pass2", IsAdmin: true},
		}

		for _, user := range users {
			_, err := userService.Create(user)
			assert.NoError(t, err)
		}

		// Obtener todos los usuarios
		result, err := userService.GetAll()

		// Verificar resultados
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("Update User", func(t *testing.T) {
		// Crear un nuevo usuario
		newUser := domain.User{
			Email:    "update@test.com",
			Password: "password123",
			IsAdmin:  false,
		}
		id, err := userService.Create(newUser)
		assert.NoError(t, err)

		// Verificar que se creó correctamente
		createdUser, err := userService.GetByID(id)
		assert.NoError(t, err)
		assert.False(t, createdUser.IsAdmin)

		// Actualizar el usuario
		updatedUser := domain.User{
			ID:       id,
			Email:    "updated@test.com",
			Password: "newpassword",
			IsAdmin:  true,
		}

		// Ejecutar actualización
		err = userService.Update(updatedUser)
		assert.NoError(t, err)

		// Verificar la actualización
		result, err := userService.GetByID(id)
		assert.NoError(t, err)
		assert.Equal(t, "updated@test.com", result.Email)
		assert.True(t, result.IsAdmin, "IsAdmin debería ser true después de la actualización")
	})

	t.Run("Delete User", func(t *testing.T) {
		// Primero crear un usuario
		newUser := domain.User{
			Email:    "delete@test.com",
			Password: "password123",
			IsAdmin:  false,
		}
		id, _ := userService.Create(newUser)

		// Ejecutar eliminación
		err := userService.Delete(id)
		assert.NoError(t, err)

		// Intentar obtener el usuario eliminado
		_, err = userService.GetByID(id)
		assert.Error(t, err)
	})

	t.Run("Login Failed - Wrong Password", func(t *testing.T) {
		// Crear usuario
		newUser := domain.User{
			Email:    "login@test.com",
			Password: "correctpassword",
			IsAdmin:  false,
		}
		_, _ = userService.Create(newUser)

		// Intentar login con contraseña incorrecta
		loginRequest := domain.Login_Request{
			Email:    "login@test.com",
			Password: "wrongpassword",
		}

		response, token, err := userService.Login(loginRequest)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Empty(t, response)
	})

	t.Run("Login Failed - User Not Found", func(t *testing.T) {
		loginRequest := domain.Login_Request{
			Email:    "nonexistent@test.com",
			Password: "password123",
		}

		response, token, err := userService.Login(loginRequest)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Empty(t, response)
	})
}
