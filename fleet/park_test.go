package fleet

import (
	"errors"
	"math"
	"testing"
)

var safe = Condition{BrakesOK: true, FrameOK: true}

func TestInspectionBeforeAcceptance(t *testing.T) {
	park := NewPark(SafetyServiceCenter{})
	for _, tc := range []struct {
		name      string
		condition Condition
		accepted  bool
	}{
		{"safe", safe, true},
		{"bad brakes", Condition{FrameOK: true}, false},
		{"bad frame", Condition{BrakesOK: true}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			vehicle, err := NewScooter(tc.name, 8, tc.condition, 1.5)
			if err != nil {
				t.Fatal(err)
			}
			err = park.Add(vehicle)
			if tc.accepted && err != nil || !tc.accepted && !errors.Is(err, ErrInspectionFailed) {
				t.Fatalf("Add() error = %v, accepted = %v", err, tc.accepted)
			}
		})
	}
	if summary := park.Summary(); summary.Vehicles != 1 || summary.TotalDailyKWh != 1.5 {
		t.Fatalf("rejected vehicles must not affect reports: %+v", summary)
	}
}

func TestReportsIncludeOnlyAcceptedAssets(t *testing.T) {
	park := NewPark(SafetyServiceCenter{})
	scooter, _ := NewScooter("S-1", 6, safe, 1.25)
	bicycle, _ := NewBicycle("B-1", 5, safe)
	ebike, _ := NewElectricBicycle("E-1", 10, safe, 0.75)
	helmet, _ := NewHelmet("H-1")
	dock, _ := NewChargingDock("D-1", 2)
	for _, asset := range []Asset{scooter, bicycle, ebike, helmet, dock} {
		if err := park.Add(asset); err != nil {
			t.Fatal(err)
		}
	}
	if got := park.Summary(); got != (Summary{Vehicles: 3, Things: 2, TotalDailyKWh: 4}) {
		t.Fatalf("Summary() = %+v", got)
	}
	beginner := park.BeginnerVehicles()
	if len(beginner) != 2 || beginner[0].InventoryNumber() != "S-1" || beginner[1].InventoryNumber() != "E-1" {
		t.Fatalf("BeginnerVehicles() = %v", beginner)
	}
	inventory := park.Inventory()
	if len(inventory) != 5 || inventory[4].Name() != "Charging dock" {
		t.Fatalf("Inventory() = %v", inventory)
	}
	inventory[0] = helmet
	if park.Inventory()[0].InventoryNumber() != "S-1" {
		t.Fatal("Inventory() exposed internal slice")
	}
}

func TestDuplicateNumberAcrossKinds(t *testing.T) {
	park := NewPark(SafetyServiceCenter{})
	helm, _ := NewHelmet("S-1")
	scooter, _ := NewScooter("S-1", 7, safe, 1)
	if err := park.Add(helm); err != nil {
		t.Fatal(err)
	}
	if err := park.Add(scooter); !errors.Is(err, ErrDuplicateInventory) {
		t.Fatalf("Add(duplicate) = %v", err)
	}
	if got := park.Summary(); got.Things != 1 || got.Vehicles != 0 {
		t.Fatalf("duplicate changed inventory: %+v", got)
	}
}

func TestConstructorsRejectInvalidInputs(t *testing.T) {
	if _, err := NewBicycle("  ", 6, safe); !errors.Is(err, ErrInvalidInventoryNumber) {
		t.Fatalf("empty inventory number: %v", err)
	}
	for _, ease := range []int{0, 11} {
		if _, err := NewBicycle("B-1", ease, safe); !errors.Is(err, ErrInvalidEase) {
			t.Fatalf("ease %d: %v", ease, err)
		}
	}
	for _, energy := range []float64{-1, math.NaN(), math.Inf(1)} {
		if _, err := NewScooter("S-1", 6, safe, energy); !errors.Is(err, ErrInvalidEnergy) {
			t.Fatalf("scooter energy %v: %v", energy, err)
		}
		if _, err := NewElectricBicycle("E-1", 6, safe, energy); !errors.Is(err, ErrInvalidEnergy) {
			t.Fatalf("ebike energy %v: %v", energy, err)
		}
		if _, err := NewChargingDock("D-1", energy); !errors.Is(err, ErrInvalidEnergy) {
			t.Fatalf("dock energy %v: %v", energy, err)
		}
	}
	if _, err := NewHelmet(""); !errors.Is(err, ErrInvalidInventoryNumber) {
		t.Fatalf("empty helmet inventory number: %v", err)
	}
}

func TestInvalidAsset(t *testing.T) {
	park := NewPark(SafetyServiceCenter{})
	if err := park.Add(nil); !errors.Is(err, ErrInvalidAsset) {
		t.Fatalf("Add(nil) = %v", err)
	}
}
