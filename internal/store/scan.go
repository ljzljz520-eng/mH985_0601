package store

import (
	"sort"

	"go.etcd.io/bbolt"
	"trainingdesk/internal/model"
)

func (s *Store) ScanRecords() ([]model.Record, error) {
	items := make([]model.Record, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r model.Record
			if err := decode(v, &r); err != nil {
				return err
			}
			items = append(items, r)
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, err
}

func (s *Store) FindWorkflowForRecord(recordID string) (model.Workflow, error) {
	var found model.Workflow
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("workflows")).ForEach(func(_, v []byte) error {
			var w model.Workflow
			if err := decode(v, &w); err != nil {
				return err
			}
			if w.RecordID == recordID {
				found = w
			}
			return nil
		})
	})
	if err != nil {
		return model.Workflow{}, err
	}
	if found.ID == "" {
		return model.Workflow{}, ErrNotFound
	}
	return found, nil
}
