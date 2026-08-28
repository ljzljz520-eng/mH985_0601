package validation

import (
	"errors"
	"strings"

	"trainingdesk/internal/model"
)

func RecordInput(row model.ImportRow) error {
	if strings.TrimSpace(row.ID) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(row.StoreID) == "" {
		return errors.New("store id is required")
	}
	if strings.TrimSpace(row.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(row.Content) == "" {
		return errors.New("content is required")
	}
	return nil
}

func SearchInput(q model.SearchQuery) error {
	if q.Limit < 0 {
		return errors.New("limit cannot be negative")
	}
	if q.Limit > 500 {
		return errors.New("limit is too large")
	}
	return nil
}
