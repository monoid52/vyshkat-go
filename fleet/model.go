package fleet

import (
	"errors"
	"math"
	"strings"
)

// Asset is anything the company tracks by inventory number.
type Asset interface {
	InventoryNumber() string
	Name() string
}

// EnergyConsumer reports daily electricity consumption in kWh.
type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

// Vehicle is rentable equipment. Inspection and beginner reports use this interface.
type Vehicle interface {
	Asset
	EaseOfUse() int
	Condition() Condition
}

// Condition is the observable input to a safety inspection.
type Condition struct {
	BrakesOK bool
	FrameOK  bool
}

var (
	ErrInvalidInventoryNumber = errors.New("inventory number cannot be empty")
	ErrInvalidEase            = errors.New("ease of use must be between 1 and 10")
	ErrInvalidEnergy          = errors.New("daily energy must be finite and non-negative")
)

type vehicleDetails struct {
	number    string
	ease      int
	condition Condition
}

func newVehicleDetails(number string, ease int, condition Condition) (vehicleDetails, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return vehicleDetails{}, ErrInvalidInventoryNumber
	}
	if ease < 1 || ease > 10 {
		return vehicleDetails{}, ErrInvalidEase
	}
	return vehicleDetails{number: number, ease: ease, condition: condition}, nil
}

func (v vehicleDetails) InventoryNumber() string { return v.number }
func (v vehicleDetails) EaseOfUse() int          { return v.ease }
func (v vehicleDetails) Condition() Condition    { return v.condition }

func validateEnergy(kwh float64) error {
	if kwh < 0 || math.IsNaN(kwh) || math.IsInf(kwh, 0) {
		return ErrInvalidEnergy
	}
	return nil
}

// Scooter and ElectricBicycle consume electricity; Bicycle does not.
type Scooter struct {
	vehicleDetails
	dailyKWh float64
}

func NewScooter(number string, ease int, condition Condition, dailyKWh float64) (*Scooter, error) {
	details, err := newVehicleDetails(number, ease, condition)
	if err != nil {
		return nil, err
	}
	if err := validateEnergy(dailyKWh); err != nil {
		return nil, err
	}
	return &Scooter{details, dailyKWh}, nil
}

func (*Scooter) Name() string              { return "Electric scooter" }
func (s *Scooter) DailyEnergyKWh() float64 { return s.dailyKWh }

type Bicycle struct{ vehicleDetails }

func NewBicycle(number string, ease int, condition Condition) (*Bicycle, error) {
	details, err := newVehicleDetails(number, ease, condition)
	if err != nil {
		return nil, err
	}
	return &Bicycle{details}, nil
}

func (*Bicycle) Name() string { return "Bicycle" }

type ElectricBicycle struct {
	vehicleDetails
	dailyKWh float64
}

func NewElectricBicycle(number string, ease int, condition Condition, dailyKWh float64) (*ElectricBicycle, error) {
	details, err := newVehicleDetails(number, ease, condition)
	if err != nil {
		return nil, err
	}
	if err := validateEnergy(dailyKWh); err != nil {
		return nil, err
	}
	return &ElectricBicycle{details, dailyKWh}, nil
}

func (*ElectricBicycle) Name() string              { return "Electric bicycle" }
func (b *ElectricBicycle) DailyEnergyKWh() float64 { return b.dailyKWh }

// Helmet is inventoried without consuming electricity.
type Helmet struct{ number string }

func NewHelmet(number string) (*Helmet, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return nil, ErrInvalidInventoryNumber
	}
	return &Helmet{number}, nil
}

func (*Helmet) Name() string              { return "Helmet" }
func (h *Helmet) InventoryNumber() string { return h.number }

// ChargingDock is a non-rentable asset that consumes electricity.
type ChargingDock struct {
	number   string
	dailyKWh float64
}

func NewChargingDock(number string, dailyKWh float64) (*ChargingDock, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return nil, ErrInvalidInventoryNumber
	}
	if err := validateEnergy(dailyKWh); err != nil {
		return nil, err
	}
	return &ChargingDock{number, dailyKWh}, nil
}

func (*ChargingDock) Name() string              { return "Charging dock" }
func (d *ChargingDock) InventoryNumber() string { return d.number }
func (d *ChargingDock) DailyEnergyKWh() float64 { return d.dailyKWh }
