package console

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/monoid52/vyshkat-go/fleet"
)

type App struct {
	park *fleet.Park
	in   *bufio.Scanner
	out  io.Writer
}

func New(park *fleet.Park, in io.Reader, out io.Writer) *App {
	return &App{park: park, in: bufio.NewScanner(in), out: out}
}

func (a *App) Run() error {
	for {
		fmt.Fprintln(a.out, "\nVyshKat fleet inventory")
		fmt.Fprintln(a.out, "1) Electric scooter  2) Bicycle  3) Electric bicycle")
		fmt.Fprintln(a.out, "4) Helmet  5) Charging dock")
		fmt.Fprintln(a.out, "6) Summary  7) Beginner vehicles  8) Inventory  0) Exit")
		choice, err := a.prompt("Choose an option: ")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch choice {
		case "0":
			return nil
		case "1", "2", "3", "4", "5":
			asset, err := a.readAsset(choice)
			if err == nil {
				err = a.park.Add(asset)
			}
			if err != nil {
				fmt.Fprintf(a.out, "Not added: %v\n", err)
				if err == io.EOF {
					return nil
				}
				continue
			}
			fmt.Fprintf(a.out, "Added: %s [%s]\n", asset.Name(), asset.InventoryNumber())
		case "6":
			a.printSummary()
		case "7":
			a.printBeginners()
		case "8":
			a.printInventory()
		default:
			fmt.Fprintln(a.out, "Unknown option. Choose 0-8.")
		}
	}
}

func (a *App) prompt(label string) (string, error) {
	if _, err := fmt.Fprint(a.out, label); err != nil {
		return "", err
	}
	if !a.in.Scan() {
		if err := a.in.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(a.in.Text()), nil
}

func (a *App) readInt(label string) (int, error) {
	text, err := a.prompt(label)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(text)
}

func (a *App) readEnergy() (float64, error) {
	text, err := a.prompt("Daily electricity use (kWh): ")
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(text, 64)
}

func (a *App) readYesNo(label string) (bool, error) {
	text, err := a.prompt(label)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(text) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("enter yes or no, got %q", text)
	}
}

func (a *App) readAsset(choice string) (fleet.Asset, error) {
	number, err := a.prompt("Inventory number: ")
	if err != nil {
		return nil, err
	}
	switch choice {
	case "4":
		return fleet.NewHelmet(number)
	case "5":
		energy, err := a.readEnergy()
		if err != nil {
			return nil, err
		}
		return fleet.NewChargingDock(number, energy)
	}
	ease, err := a.readInt("Ease of use (1-10): ")
	if err != nil {
		return nil, err
	}
	brakes, err := a.readYesNo("Brakes OK? (yes/no): ")
	if err != nil {
		return nil, err
	}
	frame, err := a.readYesNo("Frame OK? (yes/no): ")
	if err != nil {
		return nil, err
	}
	condition := fleet.Condition{BrakesOK: brakes, FrameOK: frame}
	if choice == "2" {
		return fleet.NewBicycle(number, ease, condition)
	}
	energy, err := a.readEnergy()
	if err != nil {
		return nil, err
	}
	if choice == "1" {
		return fleet.NewScooter(number, ease, condition, energy)
	}
	return fleet.NewElectricBicycle(number, ease, condition, energy)
}

func (a *App) printSummary() {
	s := a.park.Summary()
	fmt.Fprintf(a.out, "Vehicles: %d | Other assets: %d | Total assets: %d | Daily energy: %.2f kWh\n",
		s.Vehicles, s.Things, s.Vehicles+s.Things, s.TotalDailyKWh)
}

func (a *App) printBeginners() {
	fmt.Fprintln(a.out, "Beginner-friendly vehicles (ease >= 6):")
	for _, vehicle := range a.park.BeginnerVehicles() {
		fmt.Fprintf(a.out, "- %s [%s], ease %d\n", vehicle.Name(), vehicle.InventoryNumber(), vehicle.EaseOfUse())
	}
}

func (a *App) printInventory() {
	fmt.Fprintln(a.out, "Inventory:")
	for _, asset := range a.park.Inventory() {
		fmt.Fprintf(a.out, "- %s [%s]\n", asset.Name(), asset.InventoryNumber())
	}
}

// RunDemo exercises acceptance, rejection, duplicate IDs and every report.
func (a *App) RunDemo() error {
	good := fleet.Condition{BrakesOK: true, FrameOK: true}
	bad := fleet.Condition{BrakesOK: false, FrameOK: true}
	scooter, err := fleet.NewScooter("S-001", 6, good, 1.2)
	if err != nil {
		return err
	}
	rejected, err := fleet.NewScooter("S-002", 9, bad, 0.8)
	if err != nil {
		return err
	}
	bicycle, err := fleet.NewBicycle("B-001", 5, good)
	if err != nil {
		return err
	}
	ebike, err := fleet.NewElectricBicycle("E-001", 8, good, 0.6)
	if err != nil {
		return err
	}
	helmet, err := fleet.NewHelmet("H-001")
	if err != nil {
		return err
	}
	dock, err := fleet.NewChargingDock("D-001", 2.5)
	if err != nil {
		return err
	}
	duplicate, err := fleet.NewHelmet("S-001")
	if err != nil {
		return err
	}

	fmt.Fprintln(a.out, "VyshKat demonstration")
	for _, asset := range []fleet.Asset{scooter, rejected, bicycle, ebike, helmet, dock, duplicate} {
		if err := a.park.Add(asset); err != nil {
			fmt.Fprintf(a.out, "Rejected: %s [%s]: %v\n", asset.Name(), asset.InventoryNumber(), err)
		} else {
			fmt.Fprintf(a.out, "Accepted: %s [%s]\n", asset.Name(), asset.InventoryNumber())
		}
	}
	a.printSummary()
	a.printBeginners()
	a.printInventory()
	return nil
}
