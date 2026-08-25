package lifecycle

import (
	"errors"
	"fmt"
	"strings"

	"trainingdesk/internal/model"
)

type Event struct {
	Action string
	From   model.RecordStatus
	To     model.RecordStatus
	Actor  string
	Seq    int64
	Note   string
}

type Timeline struct {
	RecordID string
	Events   []Event
}

func New(recordID string) (*Timeline, error) {
	if strings.TrimSpace(recordID) == "" {
		return nil, errors.New("record id is required")
	}
	return &Timeline{RecordID: recordID, Events: make([]Event, 0)}, nil
}

func (t *Timeline) Add(event Event) error {
	if event.Actor == "" || event.Seq <= 0 {
		return errors.New("event actor and sequence are required")
	}
	if !Allowed(event.From, event.To) {
		return fmt.Errorf("transition %s to %s is not allowed", event.From, event.To)
	}
	if len(t.Events) > 0 {
		last := t.Events[len(t.Events)-1]
		if event.Seq <= last.Seq {
			return errors.New("event sequence must increase")
		}
		if event.From != last.To {
			return fmt.Errorf("event starts at %s, expected %s", event.From, last.To)
		}
	}
	if strings.TrimSpace(event.Action) == "" {
		event.Action = actionFor(event.To)
	}
	t.Events = append(t.Events, event)
	return nil
}

func Allowed(from, to model.RecordStatus) bool {
	if from == model.StatusDraft && to == model.StatusPending {
		return true
	}
	if from == model.StatusPending && (to == model.StatusApproved || to == model.StatusRejected) {
		return true
	}
	if from == model.StatusRejected && to == model.StatusPending {
		return true
	}
	if from == model.StatusApproved && to == model.StatusArchived {
		return true
	}
	if from == model.StatusArchived && to == model.StatusApproved {
		return true
	}
	return false
}

func actionFor(status model.RecordStatus) string {
	switch status {
	case model.StatusPending:
		return "submitted"
	case model.StatusApproved:
		return "approved"
	case model.StatusRejected:
		return "rejected"
	case model.StatusArchived:
		return "archived"
	default:
		return "changed"
	}
}

func (t *Timeline) Current(initial model.RecordStatus) model.RecordStatus {
	if len(t.Events) == 0 {
		return initial
	}
	return t.Events[len(t.Events)-1].To
}

func (t *Timeline) LastActor() string {
	if len(t.Events) == 0 {
		return ""
	}
	return t.Events[len(t.Events)-1].Actor
}

func (t *Timeline) Actions() []string {
	actions := make([]string, 0, len(t.Events))
	for _, event := range t.Events {
		actions = append(actions, event.Action)
	}
	return actions
}

func (t *Timeline) Contains(action string) bool {
	for _, event := range t.Events {
		if event.Action == action {
			return true
		}
	}
	return false
}

func Describe(t *Timeline, initial model.RecordStatus) string {
	parts := []string{string(initial)}
	for _, event := range t.Events {
		parts = append(parts, string(event.To))
	}
	return strings.Join(parts, " -> ")
}
