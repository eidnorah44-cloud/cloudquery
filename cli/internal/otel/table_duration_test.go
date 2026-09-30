package otel

import (
	"fmt"
	"testing"
	"time"
)

func TestTableDurations(t *testing.T) {
	td := NewTableDurations()

	td.Set("table1", 100*time.Millisecond)
	td.Set("table1", 200*time.Millisecond)
	td.Set("table1", 150*time.Millisecond)

	td.Set("table2", 50*time.Millisecond)

	all := td.GetAll()
	if all["table1"] != 200*time.Millisecond {
		t.Errorf("expected table1 duration 200ms, got %v", all["table1"])
	}
	if all["table2"] != 50*time.Millisecond {
		t.Errorf("expected table2 duration 50ms, got %v", all["table2"])
	}
}

func BenchmarkTableDurations_Set(b *testing.B) {
	td := NewTableDurations()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			table := fmt.Sprintf("table-%d", i%50)
			td.Set(table, time.Duration(i)*time.Millisecond)
		}
	})
}
