package domain

import (
	"time"
)

// Voucher defines the structure of a voucher in the system
type Voucher struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	CrewID       string
	FlightNumber string `gorm:"uniqueIndex:idx_flight_date"`
	FlightDate   string `gorm:"uniqueIndex:idx_flight_date"`
	Aircraft     string
	Seats        string    // comma separated
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
