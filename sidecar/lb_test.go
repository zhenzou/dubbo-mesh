package sidecar

import (
	"testing"

	"dubbo-mesh/registry"
	"dubbo-mesh/util"
)

func TestWeightRoundRobin(t *testing.T) {
	wrr := &WeightRoundRobin{}
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

func BenchmarkWeightRoundRobin(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	wrr := &WeightRoundRobin{}
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
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		end := wrr.Elect(endpoints)
		count[end.System.TotalMemory] = count[end.System.TotalMemory] + 1
	}
	println(util.ToJsonStr(count))
}

func TestWeightRandom(t *testing.T) {
	wr := &WeightRandom{}
	type args struct {
		endpoints []*Endpoint
	}

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
	wr.Init(endpoints)
	count := map[int]int{}
	for i := 0; i < 100000; i++ {
		end := wr.Elect(endpoints)
		count[end.System.TotalMemory] = count[end.System.TotalMemory] + 1
	}
	println(util.ToJsonStr(count))
}
func BenchmarkWeightRandom(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	wr := &WeightRandom{}

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
	wr.Init(endpoints)
	count := map[int]int{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		end := wr.Elect(endpoints)
		count[end.System.TotalMemory] = count[end.System.TotalMemory] + 1
	}
	println(util.ToJsonStr(count))
}
