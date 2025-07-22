package handler

import (
	"bookcabin-test/backend/internal/usecase"
	"bookcabin-test/backend/internal/utils"
	"fmt"
	"net/http"
	"strconv"
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

// ReGenerateRequest defines the structure of the request for regenerating a voucher
type ReGenerateRequest struct {
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
	SeatIndex    string `json:"seatIndex"` // Index of the seat to regenerate (1, 2, or 3)
}

// CheckResponse defines the structure of the response for voucher existence check
type CheckResponse struct {
	Exists bool     `json:"exists"`
	Seats  []string `json:"seats"`
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
	r.POST("/api/re-generate", h.ReGenerate)
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
	isExists := false
	var seatData []string
	existsData, _ := h.Usecase.CheckVoucher(req.FlightNumber, req.Date)
	if existsData != nil {
		isExists = true
		seatData = append(seatData, existsData.Seat1, existsData.Seat2, existsData.Seat3)
	}

	c.JSON(http.StatusOK, CheckResponse{Exists: isExists, Seats: seatData})
}

// Generate handles the request to generate a voucherseatData
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

func (h *VoucherHandler) ReGenerate(c *gin.Context) {
	var req ReGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Simple input validation: check not empty
	if strings.TrimSpace(req.FlightNumber) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.SeatIndex) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flight number, date, and seat index are required."})
		return
	}
	// Validate seat index
	if req.SeatIndex != "1" && req.SeatIndex != "2" && req.SeatIndex != "3" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seat index. Must be 1, 2, or 3."})
		return
	}

	seatIndex, err := strconv.Atoi(req.SeatIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seat index."})
		return
	}

	Seats, err := h.Usecase.ReGenerateVoucher(req.FlightNumber, req.Date, seatIndex)
	if err != nil {
		// If voucher not found, return empty seats
		if strings.Contains(err.Error(), "record not found") {
			c.JSON(http.StatusOK, GenerateResponse{Success: true, Seats: []string{}})
			return
		}
		c.JSON(http.StatusBadRequest, GenerateResponse{Success: false, Error: err.Error()})
		return
	}

	fmt.Println(Seats)

	// If successful, return the regenerated seats
	c.JSON(http.StatusOK, GenerateResponse{Success: true, Seats: Seats})
}
