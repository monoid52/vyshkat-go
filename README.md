# VyshKat

A small Go console app for tracking rental vehicles and other equipment. Data is kept in memory and resets when the app exits.

## Run

Requires Go 1.22 or newer.

```sh
go run ./cmd/vyshkat          # interactive menu
go run ./cmd/vyshkat -demo    # all main scenarios without typing
go test ./... -cover          # tests and coverage
```

The menu can add an electric scooter, bicycle, electric bicycle, helmet or charging dock. Vehicles are accepted only if their brakes and frame pass inspection. Inventory numbers must be unique. The summary counts accepted vehicles and other assets and adds daily electricity use from scooters, electric bicycles and charging docks. The beginner list includes **vehicles** rated at least 6 out of 10; the inventory lists both vehicles and other assets.

## Design and course principles

- **Single responsibility:** [`fleet`](fleet/) owns the rules and [`internal/console`](internal/console/) handles input and output. The small [`Park`](fleet/park.go) both accepts assets and calculates reports; this is a deliberate limit on splitting responsibilities, because both operations use the same in-memory list.
- **Open/closed:** [`Park`](fleet/park.go) uses [`Vehicle`](fleet/model.go), [`Asset`](fleet/model.go) and [`EnergyConsumer`](fleet/model.go), not a switch over vehicle types. A new type needs to implement the relevant methods; only the menu needs a new choice if users should create it interactively.
- **Substitution:** scooters, bicycles and electric bicycles can each be passed to `Park.Add` as a `Vehicle`; no type-specific inspection or report path is needed.
- **Small interfaces:** a helmet needs an inventory number but no energy method; a charging dock uses energy but is not a rental vehicle.
- **Dependency inversion (most visible):** [`Park`](fleet/park.go) receives a `ServiceCenter` interface rather than constructing an inspection service. [`internal/wiring`](internal/wiring/container.go) uses the `dig` DI container to assemble the real app. The container does not appear in the fleet rules.

An employee should be modelled separately: employees are not inventory assets. A parts warehouse would need its own stock rules; it should not be added to `Park` just to reuse its reports. The current vehicle and asset reports do not need to change for either addition.

[`TestInspectionBeforeAcceptance`](fleet/park_test.go) checks that failed inspection excludes a vehicle from reports; [`TestReportsIncludeOnlyAcceptedAssets`](fleet/park_test.go) covers the beginner threshold and energy totals. [`TestContainerReplacesServiceCenter`](internal/wiring/container_test.go) replaces the real inspector with a test stub through the same DI container, then checks the resulting park state.

**AI use:** An AI assistant drafted the initial code, tests and README. No independent student changes are claimed here yet; update this statement after reviewing or changing the project yourself.
