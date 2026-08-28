package flow008

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/archive"
	"trainingdesk/internal/catalog"
	"trainingdesk/internal/importer"
	"trainingdesk/internal/model"
	"trainingdesk/internal/review"
	"trainingdesk/internal/store"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, c, rv, ar := fixture(t)
	defer s.Close()
	if _, err := c.Register(model.ImportRow{ID: "guide-1", StoreID: "north", Title: "Opening", Content: "Open safely", Category: "operations", SortKey: 1}, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := rv.Submit("guide-1", 2); err != nil {
		t.Fatal(err)
	}
	if _, err := rv.Decide(model.ReviewRequest{RecordID: "guide-1", Reviewer: "reviewer", Approve: true, Seq: 3}); err != nil {
		t.Fatal(err)
	}
	got, err := ar.Archive("guide-1", "reviewer", 4)
	if err != nil || got.Status != model.StatusArchived {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, c, rv, _ := fixture(t)
	defer s.Close()
	if _, err := c.Register(model.ImportRow{ID: "guide-2", StoreID: "south", Title: "Closing", Content: "Close", Category: "operations", SortKey: 1}, 1); err != nil {
		t.Fatal(err)
	}
	items, err := c.Search(model.SearchQuery{Text: "closing"})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%#v err=%v", items, err)
	}
	if _, err := c.Change(model.ChangeRequest{RecordID: items[0].ID, Title: "Closing v2", Actor: "editor", Seq: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err := rv.Decide(model.ReviewRequest{RecordID: "guide-2", Reviewer: "reviewer", Approve: true, Seq: 3}); err != nil {
		t.Fatal(err)
	}
	got, err := c.Get("guide-2")
	if err != nil || got.Status != model.StatusApproved || got.Title != "Closing v2" {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	s, c, _, _ := fixture(t)
	defer s.Close()
	imp := importer.New(c)
	report := imp.Import([]model.ImportRow{{ID: "guide-3", StoreID: "north", Title: "Cash", Content: "Count"}, {ID: "", StoreID: "north", Title: "Broken", Content: "Missing"}}, 10)
	if report.Accepted != 1 || report.Rejected != 1 || len(report.Errors) != 1 {
		t.Fatalf("report=%#v", report)
	}
	if _, err := c.Get("guide-3"); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) (*store.Store, *catalog.Catalog, *review.Service, *archive.Service) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "flow.db"))
	if err != nil {
		t.Fatal(err)
	}
	c := catalog.New(s)
	rv := review.New(c, s)
	return s, c, rv, archive.New(c, s)
}
