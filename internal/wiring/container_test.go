package wiring

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/monoid52/vyshkat-go/fleet"
	"github.com/monoid52/vyshkat-go/internal/console"
)

type stubCenter struct{ result error }

func (s stubCenter) Inspect(fleet.Vehicle) error { return s.result }

// The same park behavior changes only by replacing the container binding.
func TestContainerReplacesServiceCenter(t *testing.T) {
	unsafe := fleet.Condition{BrakesOK: false, FrameOK: false}
	for _, tc := range []struct {
		name     string
		decision error
		want     int
	}{
		{"approve unsafe vehicle", nil, 1},
		{"reject vehicle", fleet.ErrInspectionFailed, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			container, err := NewContainer(strings.NewReader(""), &bytes.Buffer{}, func() fleet.ServiceCenter {
				return stubCenter{result: tc.decision}
			})
			if err != nil {
				t.Fatal(err)
			}
			err = container.Invoke(func(park *fleet.Park) error {
				vehicle, err := fleet.NewBicycle("B-1", 7, unsafe)
				if err != nil {
					return err
				}
				err = park.Add(vehicle)
				if tc.decision == nil && err != nil || tc.decision != nil && !errors.Is(err, tc.decision) {
					t.Errorf("Add() error = %v", err)
				}
				if got := park.Summary().Vehicles; got != tc.want {
					t.Errorf("accepted vehicles = %d, want %d", got, tc.want)
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestContainerConstructsConsoleApp(t *testing.T) {
	var out bytes.Buffer
	c, err := NewContainer(strings.NewReader("0\n"), &out, func() fleet.ServiceCenter {
		return fleet.SafetyServiceCenter{}
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Invoke(func(app *console.App) error { return app.Run() }); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "VyshKat fleet inventory") {
		t.Fatalf("console was not constructed: %s", out.String())
	}
}
