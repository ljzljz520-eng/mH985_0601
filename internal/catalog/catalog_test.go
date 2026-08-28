package catalog

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

func TestCatalogRegisterChangeSearch(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "catalog.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	c := New(s)
	if _, err := c.Register(model.ImportRow{ID: "b", StoreID: "north", Title: "Fire", Content: "Exit", Category: "safety", Owner: "amy", SortKey: 2}, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Register(model.ImportRow{ID: "a", StoreID: "north", Title: "Food", Content: "Wash", Category: "hygiene", Owner: "amy", SortKey: 1}, 2); err != nil {
		t.Fatal(err)
	}
	changed, err := c.Change(model.ChangeRequest{RecordID: "a", Title: "Food handling", Actor: "amy", Seq: 3})
	if err != nil || changed.Status != model.StatusPending {
		t.Fatalf("change = %#v err=%v", changed, err)
	}
	items, err := c.Search(model.SearchQuery{StoreID: "north", Text: "food"})
	if err != nil || len(items) != 1 || items[0].ID != "a" {
		t.Fatalf("search = %#v err=%v", items, err)
	}
	if _, err := c.Change(model.ChangeRequest{RecordID: "a", Seq: 3}); err != store.ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestCatalogAttachmentAndDetail(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "detail.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	c := New(s)
	if _, err := c.Register(model.ImportRow{ID: "r", StoreID: "s", Title: "Guide", Content: "Body"}, 1); err != nil {
		t.Fatal(err)
	}
	if err := c.Attach(model.Attachment{ID: "att", RecordID: "r", Name: "photo", Kind: "image", Size: 4, Digest: "abcd"}); err != nil {
		t.Fatal(err)
	}
	detail, err := c.Detail("r")
	if err != nil || len(detail.Audit) != 1 || len(detail.Attachments) != 1 {
		t.Fatalf("detail = %#v err=%v", detail, err)
	}
}
