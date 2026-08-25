package ledger

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestLedgerTransitions(t *testing.T) {
	l := New()
	r := model.Record{ID: "r", Status: model.StatusDraft}
	if _, err := l.Append(r, "draft", "pending", "owner", "submit", 1); err != nil {
		t.Fatal(err)
	}
	r.Status = model.StatusPending
	if _, err := l.Append(r, "pending", "approved", "reviewer", "ok", 2); err != nil {
		t.Fatal(err)
	}
	if err := l.Verify("r", model.StatusDraft); err != nil {
		t.Fatal(err)
	}
	r.Status = model.StatusApproved
	result := l.Reconcile(r)
	if !result.Balanced || len(l.TransitionPath("r")) != 3 || !ValidTransition(model.StatusPending, model.StatusApproved) {
		t.Fatalf("result=%#v path=%v", result, l.TransitionPath("r"))
	}
}
