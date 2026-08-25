package maintenance

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestChecksAndNormalize(t *testing.T) {
	records := []model.Record{{ID: "b", Title: " B ", Content: "body", Status: model.StatusApproved}, {ID: "a", Status: model.StatusDraft}}
	checks := CheckRecords(records)
	if !Healthy(checks) {
		t.Fatalf("checks=%#v", checks)
	}
	normalized, err := Normalize(records)
	if err != nil || normalized[0].ID != "a" || normalized[1].Title != "B" {
		t.Fatalf("normalized=%#v err=%v", normalized, err)
	}
	snapshot := SnapshotRecords(records, map[string]bool{"b": true})
	if snapshot.Records != 2 || snapshot.Audited != 1 || snapshot.WithContent != 1 {
		t.Fatalf("snapshot=%#v", snapshot)
	}
}
