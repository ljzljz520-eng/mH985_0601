package query

import (
	"sort"

	"trainingdesk/internal/model"
)

func Apply(records []model.Record, request Request) []model.Record {
	filtered := make([]model.Record, 0, len(records))
	for _, record := range records {
		if request.Match(record) {
			filtered = append(filtered, record)
		}
	}
	if request.SortField != "" {
		sort.SliceStable(filtered, func(i, j int) bool {
			left := fieldValue(request.SortField, filtered[i])
			right := fieldValue(request.SortField, filtered[j])
			if request.Descending {
				return left > right
			}
			return left < right
		})
	}
	start := request.Offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + request.Limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end]
}

func Summary(records []model.Record) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		counts[string(record.Status)]++
	}
	return counts
}
