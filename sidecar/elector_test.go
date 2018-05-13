package sidecar

import "testing"

func TestDrrRandom_dw(t *testing.T) {
	type fields struct {
		weights map[*Endpoint]*int
		total   int
	}
	type args struct {
		status *Status
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &DrrRandom{
				weights: tt.fields.weights,
				total:   tt.fields.total,
			}
			if got := this.dw(tt.args.status); got != tt.want {
				t.Errorf("DrrRandom.dw() = %v, want %v", got, tt.want)
			}
		})
	}
}
