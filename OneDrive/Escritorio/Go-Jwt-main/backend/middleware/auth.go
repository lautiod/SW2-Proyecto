package middleware

import (
	"backend/clients"
	"backend/dao"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	//VALIDO QUE EL METODO DE PARSEO SEA EL CORRECTO Y DEVUELVO LA VAR SECRET CONVERTIDA A UN ARRAY DE BYTES
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user dao.User
		clients.Db.First(&user, claims["sub"])

		if user.UserID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)
		c.Set("userID", user.UserID)
		c.Set("isAdmin", user.IsAdmin)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

	}

}
