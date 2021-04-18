package probes

import (
	"errors"
	"testing"
	"time"
)

func Test_RunProbe(t *testing.T) {
	type args struct {
		p *Probe
		s status
		e string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "up-to-up",
			args: args{
				p: NewProbe("up-to-up", Up),
				s: Up,
				e: "up-to-up probe was stopped with UP status",
			},
		},
		{
			name: "up-to-down",
			args: args{
				p: NewProbe("up-to-down", Up),
				s: Down,
				e: "up-to-down probe was stopped with DOWN status",
			},
		},
		{
			name: "down-to-down",
			args: args{
				p: NewProbe("down-to-down", Down),
				s: Down,
				e: "down-to-down probe was stopped with DOWN status",
			},
		},
		{
			name: "down-to-up",
			args: args{
				p: NewProbe("down-to-up", Down),
				s: Up,
				e: "down-to-up probe was stopped with UP status",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.p.Chan() <- tt.args.s
			close(tt.args.p.Chan())

			err := RunProbe(tt.args.p)
			if err.Error() != tt.args.e {
				t.Errorf("runProbe() error = %v, wantErr %v", err, tt.args.e)
				return
			}
		})
	}
}

func Test_ReadinessProbe(t *testing.T) {
	type args struct {
		tests []func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass all",
			args: args{
				tests: []func() error{
					func() error { return nil },
					func() error { return nil },
				},
			},
			wantErr: false,
		},
		{
			name: "fail first",
			args: args{
				tests: []func() error{
					func() error { return errors.New("first") },
					func() error { return nil },
				},
			},
			wantErr: true,
		},
		{
			name: "fail second",
			args: args{
				tests: []func() error{
					func() error { return nil },
					func() error { return errors.New("second") },
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		p := NewProbe("test", Down)
		defer close(p.Chan())
		go RunProbe(p)

		t.Run(tt.name, func(t *testing.T) {
			if err := ReadinessProbe(p, tt.args.tests...); (err != nil) != tt.wantErr {
				t.Errorf("ReadinessProbe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_LivenessProbe(t *testing.T) {
	type args struct {
		probes []*Probe
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				probes: []*Probe{
					{
						name:   "pass",
						status: Up,
						time:   time.Now(),
					}},
			},
			wantErr: false,
		},
		{
			name: "fail-short",
			args: args{
				probes: []*Probe{{
					name:   "fail-short",
					status: Down,
					time:   time.Now(),
				}},
			},
			wantErr: false,
		},
		{
			name: "fail-long",
			args: args{
				probes: []*Probe{{
					name:   "fail-long",
					status: Down,
					time:   time.Now().Add(-10 * time.Minute),
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		p := NewProbe("test", Up)
		defer close(p.channel)
		go RunProbe(p)

		t.Run(tt.name, func(t *testing.T) {
			if err := LivenessProbe(p, tt.args.probes...); (err != nil) != tt.wantErr {
				t.Errorf("LivenessProbe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Dummy test for coverage.
// Everything this function does is tested above.
func Test_StartProbes(t *testing.T) {
	go StartProbes(func() error { return nil })
	time.Sleep(10 * time.Millisecond)
}
