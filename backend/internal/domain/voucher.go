package domain

import (
	"time"
)

// Voucher defines the structure of a voucher in the system
type Voucher struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	CrewID       string
	FlightNumber string
	FlightDate   string
	Aircraft     string
	Seat1        string
	Seat2        string
	Seat3        string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
