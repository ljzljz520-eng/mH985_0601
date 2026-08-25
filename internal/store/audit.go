package store

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
	"trainingdesk/internal/model"
)

func auditKey(recordID string, seq int64) string {
	return fmt.Sprintf("%s:%012d", recordID, seq)
}

func (s *Store) PutAudit(e model.AuditEvent) error {
	if e.RecordID == "" || e.ID == "" {
		return errors.New("audit identity is required")
	}
	return s.put([]byte("audits"), e.ID, e)
}

func (s *Store) ListAudit(recordID string) ([]model.AuditEvent, error) {
	if recordID == "" {
		return nil, errors.New("record id is required")
	}
	items := make([]model.AuditEvent, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte("audits")).Cursor()
		prefix := []byte(recordID + ":")
		for k, v := c.Seek(prefix); k != nil && len(k) >= len(prefix) && string(k[:len(prefix)]) == string(prefix); k, v = c.Next() {
			var e model.AuditEvent
			if err := decode(v, &e); err != nil {
				return err
			}
			items = append(items, e)
		}
		return nil
	})
	return items, err
}

func (s *Store) PutWorkflow(w model.Workflow) error {
	if w.ID == "" || w.RecordID == "" {
		return errors.New("workflow identity is required")
	}
	return s.put([]byte("workflows"), w.ID, w)
}

func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	var w model.Workflow
	err := s.get([]byte("workflows"), id, &w)
	return w, err
}

func (s *Store) PutAttachment(a model.Attachment) error {
	if a.ID == "" || a.RecordID == "" {
		return errors.New("attachment identity is required")
	}
	return s.put([]byte("attachments"), a.ID, a)
}

func (s *Store) ListAttachments(recordID string) ([]model.Attachment, error) {
	if recordID == "" {
		return nil, errors.New("record id is required")
	}
	items := make([]model.Attachment, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("attachments")).ForEach(func(_, v []byte) error {
			var a model.Attachment
			if err := decode(v, &a); err != nil {
				return err
			}
			if a.RecordID == recordID {
				items = append(items, a)
			}
			return nil
		})
	})
	return items, err
}
