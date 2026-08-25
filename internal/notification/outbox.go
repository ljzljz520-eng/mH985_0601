package notification

import (
	"errors"
	"fmt"
	"sort"

	"trainingdesk/internal/model"
)

type Message struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Kind     string `json:"kind"`
	Target   string `json:"target"`
	Body     string `json:"body"`
	Sent     bool   `json:"sent"`
	Seq      int64  `json:"seq"`
}

type Outbox struct {
	messages map[string]Message
}

func New() *Outbox {
	return &Outbox{messages: make(map[string]Message)}
}

func (o *Outbox) Enqueue(record model.Record, kind, target string, seq int64) (Message, error) {
	if record.ID == "" || kind == "" || target == "" {
		return Message{}, errors.New("notification identity is required")
	}
	if seq <= 0 {
		return Message{}, errors.New("notification sequence is required")
	}
	id := fmt.Sprintf("notice-%s-%d", record.ID, seq)
	message := Message{ID: id, RecordID: record.ID, Kind: kind, Target: target, Body: fmt.Sprintf("%s is now %s", record.Title, record.Status), Seq: seq}
	if existing, ok := o.messages[id]; ok {
		return existing, nil
	}
	o.messages[id] = message
	return message, nil
}

func (o *Outbox) Pending() []Message {
	items := make([]Message, 0)
	for _, message := range o.messages {
		if !message.Sent {
			items = append(items, message)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (o *Outbox) MarkSent(id string) error {
	message, ok := o.messages[id]
	if !ok {
		return errors.New("notification not found")
	}
	message.Sent = true
	o.messages[id] = message
	return nil
}

func (o *Outbox) ForRecord(recordID string) []Message {
	items := make([]Message, 0)
	for _, message := range o.messages {
		if message.RecordID == recordID {
			items = append(items, message)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Seq < items[j].Seq })
	return items
}

func (o *Outbox) DeliverAll() int {
	delivered := 0
	for _, message := range o.Pending() {
		if o.MarkSent(message.ID) == nil {
			delivered++
		}
	}
	return delivered
}
