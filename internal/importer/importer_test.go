package importer

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

func TestImportReport(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "import.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	i := New(catalog.New(s))
	report := i.Import([]model.ImportRow{{ID: "ok", StoreID: "s", Title: "One", Content: "x"}, {ID: "", StoreID: "s", Title: "Bad", Content: "x"}}, 10)
	if report.Accepted != 1 || report.Rejected != 1 || len(report.IDs) != 1 {
		t.Fatalf("report=%#v", report)
	}
	if errs := i.ValidateRows([]model.ImportRow{{ID: "x", StoreID: "s"}, {ID: "x", StoreID: "", SortKey: -1}}); len(errs) != 3 {
		t.Fatalf("validation errors=%v", errs)
	}
}
