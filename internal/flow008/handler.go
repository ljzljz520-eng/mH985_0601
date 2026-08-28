package flow008

import (
	"errors"
	"sort"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
	"trainingdesk/internal/review"
)

type Handler struct {
	catalog *catalog.Catalog
	review  *review.Service
}

func New(c *catalog.Catalog, r *review.Service) *Handler {
	return &Handler{catalog: c, review: r}
}

func (h *Handler) Ordered(storeID string) ([]model.Record, error) {
	items, err := h.catalog.Search(model.SearchQuery{StoreID: storeID})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].SortKey == items[j].SortKey {
			return items[i].ID > items[j].ID
		}
		return items[i].SortKey < items[j].SortKey
	})
	return items, nil
}

func (h *Handler) DetailAt(storeID string, index int) (model.Record, error) {
	if index < 0 {
		return model.Record{}, errors.New("index cannot be negative")
	}
	items, err := h.Ordered(storeID)
	if err != nil {
		return model.Record{}, err
	}
	if index >= len(items) {
		return model.Record{}, errors.New("record index out of range")
	}
	return h.catalog.Get(items[index].ID)
}

func (h *Handler) ReviewAt(storeID string, index int, reviewer string, approve bool, seq int64) (model.Record, error) {
	r, err := h.DetailAt(storeID, index)
	if err != nil {
		return model.Record{}, err
	}
	return h.review.Decide(model.ReviewRequest{RecordID: r.ID, Reviewer: reviewer, Approve: approve, Seq: seq})
}

func (h *Handler) VisibleCount(storeID string) (int, error) {
	items, err := h.Ordered(storeID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range items {
		if item.Status == model.StatusApproved || item.Status == model.StatusPending {
			count++
		}
	}
	return count, nil
}
