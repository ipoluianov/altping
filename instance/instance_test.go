package instance

import (
	"testing"
	"time"
)

func TestSecondStartShowsTheFirst(t *testing.T) {
	dir := t.TempDir()

	first, ok := Acquire(dir)
	if !ok {
		t.Fatal("the first start must run")
	}
	shown := make(chan struct{}, 1)
	first.Serve(func() { shown <- struct{}{} })

	if _, ok := Acquire(dir); ok {
		t.Fatal("the second start must quit while the first runs")
	}
	select {
	case <-shown:
	case <-time.After(3 * time.Second):
		t.Fatal("the first copy was not asked to show itself")
	}

	first.Close()
	again, ok := Acquire(dir)
	if !ok {
		t.Fatal("a start after the first has quit must run")
	}
	again.Close()
}
