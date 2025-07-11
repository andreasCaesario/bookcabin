package usecase

import (
	"bookcabin-test/backend/internal/domain"
	"bookcabin-test/backend/internal/repository"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type VoucherUsecaseInterface interface {
	CheckVoucher(flightNumber, date string) (bool, error)
	GenerateVoucher(name, crewID, flightNumber, date, aircraft string) ([]string, error)
}

type voucherUsecase struct {
	repo repository.VoucherRepository
}

func NewVoucherUsecase(repo repository.VoucherRepository) VoucherUsecaseInterface {
	return &voucherUsecase{repo}
}

// CheckVoucher checks if a voucher exists for the given flight number and date
func (u *voucherUsecase) CheckVoucher(flightNumber, date string) (bool, error) {
	// find voucher by flight number and date
	voucher, err := u.repo.FindByFlightAndDate(flightNumber, date)
	if err != nil {
		return false, nil // not found is not an error
	}
	return voucher != nil, nil
}

// GenerateVoucher creates a new voucher with random seats for the specified flight
func (u *voucherUsecase) GenerateVoucher(name, crewID, flightNumber, date, aircraft string) ([]string, error) {
	// Validate aircraft type
	aircraftConf, exists := domain.AircraftConfigs[aircraft]
	if !exists {
		return nil, errors.New("unknown or invalid aircraft type")
	}

	// generate random seats
	seats := getRandomSeats(aircraftConf.Rows, aircraftConf.SeatLetters, 3)

	// Create a new voucher data structure
	newVoucher := &domain.Voucher{
		Name:         name,
		CrewID:       crewID,
		FlightNumber: flightNumber,
		FlightDate:   date,
		Aircraft:     aircraft,
		Seats:        joinSeats(seats),
	}

	// Save the voucher to the repository
	err := u.repo.Create(newVoucher)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

// getRandomSeats generates a list of random seats based on the aircraft configuration
func getRandomSeats(rows []int, seatLetters []string, seatCount int) []string {
	lenRows := len(rows)
	lenSeatLetters := len(seatLetters)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	seats := make(map[string]bool)
	var result []string
	for len(result) < seatCount {
		num := rows[r.Intn(lenRows)]
		letter := seatLetters[r.Intn(lenSeatLetters)]
		seat := fmt.Sprintf("%d%s", num, letter)
		if _, exists := seats[seat]; !exists {
			seats[seat] = true
			result = append(result, seat)
		}
	}
	return result
}

// joinSeats converts a slice of seat strings into a comma-separated string
func joinSeats(seats []string) string {
	return strings.Join(seats, ",")
}
