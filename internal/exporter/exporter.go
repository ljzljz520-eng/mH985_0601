package exporter

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strings"

	"trainingdesk/internal/model"
)

type Options struct {
	IncludeContent bool
	Statuses       []model.RecordStatus
	StoreID        string
}

func Filter(records []model.Record, options Options) []model.Record {
	allowed := make(map[model.RecordStatus]bool)
	for _, status := range options.Statuses {
		allowed[status] = true
	}
	filtered := make([]model.Record, 0, len(records))
	for _, r := range records {
		if options.StoreID != "" && r.StoreID != options.StoreID {
			continue
		}
		if len(allowed) > 0 && !allowed[r.Status] {
			continue
		}
		filtered = append(filtered, r)
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].StoreID == filtered[j].StoreID {
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].StoreID < filtered[j].StoreID
	})
	return filtered
}

func CSV(w io.Writer, records []model.Record, options Options) error {
	if w == nil {
		return errors.New("writer is required")
	}
	writer := csv.NewWriter(w)
	header := []string{"id", "store_id", "title", "category", "status", "version"}
	if options.IncludeContent {
		header = append(header, "content")
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, r := range Filter(records, options) {
		row := []string{r.ID, r.StoreID, r.Title, r.Category, string(r.Status), itoa(r.Version)}
		if options.IncludeContent {
			row = append(row, r.Content)
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func JSON(records []model.Record, options Options) ([]byte, error) {
	type safeRecord struct {
		ID       string             `json:"id"`
		StoreID  string             `json:"store_id"`
		Title    string             `json:"title"`
		Category string             `json:"category"`
		Status   model.RecordStatus `json:"status"`
		Version  int                `json:"version"`
		Content  string             `json:"content,omitempty"`
	}
	items := make([]safeRecord, 0)
	for _, r := range Filter(records, options) {
		item := safeRecord{ID: r.ID, StoreID: r.StoreID, Title: r.Title, Category: r.Category, Status: r.Status, Version: r.Version}
		if options.IncludeContent {
			item.Content = r.Content
		}
		items = append(items, item)
	}
	return json.Marshal(items)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 8)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

func Lines(records []model.Record, options Options) string {
	parts := make([]string, 0)
	for _, r := range Filter(records, options) {
		parts = append(parts, strings.Join([]string{r.ID, r.StoreID, r.Title, string(r.Status)}, "|"))
	}
	return strings.Join(parts, "\n")
}
