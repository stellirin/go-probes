package probes

import (
	"fmt"
	"time"
)

type status string

const (
	Up   status = "UP"
	Down status = "DOWN"
)

type Probe struct {
	name    string
	status  status
	channel chan status
	time    time.Time
}

type ProbeFunc func() error

var (
	Liveness  *Probe
	Readiness *Probe
)

func init() {
	Liveness = New("liveness", Up)
	Readiness = New("readiness", Down)
}

// New initializes a new probe with an initial status.
func New(n string, s status) *Probe {
	return &Probe{
		name:    n,
		status:  s,
		channel: make(chan status, 1),
		time:    time.Now(),
	}
}

// Chan exposes the write end of the probe channel.
func (p *Probe) Chan() chan<- status {
	return p.channel
}

// IsUp tests if the probe status is UP.
func (p *Probe) IsUp() bool {
	return p.status == Up
}

// Downtime is the duration since the probe went DOWN.
func (p *Probe) Downtime() time.Duration {
	return time.Since(p.time)
}

// RunProbe waits for status messages on the probe channel.
//
// The probe status is updated to reflect the current status.
// The probe timestamp is also updated when the status is UP.
//
// This function typically runs in its own goroutine.
// The return parameter may be used for tests.
func RunProbe(p *Probe) error {
	for s := range p.channel {
		if s == Up {
			p.time = time.Now()
		}

		p.status = s
	}

	return fmt.Errorf("%s probe was stopped with %s status", p.name, p.status)
}

// ReadinessProbe runs any probe functions passed to it.
// If any of the tests fail the given probe is set to DOWN.
//
// This function typically runs in its own goroutine.
// The return parameter may be used for tests.
func ReadinessProbe(p *Probe, pf ...ProbeFunc) error {
	for _, probeFunc := range pf {
		if err := probeFunc(); err != nil {
			p.Chan() <- Down
			return err
		}
	}

	// All tests passed
	p.Chan() <- Up
	return nil
}

// LivenessProbe sets the given Liveness probe to DOWN
// if any given Readiness probe is DOWN for more than 5 minutes.
//
// This function typically runs in its own goroutine.
// The return parameter may be used for tests.
func LivenessProbe(liveness *Probe, readiness ...*Probe) error {
	for _, p := range readiness {
		if p.IsUp() {
			continue
		}
		if p.Downtime() < 5*time.Minute {
			continue
		}

		liveness.Chan() <- Down
		return fmt.Errorf("%s probe down for too long", p.name)
	}

	// All tests passed
	liveness.Chan() <- Up
	return nil
}

// StartProbes is a convenience function to run the default Readiness
// and Liveness probes every 3 seconds using the given probe functions.
func StartProbes(pf ...ProbeFunc) {
	go RunProbe(Liveness)
	go RunProbe(Readiness)

	for ; true; <-time.NewTicker(3 * time.Second).C {
		LivenessProbe(Liveness)
		ReadinessProbe(Readiness, pf...)
	}
}
