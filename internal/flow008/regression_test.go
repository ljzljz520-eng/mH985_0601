package flow008

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/review"
	"trainingdesk/internal/store"
)

func Test985BusinessRegression(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "regression.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	c := catalog.New(s)
	rv := review.New(c, s)
	h := New(c, rv)
	for _, row := range []model.ImportRow{{ID: "record-a", StoreID: "flagship", Title: "Opening", Content: "Open", SortKey: 1}, {ID: "record-b", StoreID: "flagship", Title: "Closing", Content: "Close", SortKey: 1}} {
		if _, err := c.Register(row, 1); err != nil {
			t.Fatal(err)
		}
	}
	got, err := h.DetailAt("flagship", 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "record-b" || got.Title != "Closing" {
		t.Fatalf("selected record = %#v", got)
	}
}
