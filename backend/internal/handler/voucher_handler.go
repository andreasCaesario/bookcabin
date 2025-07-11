package handler

import (
	"bookcabin-test/backend/internal/usecase"
	"bookcabin-test/backend/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	Usecase usecase.VoucherUsecaseInterface
}

func NewVoucherHandler(u usecase.VoucherUsecaseInterface) *VoucherHandler {
	return &VoucherHandler{Usecase: u}
}

// CheckRequest defines the structure of the request for checking voucher existence
type CheckRequest struct {
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
}

// CheckResponse defines the structure of the response for voucher existence check
type CheckResponse struct {
	Exists bool `json:"exists"`
}

// GenerateRequest defines the structure of the request for generating a voucher
type GenerateRequest struct {
	Name         string `json:"name"`
	CrewID       string `json:"id"`
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
	Aircraft     string `json:"aircraft"`
}

// GenerateResponse defines the structure of the response for voucher generation
type GenerateResponse struct {
	Success bool     `json:"success"`
	Seats   []string `json:"seats"`
	Error   string   `json:"error,omitempty"`
}

// RegisterRoutes sets up the routes for the VoucherHandler
func (h *VoucherHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/check", h.Check)
	r.POST("/api/generate", h.Generate)
}

// Check handles the request to check if a voucher exists for a given flight number and date
func (h *VoucherHandler) Check(c *gin.Context) {
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Simple input validation: check not empty
	if strings.TrimSpace(req.FlightNumber) == "" || strings.TrimSpace(req.Date) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flight number and date are required."})
		return
	}
	exists, _ := h.Usecase.CheckVoucher(req.FlightNumber, req.Date)
	c.JSON(http.StatusOK, CheckResponse{Exists: exists})
}

// Generate handles the request to generate a voucher
func (h *VoucherHandler) Generate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Simple input validation: check not empty
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.CrewID) == "" ||
		strings.TrimSpace(req.FlightNumber) == "" || strings.TrimSpace(req.Date) == "" ||
		strings.TrimSpace(req.Aircraft) == "" {
		c.JSON(http.StatusBadRequest, GenerateResponse{Success: false, Error: "All fields are required."})
		return
	}
	// Validate date is not in the past (expects DD - MM - YY)
	if !utils.ValidateDateNotPast(req.Date) {
		c.JSON(http.StatusBadRequest, GenerateResponse{Success: false, Error: "Date cannot be in the past."})
		return
	}
	seats, err := h.Usecase.GenerateVoucher(req.Name, req.CrewID, req.FlightNumber, req.Date, req.Aircraft)
	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusBadRequest, GenerateResponse{Success: false, Error: "voucher already generated for this flight number and for the date you choose."})
			return
		}
		c.JSON(http.StatusBadRequest, GenerateResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, GenerateResponse{Success: true, Seats: seats})
}
