package ledger

import (
	"errors"
	"fmt"
	"sort"

	"trainingdesk/internal/model"
)

type Entry struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Kind     string `json:"kind"`
	From     string `json:"from"`
	To       string `json:"to"`
	Actor    string `json:"actor"`
	Seq      int64  `json:"seq"`
	Note     string `json:"note"`
}

type Ledger struct {
	entries map[string]Entry
}

func New() *Ledger {
	return &Ledger{entries: make(map[string]Entry)}
}

func (l *Ledger) Append(record model.Record, from, to, actor, note string, seq int64) (Entry, error) {
	if record.ID == "" {
		return Entry{}, errors.New("record is required")
	}
	if from == to {
		return Entry{}, errors.New("ledger transition must change status")
	}
	if actor == "" || seq <= 0 {
		return Entry{}, errors.New("ledger actor and sequence are required")
	}
	id := fmt.Sprintf("entry-%s-%012d", record.ID, seq)
	if existing, ok := l.entries[id]; ok {
		return existing, nil
	}
	entry := Entry{ID: id, RecordID: record.ID, Kind: "status-change", From: from, To: to, Actor: actor, Seq: seq, Note: note}
	l.entries[id] = entry
	return entry, nil
}

func (l *Ledger) ForRecord(recordID string) []Entry {
	items := make([]Entry, 0)
	for _, entry := range l.entries {
		if entry.RecordID == recordID {
			items = append(items, entry)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Seq == items[j].Seq {
			return items[i].ID < items[j].ID
		}
		return items[i].Seq < items[j].Seq
	})
	return items
}

func (l *Ledger) Latest(recordID string) (Entry, error) {
	items := l.ForRecord(recordID)
	if len(items) == 0 {
		return Entry{}, errors.New("ledger entry not found")
	}
	return items[len(items)-1], nil
}

func (l *Ledger) Verify(recordID string, initial model.RecordStatus) error {
	items := l.ForRecord(recordID)
	current := string(initial)
	var previous int64
	for _, entry := range items {
		if entry.Seq <= previous {
			return errors.New("ledger sequence is not increasing")
		}
		if entry.From != current {
			return fmt.Errorf("ledger expected from %s, got %s", current, entry.From)
		}
		current = entry.To
		previous = entry.Seq
	}
	return nil
}

func (l *Ledger) Count(kind string) int {
	count := 0
	for _, entry := range l.entries {
		if kind == "" || entry.Kind == kind {
			count++
		}
	}
	return count
}

func (l *Ledger) RemoveRecord(recordID string) int {
	removed := 0
	for id, entry := range l.entries {
		if entry.RecordID == recordID {
			delete(l.entries, id)
			removed++
		}
	}
	return removed
}
