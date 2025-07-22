package usecase

import (
	"bookcabin-test/backend/internal/domain"
	"bookcabin-test/backend/internal/repository"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type VoucherUsecaseInterface interface {
	CheckVoucher(flightNumber, date string) (*domain.Voucher, error)
	GenerateVoucher(name, crewID, flightNumber, date, aircraft string) ([]string, error)
	ReGenerateVoucher(flightNumber, flightDate string, seatIndex int) ([]string, error)
}

type voucherUsecase struct {
	repo repository.VoucherRepository
}

func NewVoucherUsecase(repo repository.VoucherRepository) VoucherUsecaseInterface {
	return &voucherUsecase{repo}
}

// CheckVoucher checks if a voucher exists for the given flight number and date
func (u *voucherUsecase) CheckVoucher(flightNumber, date string) (*domain.Voucher, error) {
	// find voucher by flight number and date
	voucher, err := u.repo.FindByFlightAndDate(flightNumber, date)
	if err != nil {
		return nil, nil // not found is not an error
	}
	return voucher, nil
}

// GenerateVoucher creates a new voucher with random seats for the specified flight
func (u *voucherUsecase) GenerateVoucher(name, crewID, flightNumber, date, aircraft string) ([]string, error) {
	// Validate aircraft type
	aircraftConf, exists := domain.AircraftConfigs[aircraft]
	if !exists {
		return nil, errors.New("unknown or invalid aircraft type")
	}

	// generate random seats
	seats := getRandomSeats(aircraftConf.Rows, aircraftConf.SeatLetters, 3, nil)

	// Create a new voucher data structure
	newVoucher := &domain.Voucher{
		Name:         name,
		CrewID:       crewID,
		FlightNumber: flightNumber,
		FlightDate:   date,
		Aircraft:     aircraft,
		Seat1:        seats[0],
		Seat2:        seats[1],
		Seat3:        seats[2],
	}

	// Save the voucher to the repository
	err := u.repo.Create(newVoucher)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

func (u *voucherUsecase) ReGenerateVoucher(flightNumber, flightDate string, seatIndex int) ([]string, error) {
	// Get existing voucher
	voucher, err := u.repo.FindByFlightAndDate(flightNumber, flightDate)
	if err != nil {
		return []string{}, err // not found is not an error
	}

	// Get aircraft type and seats confoguration
	aircraftConf, exists := domain.AircraftConfigs[voucher.Aircraft]
	if !exists {
		return []string{}, errors.New("unknown or invalid aircraft type")
	}

	existingSeats := map[string]bool{
		voucher.Seat1: true,
		voucher.Seat2: true,
		voucher.Seat3: true,
	}
	// generate random seats
	seats := getRandomSeats(aircraftConf.Rows, aircraftConf.SeatLetters, 1, existingSeats)

	if seatIndex == 1 {
		voucher.Seat1 = seats[0]
	} else if seatIndex == 2 {
		voucher.Seat2 = seats[0]
	} else if seatIndex == 3 {
		voucher.Seat3 = seats[0]
	} else {
		return []string{}, errors.New("invalid seat index")
	}

	// Save the voucher to the repository
	err = u.repo.Create(voucher)
	if err != nil {
		return nil, err
	}
	return append([]string{}, voucher.Seat1, voucher.Seat2, voucher.Seat3), nil
}

// getRandomSeats generates a list of random seats based on the aircraft configuration
func getRandomSeats(rows []int, seatLetters []string, seatCount int, existingSeatConfig map[string]bool) []string {
	lenRows := len(rows)
	lenSeatLetters := len(seatLetters)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	seats := make(map[string]bool)
	if seats == nil || len(seats) == 0 {
		seats = make(map[string]bool)
	}

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
