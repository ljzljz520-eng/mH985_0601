package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"go.etcd.io/bbolt"
	"trainingdesk/internal/model"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record version conflict")
)

var bucketNames = [][]byte{
	[]byte("records"), []byte("audits"), []byte("workflows"), []byte("attachments"),
}

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bbolt: %w", err)
	}
	s := &Store{db: db}
	if err := s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func encode(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode value: %w", err)
	}
	return b, nil
}

func decode(data []byte, dst any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode value: %w", err)
	}
	return nil
}

func (s *Store) put(bucket []byte, key string, value any) error {
	b, err := encode(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(bucket).Put([]byte(key), b); err != nil {
			return fmt.Errorf("put %s: %w", key, err)
		}
		return nil
	})
}

func (s *Store) get(bucket []byte, key string, dst any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucket).Get([]byte(key))
		if data == nil {
			return ErrNotFound
		}
		copyData := append([]byte(nil), data...)
		return decode(copyData, dst)
	})
}

func (s *Store) DeleteRecord(id string) error {
	if id == "" {
		return errors.New("record id is required")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket([]byte("records")).Delete([]byte(id)); err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) PutRecord(r model.Record) error {
	if r.ID == "" {
		return errors.New("record id is required")
	}
	return s.put([]byte("records"), r.ID, r)
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	var r model.Record
	err := s.get([]byte("records"), id, &r)
	return r, err
}
