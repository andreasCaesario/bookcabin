package handler

import (
	"bookcabin-test/backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AircraftHandler struct {
	Usecase usecase.AircraftUsecaseInterface
}

func NewAircraftHandler(u usecase.AircraftUsecaseInterface) *AircraftHandler {
	return &AircraftHandler{Usecase: u}
}

// AircraftResponse defines the structure of aircraft type response
type AircraftResponse struct {
	Type string `json:"type"`
}

// RegisterRoutes sets up the routes for the AircraftHandler
func (h *AircraftHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/aircraft-list", h.getAircraftList)
}

func (h *AircraftHandler) getAircraftList(c *gin.Context) {
	var result []AircraftResponse
	for _, typ := range h.Usecase.GetAircraftList() {
		result = append(result, AircraftResponse{
			Type: typ,
		})
	}
	c.JSON(http.StatusOK, result)
}
