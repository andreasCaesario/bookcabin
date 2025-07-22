package repository

import (
	"bookcabin-test/backend/internal/domain"

	"gorm.io/gorm"
)

type VoucherRepository interface {
	FindByFlightAndDate(flightNumber, date string) (*domain.Voucher, error)
	Create(voucher *domain.Voucher) error
}

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db}
}

// FindByFlightAndDate retrieves a voucher by flight number and date
func (r *voucherRepository) FindByFlightAndDate(flightNumber, date string) (*domain.Voucher, error) {
	var voucher domain.Voucher
	err := r.db.Where("flight_number = ? AND flight_date = ?", flightNumber, date).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

// Create saves a new voucher to the database
func (r *voucherRepository) Create(voucher *domain.Voucher) error {
	return r.db.Save(voucher).Error
}
