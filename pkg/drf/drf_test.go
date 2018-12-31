package drf_test

import (
	"testing"

	"github.com/andy2046/gopie/pkg/drf"
)

func TestDRF(t *testing.T) {
	nodeA := drf.NewNode(map[drf.Typ]float64{
		drf.CPU:    1,
		drf.MEMORY: 4,
	})
	nodeB := drf.NewNode(map[drf.Typ]float64{
		drf.CPU:    3,
		drf.MEMORY: 1,
	})
	cluster, _ := drf.New(map[drf.Typ]float64{
		drf.CPU:    9,
		drf.MEMORY: 18,
	}, nodeA, nodeB)
	var err error
	for err == nil {
		err = cluster.NextTask()
	}
	t.Logf("Error from NextTask %v", err)
	t.Logf("Resource for cluster %v", cluster.Resource())
	t.Logf("Consumed by cluster %v", cluster.Consumed())
	allcNodeA := nodeA.Allocated()
	t.Logf("Allocated for node A %v", allcNodeA)
	if allcNodeA[drf.CPU] != 3 {
		t.Fatalf("got %v, want %v", allcNodeA[drf.CPU], 3)
	}
	allcNodeB := nodeB.Allocated()
	t.Logf("Allocated for node B %v", allcNodeB)
	if allcNodeB[drf.CPU] != 6 {
		t.Fatalf("got %v, want %v", allcNodeB[drf.CPU], 6)
	}
}

func BenchmarkDRF(b *testing.B) {
	nodeA := drf.NewNode(map[drf.Typ]float64{
		drf.CPU:    1,
		drf.MEMORY: 4,
	})
	nodeB := drf.NewNode(map[drf.Typ]float64{
		drf.CPU:    3,
		drf.MEMORY: 1,
	})
	cluster, _ := drf.New(map[drf.Typ]float64{
		drf.CPU:    9,
		drf.MEMORY: 18,
	}, nodeA, nodeB)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cluster.NextTask()
	}
}
