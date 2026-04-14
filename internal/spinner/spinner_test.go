package spinner

import (
	"testing"
)

func TestSpinner_NewNotNil(t *testing.T) {
	s := New("loading")
	if s == nil {
		t.Fatal("expected non-nil spinner")
	}
	if s.message != "loading" {
		t.Errorf("expected message %q, got %q", "loading", s.message)
	}
}

func TestSpinner_StopWithMessage(t *testing.T) {
	s := New("working")
	s.Start()
	// Should not panic or hang.
	s.Stop("done")
}

func TestSpinner_Update(t *testing.T) {
	s := New("step 1")
	s.Start()
	// Should not panic.
	s.Update("step 2")
	if s.message != "step 2" {
		t.Errorf("expected message %q after update, got %q", "step 2", s.message)
	}
	s.Stop("")
}
