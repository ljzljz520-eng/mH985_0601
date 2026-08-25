package report

import (
	"strings"
	"testing"

	"trainingdesk/internal/model"
)

func TestBuildAndRender(t *testing.T) {
	summary := Build([]model.Record{{ID: "1", StoreID: "north", Category: "safety", Status: model.StatusApproved}, {ID: "2", StoreID: "north", Category: "hygiene", Status: model.StatusPending}, {ID: "3", StoreID: "south", Category: "safety", Status: model.StatusArchived}})
	if summary.Total != 3 || summary.Stores != 2 || summary.Approved != 1 || summary.Archived != 1 {
		t.Fatalf("summary=%#v", summary)
	}
	rendered := Render(summary)
	if !strings.Contains(rendered, "north=2") || !strings.Contains(rendered, "pending=1") {
		t.Fatalf("rendered=%q", rendered)
	}
	if rows := ExportRows([]model.Record{{ID: "b"}, {ID: "a"}}); rows[0] != "a||||" {
		t.Fatalf("rows=%v", rows)
	}
}
