package review

import (
	"errors"
	"fmt"
	"strings"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/notification"
	"trainingdesk/internal/store"
)

type Service struct {
	catalog *catalog.Catalog
	store   *store.Store
	outbox  *notification.Outbox
}

func New(c *catalog.Catalog, s *store.Store) *Service {
	return &Service{catalog: c, store: s, outbox: notification.New()}
}

func (s *Service) Submit(id string, seq int64) (model.Record, error) {
	r, err := s.catalog.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if r.Status != model.StatusDraft && r.Status != model.StatusRejected {
		return model.Record{}, errors.New("only draft or rejected records can be submitted")
	}
	r.Status = model.StatusPending
	r.Version++
	r.UpdatedSeq = seq
	if err := s.store.PutRecord(r); err != nil {
		return model.Record{}, err
	}
	if err := s.store.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s:%012d", r.ID, seq), RecordID: r.ID, Action: "submitted", Seq: seq}); err != nil {
		return model.Record{}, err
	}
	return r, nil
}

func (s *Service) Decide(req model.ReviewRequest) (model.Record, error) {
	if strings.TrimSpace(req.Reviewer) == "" {
		return model.Record{}, errors.New("reviewer is required")
	}
	r, err := s.catalog.Get(req.RecordID)
	if err != nil {
		return model.Record{}, err
	}
	if r.Status != model.StatusPending {
		return model.Record{}, errors.New("record is not pending")
	}
	if req.Seq <= r.UpdatedSeq {
		return model.Record{}, store.ErrConflict
	}
	r.Reviewer = req.Reviewer
	r.UpdatedSeq = req.Seq
	r.Version++
	action := "rejected"
	if req.Approve {
		r.Status = model.StatusApproved
		r.PublishedSeq = req.Seq
		action = "approved"
	}
	if req.Note == "" {
		req.Note = action
	}
	if err := s.store.PutRecord(r); err != nil {
		return model.Record{}, err
	}
	if err := s.store.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s:%012d", r.ID, req.Seq), RecordID: r.ID, Action: action, Actor: req.Reviewer, Note: req.Note, Seq: req.Seq}); err != nil {
		return model.Record{}, err
	}
	target := r.Owner
	if target == "" {
		target = req.Reviewer
	}
	if _, err := s.outbox.Enqueue(r, action, target, req.Seq); err != nil {
		return model.Record{}, err
	}
	return r, nil
}

func (s *Service) CanReview(r model.Record) bool {
	return r.Status == model.StatusPending && r.Reviewer == ""
}

func (s *Service) Pending(candidates []model.Record) []model.Record {
	items := make([]model.Record, 0, len(candidates))
	for _, r := range candidates {
		if r.Status == model.StatusPending {
			items = append(items, r)
		}
	}
	return items
}

func (s *Service) Notifications() []notification.Message {
	return s.outbox.Pending()
}
