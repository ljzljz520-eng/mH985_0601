package lifecycle

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestChecklist(t *testing.T) {
	c, err := NewChecklist("r", []string{"review", "content", "review"})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Complete("content", 1); err != nil {
		t.Fatal(err)
	}
	if len(c.Pending()) != 1 || c.Ready() {
		t.Fatalf("pending=%v ready=%v", c.Pending(), c.Ready())
	}
	r := model.Record{ID: "r", Content: "body", Status: model.StatusApproved, UpdatedSeq: 2, PublishedSeq: 3}
	if !BuildChecklist(r).Ready() {
		t.Fatal("approved published record should have complete checklist")
	}
}
