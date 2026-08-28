package validation

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestInputAndPolicy(t *testing.T) {
	if err := RecordInput(model.ImportRow{ID: "x", StoreID: "s", Title: "T", Content: "C"}); err != nil {
		t.Fatal(err)
	}
	if err := SearchInput(model.SearchQuery{Limit: 501}); err == nil {
		t.Fatal("expected limit error")
	}
	r := model.Record{Status: model.StatusApproved, Owner: "owner", Content: "body"}
	if err := CanChange(r, "other"); err == nil {
		t.Fatal("expected ownership error")
	}
	if !CanPublish(r) {
		t.Fatal("approved record should publish")
	}
}
