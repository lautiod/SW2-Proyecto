package middleware

import (
	//"backend/clients"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	dao "users-api/dao/users"
	repositories "users-api/repositories/users"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Validar el token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Obtener el ID del usuario desde claims["sub"]
		userID, ok := claims["sub"].(float64) // JWT usa float64 por defecto
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user dao.User
		// Usar el método GetByID para obtener el usuario
		user, err := repositories.MySQLInstance.GetByID(int64(userID))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Validar que el usuario fue encontrado
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Configurar el contexto con datos del usuario
		c.Set("user", user)
		c.Set("userID", user.ID)
		//c.Set("isAdmin", user.IsAdmin)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
