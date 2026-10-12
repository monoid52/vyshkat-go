package main

import (
	"flag"
	"log"
	"os"

	"github.com/monoid52/vyshkat-go/fleet"
	"github.com/monoid52/vyshkat-go/internal/console"
	"github.com/monoid52/vyshkat-go/internal/wiring"
)

func main() {
	demo := flag.Bool("demo", false, "run a non-interactive demonstration")
	flag.Parse()

	container, err := wiring.NewContainer(os.Stdin, os.Stdout, func() fleet.ServiceCenter {
		return fleet.SafetyServiceCenter{}
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := container.Invoke(func(app *console.App) error {
		if *demo {
			return app.RunDemo()
		}
		return app.Run()
	}); err != nil {
		log.Fatal(err)
	}
}
