package catalog

import (
	"sort"
	"strings"

	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

func (c *Catalog) Search(q model.SearchQuery) ([]model.Record, error) {
	all, err := c.allRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]model.Record, 0, len(all))
	needle := strings.ToLower(strings.TrimSpace(q.Text))
	for _, r := range all {
		if q.StoreID != "" && r.StoreID != q.StoreID {
			continue
		}
		if q.Status != "" && r.Status != q.Status {
			continue
		}
		if q.Category != "" && r.Category != q.Category {
			continue
		}
		if needle != "" && !strings.Contains(strings.ToLower(r.Title), needle) && !strings.Contains(strings.ToLower(r.Content), needle) {
			continue
		}
		filtered = append(filtered, r)
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].SortKey == filtered[j].SortKey {
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].SortKey < filtered[j].SortKey
	})
	if q.Limit > 0 && len(filtered) > q.Limit {
		filtered = filtered[:q.Limit]
	}
	return filtered, nil
}

func (c *Catalog) allRecords() ([]model.Record, error) {
	return c.storage.ScanRecords()
}

func (c *Catalog) Detail(id string) (model.Detail, error) {
	r, err := c.Get(id)
	if err != nil {
		return model.Detail{}, err
	}
	audit, err := c.storage.ListAudit(id)
	if err != nil {
		return model.Detail{}, err
	}
	attachments, err := c.storage.ListAttachments(id)
	if err != nil {
		return model.Detail{}, err
	}
	wf, err := c.storage.FindWorkflowForRecord(id)
	if err != nil && err != store.ErrNotFound {
		return model.Detail{}, err
	}
	return model.Detail{Record: r, Audit: audit, Attachments: attachments, Workflow: wf}, nil
}
