package analytics

import (
	"strings"
	"testing"

	"trainingdesk/internal/model"
)

func TestDashboardMetrics(t *testing.T) {
	d := Build([]model.Record{{StoreID: "north", Category: "safety", Status: model.StatusApproved, Reviewer: "lee"}, {StoreID: "north", Category: "safety", Status: model.StatusPending}, {StoreID: "south", Category: "hygiene", Status: model.StatusArchived, Reviewer: "lee"}})
	if d.TotalRecords != 3 || d.VisibleRecords != 2 || d.PendingRecords != 1 || d.CompletionRate != 66 {
		t.Fatalf("dashboard=%#v", d)
	}
	if !strings.Contains(Headline(d), "completion=66%") {
		t.Fatalf("headline=%q", Headline(d))
	}
}
