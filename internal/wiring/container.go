package wiring

import (
	"io"

	"github.com/monoid52/vyshkat-go/fleet"
	"github.com/monoid52/vyshkat-go/internal/console"
	"go.uber.org/dig"
)

// NewContainer wires the app once. Tests can supply a fake service center.
func NewContainer(in io.Reader, out io.Writer, centerProvider func() fleet.ServiceCenter) (*dig.Container, error) {
	c := dig.New()
	if err := c.Provide(centerProvider); err != nil {
		return nil, err
	}
	if err := c.Provide(fleet.NewPark); err != nil {
		return nil, err
	}
	if err := c.Provide(func(park *fleet.Park) *console.App {
		return console.New(park, in, out)
	}); err != nil {
		return nil, err
	}
	return c, nil
}
