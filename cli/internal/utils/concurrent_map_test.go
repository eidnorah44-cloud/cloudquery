package utils

import (
	"sync"
	"testing"
)

func TestConcurrentMap(t *testing.T) {
	m := NewConcurrentMap[string, int]()

	// Test Add & Get
	m.Add("a", 1)
	m.Add("b", 2)

	val, ok := m.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1, got %d (ok=%v)", val, ok)
	}

	val, ok = m.Get("b")
	if !ok || val != 2 {
		t.Errorf("expected 2, got %d (ok=%v)", val, ok)
	}

	_, ok = m.Get("c")
	if ok {
		t.Errorf("expected not ok for missing key")
	}

	// Test SetIfGreater
	m.SetIfGreater("a", 10, func(newVal, existingVal int) bool { return newVal > existingVal })
	val, _ = m.Get("a")
	if val != 10 {
		t.Errorf("expected 10 after SetIfGreater with higher value, got %d", val)
	}

	m.SetIfGreater("a", 5, func(newVal, existingVal int) bool { return newVal > existingVal })
	val, _ = m.Get("a")
	if val != 10 {
		t.Errorf("expected 10 after SetIfGreater with lower value, got %d", val)
	}

	m.SetIfGreater("new_key", 100, func(newVal, existingVal int) bool { return newVal > existingVal })
	val, _ = m.Get("new_key")
	if val != 100 {
		t.Errorf("expected 100 for new_key, got %d", val)
	}

	// Test GetAll
	all := m.GetAll()
	if len(all) != 3 || all["a"] != 10 || all["b"] != 2 || all["new_key"] != 100 {
		t.Errorf("unexpected GetAll result: %v", all)
	}

	// Test concurrent operations
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Add("key", i)
			m.SetIfGreater("key", i+1, func(newVal, existingVal int) bool { return newVal > existingVal })
			m.Get("key")
			m.GetAll()
		}(i)
	}
	wg.Wait()
}
