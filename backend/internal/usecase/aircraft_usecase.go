package usecase

import "bookcabin-test/backend/internal/domain"

type AircraftUsecaseInterface interface {
	GetAircraftList() []string
}

type aircraftUsecase struct {
}

func NewAircraftUsecase() AircraftUsecaseInterface {
	return &aircraftUsecase{}
}

// GetAircraftList returns a list of available aircraft types
func (u *aircraftUsecase) GetAircraftList() (aircraftList []string) {
	for typ, _ := range domain.AircraftConfigs {
		aircraftList = append(aircraftList, typ)
	}
	return aircraftList
}
