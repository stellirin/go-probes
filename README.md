# Probes

[![reference](https://pkg.go.dev/badge/czechia.dev/probes.svg)](https://pkg.go.dev/czechia.dev/probes)
[![report](https://goreportcard.com/badge/czechia.dev/probes)](https://goreportcard.com/report/czechia.dev/probes)
[![tests](https://github.com/stellirin/go-probes/workflows/Go/badge.svg)](https://github.com/stellirin/go-probes/actions?query=workflow%3AGo)
[![coverage](https://codecov.io/gh/stellirin/go-probes/branch/main/graph/badge.svg?token=9ATrVTllue)](https://codecov.io/gh/stellirin/go-probes)

Probes is a simple package to implement readiness and liveness endpoints.

## ⚙️ Installation

```sh
go get -u czechia.dev/probes
```

## 👀 Example

```go
package main

import (
	"errors"
	"net/http"

	"czechia.dev/probes"
	"github.com/labstack/echo/v4"
)

const alive = true

func isAlive() error {
	if alive {
		return nil
	}
	return errors.New("dead")
}

func main() {
	go probes.StartProbes(isAlive)

	e := echo.New()
	e.GET("/liveness", probeRoute(probes.Liveness))
	e.GET("/readiness", probeRoute(probes.Readiness))
	e.Start(":8080")
}

func probeRoute(p *probes.Probe) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if p.IsUp() {
			return ctx.NoContent(http.StatusOK)
		}
		return ctx.NoContent(http.StatusServiceUnavailable)
	}
}
```
