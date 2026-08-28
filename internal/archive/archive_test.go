package archive

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/review"
	"trainingdesk/internal/store"
)

func TestArchiveLifecycle(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "archive.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	c := catalog.New(s)
	rv := review.New(c, s)
	a := New(c, s)
	if _, err := c.Register(model.ImportRow{ID: "r", StoreID: "s", Title: "Guide", Content: "Body"}, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := rv.Submit("r", 2); err != nil {
		t.Fatal(err)
	}
	if _, err := rv.Decide(model.ReviewRequest{RecordID: "r", Reviewer: "lee", Approve: true, Seq: 3}); err != nil {
		t.Fatal(err)
	}
	archived, err := a.Archive("r", "lee", 4)
	if err != nil || archived.Status != model.StatusArchived || !a.IsVisible(archived) {
		t.Fatalf("archived=%#v err=%v", archived, err)
	}
	restored, err := a.Restore("r", "lee", 5)
	if err != nil || restored.Status != model.StatusApproved {
		t.Fatalf("restored=%#v err=%v", restored, err)
	}
}
