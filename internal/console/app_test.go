package console

import (
	"bytes"
	"strings"
	"testing"

	"github.com/monoid52/vyshkat-go/fleet"
)

func TestInteractiveFlow(t *testing.T) {
	// An unsafe scooter is rejected, while a safe bicycle, helmet and dock are inventoried.
	input := strings.Join([]string{
		"1", "S-1", "9", "no", "yes", "1.5",
		"2", "B-1", "6", "yes", "yes",
		"4", "H-1",
		"5", "D-1", "2.5",
		"4", "B-1", // duplicate across kinds
		"6", "7", "8", "0",
	}, "\n") + "\n"
	var out bytes.Buffer
	app := New(fleet.NewPark(fleet.SafetyServiceCenter{}), strings.NewReader(input), &out)
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Not added: inspection", "Added: Bicycle [B-1]", "Not added: inventory number already in use",
		"Vehicles: 1 | Other assets: 2 | Total assets: 3 | Daily energy: 2.50 kWh",
		"Bicycle [B-1], ease 6", "Helmet [H-1]", "Charging dock [D-1]",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "Accepted: Electric scooter") {
		t.Fatal("rejected scooter appeared as accepted")
	}
}

func TestInputErrorsDoNotAddAssets(t *testing.T) {
	input := "8\n7\n6\n9\n2\nB-1\ninvalid\n4\nH-1\n0\n"
	var out bytes.Buffer
	park := fleet.NewPark(fleet.SafetyServiceCenter{})
	if err := New(park, strings.NewReader(input), &out).Run(); err != nil {
		t.Fatal(err)
	}
	if got := park.Summary(); got.Vehicles != 0 || got.Things != 1 {
		t.Fatalf("invalid input changed park: %+v", got)
	}
	for _, want := range []string{"Inventory:", "Beginner-friendly vehicles", "Unknown option", "Not added: strconv.Atoi"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestDemo(t *testing.T) {
	var out bytes.Buffer
	app := New(fleet.NewPark(fleet.SafetyServiceCenter{}), strings.NewReader(""), &out)
	if err := app.RunDemo(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Accepted: Electric scooter", "Rejected: Electric scooter [S-002]",
		"Rejected: Helmet [S-001]", "Vehicles: 3 | Other assets: 2 | Total assets: 5 | Daily energy: 4.30 kWh",
		"Electric scooter [S-001], ease 6", "Electric bicycle [E-001], ease 8",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("demo missing %q:\n%s", want, out.String())
		}
	}
}
