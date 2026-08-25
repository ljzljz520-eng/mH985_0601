package query

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestParseAndApply(t *testing.T) {
	request, err := Parse("store=north status=approved sort=-title limit=1")
	if err != nil {
		t.Fatal(err)
	}
	records := []model.Record{{ID: "a", StoreID: "north", Title: "Alpha", Status: model.StatusApproved}, {ID: "b", StoreID: "north", Title: "Beta", Status: model.StatusApproved}, {ID: "c", StoreID: "south", Title: "Gamma", Status: model.StatusApproved}}
	items := Apply(records, request)
	if len(items) != 1 || items[0].ID != "b" {
		t.Fatalf("items=%#v", items)
	}
	if _, err := Parse("unknown=x"); err == nil {
		t.Fatal("expected invalid expression")
	}
}
