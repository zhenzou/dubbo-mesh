package sidecar

import (
	"reflect"
	"testing"
	"dubbo-mesh/registry"
	"dubbo-mesh/util"
)

func TestWrrRandom_Elect(t *testing.T) {
	wrr := &WrrRandom{}
	endpoints := []*Endpoint{
		&Endpoint{
			Endpoint: &registry.Endpoint{System: &registry.System{TotalMemory: 2048}},
		},
		&Endpoint{
			Endpoint: &registry.Endpoint{System: &registry.System{TotalMemory: 4096}},
		},
		&Endpoint{
			Endpoint: &registry.Endpoint{System: &registry.System{TotalMemory: 6144}},
		},
	}
	wrr.Init(endpoints)
	count := map[int]int{}
	for i := 0; i < 100000; i++ {
		end := wrr.Elect(endpoints)
		count[end.System.TotalMemory] = count[end.System.TotalMemory] + 1
	}
	println(util.ToJsonStr(count))
}

func TestDrrRandom_Elect(t *testing.T) {
	type fields struct {
		weights map[*Endpoint]*int
		total   int
	}
	type args struct {
		endpoints []*Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Endpoint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &DrrRandom{
				weights: tt.fields.weights,
				total:   tt.fields.total,
			}
			if got := this.Elect(tt.args.endpoints); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DrrRandom.Elect() = %v, want %v", got, tt.want)
			}
		})
	}
}
