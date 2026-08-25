package review

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

func TestReviewApproveAndReject(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "review.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	c := catalog.New(s)
	rv := New(c, s)
	if _, err := c.Register(model.ImportRow{ID: "r", StoreID: "s", Title: "Guide", Content: "Body"}, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := rv.Submit("r", 2); err != nil {
		t.Fatal(err)
	}
	approved, err := rv.Decide(model.ReviewRequest{RecordID: "r", Reviewer: "lee", Approve: true, Seq: 3})
	if err != nil || approved.Status != model.StatusApproved {
		t.Fatalf("approved=%#v err=%v", approved, err)
	}
	if _, err := rv.Decide(model.ReviewRequest{RecordID: "r", Reviewer: "lee", Approve: false, Seq: 4}); err == nil {
		t.Fatal("expected non-pending decision error")
	}
}

func TestReviewPendingFilter(t *testing.T) {
	rv := &Service{}
	items := rv.Pending([]model.Record{{ID: "a", Status: model.StatusPending}, {ID: "b", Status: model.StatusApproved}})
	if len(items) != 1 || items[0].ID != "a" {
		t.Fatalf("pending=%#v", items)
	}
}
