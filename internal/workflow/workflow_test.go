package workflow

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/store"
)

func TestWorkflowEngine(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	e := New(s)
	w, err := e.Start("r", "training")
	if err != nil {
		t.Fatal(err)
	}
	if next, _ := NextStep(w); next != "register" {
		t.Fatalf("next=%q", next)
	}
	w, err = e.CompleteStep(w.ID, "register")
	if err != nil || w.Stage != "register" {
		t.Fatalf("workflow=%#v err=%v", w, err)
	}
	if IsComplete(w) {
		t.Fatal("workflow should not be complete")
	}
}
