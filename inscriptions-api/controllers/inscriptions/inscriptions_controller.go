package courses

import (
	"context"
	"fmt"
	inscriptionsDomain "inscriptions-api/domain/inscriptions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	EnrollUser(ctx context.Context, inscription inscriptionsDomain.Inscription) (string, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) EnrollUser(ctx *gin.Context) {
	// Parse inscription
	var inscription inscriptionsDomain.Inscription
	if err := ctx.ShouldBindJSON(&inscription); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Create hotel
	id, err := controller.service.EnrollUser(ctx.Request.Context(), inscription)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error enrrolling user: %s", err.Error()),
		})
		return
	}

	// Send ID
	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}
