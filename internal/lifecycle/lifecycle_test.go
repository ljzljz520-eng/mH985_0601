package lifecycle

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestTimeline(t *testing.T) {
	timeline, err := New("r")
	if err != nil {
		t.Fatal(err)
	}
	if err := timeline.Add(Event{From: model.StatusDraft, To: model.StatusPending, Actor: "owner", Seq: 1}); err != nil {
		t.Fatal(err)
	}
	if err := timeline.Add(Event{From: model.StatusPending, To: model.StatusApproved, Actor: "reviewer", Seq: 2, Action: "approved"}); err != nil {
		t.Fatal(err)
	}
	if timeline.Current(model.StatusDraft) != model.StatusApproved || timeline.LastActor() != "reviewer" || !timeline.Contains("approved") || Describe(timeline, model.StatusDraft) != "draft -> pending -> approved" {
		t.Fatalf("timeline=%#v", timeline)
	}
	if err := timeline.Add(Event{From: model.StatusDraft, To: model.StatusApproved, Actor: "bad", Seq: 3}); err == nil {
		t.Fatal("expected invalid transition")
	}
}
