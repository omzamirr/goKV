package store

import (
	"testing"
	"time"
)

func TestPassiveEviction(t *testing.T) {
	// 1. Initialize a new store using your factory function
	s := New()

	// 2. Set a key that should expire incredibly fast (50 milliseconds)
	s.Set("flash", "now you see me", 50*time.Millisecond)

	// 3. Immediately try to Get it. It should still be alive
	val, ok := s.Get("flash")
	if !ok || val != "now you see me" {
		t.Fatalf("Expected key to be alive immediately, got ok=%v, val=%s", ok, val)
	}

	// 4. Sleep for 60 milliseconds so the deadline passes
	time.Sleep(60 * time.Millisecond)

	// 5. Try to Get it again. Your passive eviction logic should trigger right here!
	val, ok = s.Get("flash")
	if ok {
		t.Errorf("Expected key to be passively evicted, but it still exists with value: %s", val)
	}
}
