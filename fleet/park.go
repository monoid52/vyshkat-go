package fleet

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInspectionFailed   = errors.New("vehicle failed inspection")
	ErrDuplicateInventory = errors.New("inventory number already in use")
	ErrInvalidAsset       = errors.New("asset and inventory number are required")
)

// ServiceCenter decides whether a new vehicle may enter the active fleet.
type ServiceCenter interface {
	Inspect(Vehicle) error
}

// SafetyServiceCenter checks the brakes and frame before acceptance.
type SafetyServiceCenter struct{}

func (SafetyServiceCenter) Inspect(v Vehicle) error {
	condition := v.Condition()
	if !condition.BrakesOK || !condition.FrameOK {
		return ErrInspectionFailed
	}
	return nil
}

type Park struct {
	center  ServiceCenter
	assets  []Asset
	numbers map[string]struct{}
}

func NewPark(center ServiceCenter) *Park {
	return &Park{center: center, numbers: make(map[string]struct{})}
}

// Add is the only route into the park: every vehicle is inspected, while
// non-vehicle assets can be inventoried without a vehicle inspection.
func (p *Park) Add(asset Asset) error {
	if asset == nil || strings.TrimSpace(asset.InventoryNumber()) == "" {
		return ErrInvalidAsset
	}
	number := asset.InventoryNumber()
	if _, exists := p.numbers[number]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateInventory, number)
	}
	if vehicle, ok := asset.(Vehicle); ok {
		if err := p.center.Inspect(vehicle); err != nil {
			return fmt.Errorf("inspection of %s: %w", number, err)
		}
	}
	p.assets = append(p.assets, asset)
	p.numbers[number] = struct{}{}
	return nil
}

type Summary struct {
	Vehicles      int
	Things        int
	TotalDailyKWh float64
}

func (p *Park) Summary() Summary {
	var result Summary
	for _, asset := range p.assets {
		if _, ok := asset.(Vehicle); ok {
			result.Vehicles++
		} else {
			result.Things++
		}
		if consumer, ok := asset.(EnergyConsumer); ok {
			result.TotalDailyKWh += consumer.DailyEnergyKWh()
		}
	}
	return result
}

func (p *Park) BeginnerVehicles() []Vehicle {
	var result []Vehicle
	for _, asset := range p.assets {
		if vehicle, ok := asset.(Vehicle); ok && vehicle.EaseOfUse() >= 6 {
			result = append(result, vehicle)
		}
	}
	return result
}

// Inventory returns a copy, so callers cannot mutate the park's asset list.
func (p *Park) Inventory() []Asset {
	return append([]Asset(nil), p.assets...)
}
