package validation

import (
	"errors"
	"strings"

	"trainingdesk/internal/model"
)

func CanChange(r model.Record, actor string) error {
	if strings.TrimSpace(actor) == "" {
		return errors.New("actor is required")
	}
	if r.Status == model.StatusArchived {
		return errors.New("archived record is immutable")
	}
	if r.Status == model.StatusApproved && actor != r.Owner && actor != r.Reviewer {
		return errors.New("approved record requires owner or reviewer")
	}
	return nil
}

func CanPublish(r model.Record) bool {
	return r.Status == model.StatusApproved && r.Content != ""
}
