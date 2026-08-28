package notification

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestOutboxDedupAndDelivery(t *testing.T) {
	o := New()
	r := model.Record{ID: "r", Title: "Guide", Status: model.StatusApproved}
	first, err := o.Enqueue(r, "approved", "owner", 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := o.Enqueue(r, "approved", "owner", 2)
	if err != nil || first.ID != second.ID || len(o.Pending()) != 1 {
		t.Fatalf("first=%#v second=%#v pending=%#v", first, second, o.Pending())
	}
	if delivered := o.DeliverAll(); delivered != 1 || len(o.Pending()) != 0 {
		t.Fatalf("delivered=%d pending=%#v", delivered, o.Pending())
	}
}
