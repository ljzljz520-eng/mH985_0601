package catalog

import (
	"errors"
	"fmt"
	"strings"

	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

type Catalog struct {
	storage *store.Store
}

func New(storage *store.Store) *Catalog {
	return &Catalog{storage: storage}
}

func (c *Catalog) Register(row model.ImportRow, seq int64) (model.Record, error) {
	if strings.TrimSpace(row.ID) == "" {
		return model.Record{}, errors.New("id is required")
	}
	if strings.TrimSpace(row.StoreID) == "" {
		return model.Record{}, errors.New("store id is required")
	}
	if strings.TrimSpace(row.Title) == "" || strings.TrimSpace(row.Content) == "" {
		return model.Record{}, errors.New("title and content are required")
	}
	if row.SortKey < 0 {
		return model.Record{}, errors.New("sort key cannot be negative")
	}
	r := model.Record{ID: row.ID, StoreID: row.StoreID, Title: strings.TrimSpace(row.Title), Content: row.Content, Category: strings.TrimSpace(row.Category), Status: model.StatusDraft, Version: 1, SortKey: row.SortKey, Owner: row.Owner, CreatedSeq: seq, UpdatedSeq: seq}
	if err := c.storage.PutRecord(r); err != nil {
		return model.Record{}, err
	}
	if err := c.storage.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s:%012d", r.ID, seq), RecordID: r.ID, Action: "registered", Actor: r.Owner, Seq: seq}); err != nil {
		return model.Record{}, err
	}
	return r, nil
}

func (c *Catalog) Get(id string) (model.Record, error) {
	if strings.TrimSpace(id) == "" {
		return model.Record{}, errors.New("id is required")
	}
	return c.storage.GetRecord(id)
}

func (c *Catalog) Change(req model.ChangeRequest) (model.Record, error) {
	r, err := c.Get(req.RecordID)
	if err != nil {
		return model.Record{}, err
	}
	if r.Status == model.StatusArchived {
		return model.Record{}, errors.New("archived records cannot change")
	}
	if req.Seq <= r.UpdatedSeq {
		return model.Record{}, store.ErrConflict
	}
	if strings.TrimSpace(req.Title) != "" {
		r.Title = strings.TrimSpace(req.Title)
	}
	if req.Content != "" {
		r.Content = req.Content
	}
	if strings.TrimSpace(req.Category) != "" {
		r.Category = strings.TrimSpace(req.Category)
	}
	r.Version++
	r.Status = model.StatusPending
	r.UpdatedSeq = req.Seq
	if err := c.storage.PutRecord(r); err != nil {
		return model.Record{}, err
	}
	if err := c.storage.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s:%012d", r.ID, req.Seq), RecordID: r.ID, Action: "changed", Actor: req.Actor, Seq: req.Seq}); err != nil {
		return model.Record{}, err
	}
	return r, nil
}

func (c *Catalog) Attach(a model.Attachment) error {
	if a.Size < 0 {
		return errors.New("attachment size cannot be negative")
	}
	if a.Digest == "" {
		return errors.New("attachment digest is required")
	}
	return c.storage.PutAttachment(a)
}
