package catalog

import (
	"sort"

	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

func (c *Catalog) ListByStore(storeID string) ([]model.Record, error) {
	items, err := c.Search(model.SearchQuery{StoreID: storeID})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Catalog) CountByStatus(status model.RecordStatus) (int, error) {
	items, err := c.Search(model.SearchQuery{Status: status})
	if err != nil {
		return 0, err
	}
	return len(items), nil
}

func (c *Catalog) SortForExport(items []model.Record) []model.Record {
	copyItems := append([]model.Record(nil), items...)
	sort.SliceStable(copyItems, func(i, j int) bool {
		if copyItems[i].Category == copyItems[j].Category {
			return copyItems[i].Title < copyItems[j].Title
		}
		return copyItems[i].Category < copyItems[j].Category
	})
	return copyItems
}

func (c *Catalog) Remove(id string) error {
	if id == "" {
		return store.ErrNotFound
	}
	_, err := c.Get(id)
	if err != nil {
		return err
	}
	return c.storage.DeleteRecord(id)
}
