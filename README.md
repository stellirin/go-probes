# Probes

[![codecov](https://codecov.io/gh/stellirin/go-probes/branch/main/graph/badge.svg?token=9ATrVTllue)](https://codecov.io/gh/stellirin/go-probes)
[![Test Action Status](https://github.com/stellirin/go-probes/workflows/Go/badge.svg)](https://github.com/stellirin/go-probes/actions?query=workflow%3AGo)

Probes is a simple package to implement readiness and liveness endpoints.

## ⚙️ Installation

```sh
go get -u czechia.dev/probes
```

## 👀 Example

```go
package main

import (
	"time"

	"czechia.dev/probes"
	"github.com/gofiber/fiber/v2"
)

func NewFiberRoute(p *probes.Probe) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		switch p.IsDown() {
		case false:
			return ctx.SendStatus(fiber.StatusOK)
		case true:
			return ctx.SendStatus(fiber.StatusServiceUnavailable)
		}
	}
}

func main() {
	go probes.RunProbe(probes.Liveness)
	go probes.RunProbe(probes.Readiness)

	app := fiber.New()
	app.Get("/liveness", NewFiberRoute(probes.Liveness))
	app.Get("/readiness", NewFiberRoute(probes.Readiness))
	go app.Listen(":8080")

	for ; true; <-time.NewTicker(3 * time.Second).C {
		probes.LivenessProbe(probes.Liveness)
		probes.ReadinessProbe(probes.Readiness, func() error { return nil })
	}
}
```
