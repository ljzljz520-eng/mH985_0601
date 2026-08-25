package report

import (
	"fmt"
	"sort"
	"strings"

	"trainingdesk/internal/model"
)

type Summary struct {
	Total      int
	Draft      int
	Pending    int
	Approved   int
	Rejected   int
	Archived   int
	Stores     int
	ByStore    map[string]int
	ByCategory map[string]int
}

func Build(records []model.Record) Summary {
	summary := Summary{ByStore: make(map[string]int), ByCategory: make(map[string]int)}
	for _, r := range records {
		summary.Total++
		summary.ByStore[r.StoreID]++
		summary.ByCategory[r.Category]++
		switch r.Status {
		case model.StatusDraft:
			summary.Draft++
		case model.StatusPending:
			summary.Pending++
		case model.StatusApproved:
			summary.Approved++
		case model.StatusRejected:
			summary.Rejected++
		case model.StatusArchived:
			summary.Archived++
		}
	}
	summary.Stores = len(summary.ByStore)
	return summary
}

func Render(summary Summary) string {
	stores := make([]string, 0, len(summary.ByStore))
	for storeID, count := range summary.ByStore {
		stores = append(stores, fmt.Sprintf("%s=%d", storeID, count))
	}
	sort.Strings(stores)
	return strings.Join([]string{
		fmt.Sprintf("total=%d", summary.Total),
		fmt.Sprintf("draft=%d", summary.Draft),
		fmt.Sprintf("pending=%d", summary.Pending),
		fmt.Sprintf("approved=%d", summary.Approved),
		fmt.Sprintf("rejected=%d", summary.Rejected),
		fmt.Sprintf("archived=%d", summary.Archived),
		fmt.Sprintf("stores=%s", strings.Join(stores, ",")),
	}, " ")
}

func ExportRows(records []model.Record) []string {
	rows := make([]string, 0, len(records))
	for _, r := range records {
		rows = append(rows, fmt.Sprintf("%s|%s|%s|%s|%s", r.ID, r.StoreID, r.Title, r.Category, r.Status))
	}
	sort.Strings(rows)
	return rows
}
