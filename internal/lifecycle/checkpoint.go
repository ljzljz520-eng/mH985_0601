package lifecycle

import (
	"errors"
	"sort"

	"trainingdesk/internal/model"
)

type Checkpoint struct {
	Name      string
	Required  bool
	Completed bool
	Seq       int64
}

type Checklist struct {
	RecordID string
	Items    map[string]Checkpoint
}

func NewChecklist(recordID string, names []string) (*Checklist, error) {
	if recordID == "" {
		return nil, errors.New("record id is required")
	}
	items := make(map[string]Checkpoint)
	for _, name := range names {
		if name == "" {
			continue
		}
		items[name] = Checkpoint{Name: name, Required: true}
	}
	return &Checklist{RecordID: recordID, Items: items}, nil
}

func (c *Checklist) Complete(name string, seq int64) error {
	checkpoint, ok := c.Items[name]
	if !ok {
		return errors.New("checkpoint not found")
	}
	if seq <= 0 {
		return errors.New("checkpoint sequence is required")
	}
	if checkpoint.Completed && seq < checkpoint.Seq {
		return errors.New("checkpoint sequence cannot move backward")
	}
	checkpoint.Completed = true
	checkpoint.Seq = seq
	c.Items[name] = checkpoint
	return nil
}

func (c *Checklist) Pending() []string {
	pending := make([]string, 0)
	for name, checkpoint := range c.Items {
		if checkpoint.Required && !checkpoint.Completed {
			pending = append(pending, name)
		}
	}
	sort.Strings(pending)
	return pending
}

func (c *Checklist) Ready() bool {
	return len(c.Pending()) == 0 && len(c.Items) > 0
}

func BuildChecklist(record model.Record) *Checklist {
	names := []string{"content", "review", "confirmation"}
	checklist, _ := NewChecklist(record.ID, names)
	if record.Content != "" {
		_ = checklist.Complete("content", record.UpdatedSeq)
	}
	if record.Status == model.StatusApproved || record.Status == model.StatusArchived {
		_ = checklist.Complete("review", record.UpdatedSeq)
	}
	if record.PublishedSeq > 0 {
		_ = checklist.Complete("confirmation", record.PublishedSeq)
	}
	return checklist
}

func (c *Checklist) CompletedCount() int {
	count := 0
	for _, item := range c.Items {
		if item.Completed {
			count++
		}
	}
	return count
}

func (c *Checklist) RequiredCount() int {
	count := 0
	for _, item := range c.Items {
		if item.Required {
			count++
		}
	}
	return count
}

func (c *Checklist) Progress() int {
	required := c.RequiredCount()
	if required == 0 {
		return 0
	}
	return c.CompletedCount() * 100 / required
}

func (c *Checklist) MarkOptional(name string) error {
	item, ok := c.Items[name]
	if !ok {
		return errors.New("checkpoint not found")
	}
	item.Required = false
	c.Items[name] = item
	return nil
}
