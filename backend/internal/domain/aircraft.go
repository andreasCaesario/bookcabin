package domain

import (
	"bookcabin-test/backend/internal/utils"
)

// AircraftConfig defines the seat configuration for an aircraft type
var AircraftConfigs = map[string]struct {
	Rows        []int
	SeatLetters []string
}{
	"ATR":            {Rows: utils.MakeRange(1, 18), SeatLetters: []string{"A", "C", "D", "F"}},
	"Airbus 320":     {Rows: utils.MakeRange(1, 32), SeatLetters: []string{"A", "B", "C", "D", "E", "F"}},
	"Boeing 737 Max": {Rows: utils.MakeRange(1, 32), SeatLetters: []string{"A", "B", "C", "D", "E", "F"}},
}
