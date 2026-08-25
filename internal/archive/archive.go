package archive

import (
	"errors"
	"fmt"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

type Service struct {
	catalog *catalog.Catalog
	store   *store.Store
}

func New(c *catalog.Catalog, s *store.Store) *Service {
	return &Service{catalog: c, store: s}
}

func (s *Service) Archive(id, actor string, seq int64) (model.Record, error) {
	if actor == "" {
		return model.Record{}, errors.New("actor is required")
	}
	r, err := s.catalog.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if r.Status != model.StatusApproved {
		return model.Record{}, errors.New("only approved records can be archived")
	}
	if seq <= r.UpdatedSeq {
		return model.Record{}, store.ErrConflict
	}
	r.Status = model.StatusArchived
	r.ArchivedSeq = seq
	r.UpdatedSeq = seq
	r.Version++
	if err := s.store.PutRecord(r); err != nil {
		return model.Record{}, err
	}
	if err := s.store.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s:%012d", r.ID, seq), RecordID: r.ID, Action: "archived", Actor: actor, Seq: seq}); err != nil {
		return model.Record{}, err
	}
	return r, nil
}

func (s *Service) Restore(id, actor string, seq int64) (model.Record, error) {
	r, err := s.catalog.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if r.Status != model.StatusArchived {
		return model.Record{}, errors.New("record is not archived")
	}
	if actor == "" || seq <= r.UpdatedSeq {
		return model.Record{}, errors.New("restore request is invalid")
	}
	r.Status = model.StatusApproved
	r.UpdatedSeq = seq
	r.Version++
	if err := s.store.PutRecord(r); err != nil {
		return model.Record{}, err
	}
	if err := s.store.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s:%012d", r.ID, seq), RecordID: r.ID, Action: "restored", Actor: actor, Seq: seq}); err != nil {
		return model.Record{}, err
	}
	return r, nil
}

func (s *Service) IsVisible(r model.Record) bool {
	return r.Status == model.StatusApproved || r.Status == model.StatusArchived
}
