package exporter

import (
	"bytes"
	"strings"
	"testing"

	"trainingdesk/internal/model"
)

func TestCSVAndJSONExport(t *testing.T) {
	records := []model.Record{{ID: "b", StoreID: "south", Title: "B", Status: model.StatusApproved, Version: 2, Content: "secret"}, {ID: "a", StoreID: "north", Title: "A", Status: model.StatusDraft, Version: 1, Content: "text"}}
	var out bytes.Buffer
	if err := CSV(&out, records, Options{IncludeContent: false}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "id,store_id,title,category,status,version") || strings.Contains(out.String(), "secret") {
		t.Fatalf("csv=%q", out.String())
	}
	data, err := JSON(records, Options{StoreID: "north", IncludeContent: true})
	if err != nil || !strings.Contains(string(data), "text") || strings.Contains(string(data), "secret") {
		t.Fatalf("json=%s err=%v", data, err)
	}
}
