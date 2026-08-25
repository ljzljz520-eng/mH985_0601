package analytics

import (
	"sort"
	"strings"

	"trainingdesk/internal/model"
)

type Dashboard struct {
	TotalRecords      int
	VisibleRecords    int
	PendingRecords    int
	ArchivedRecords   int
	CompletionRate    int
	StoreBreakdown    []StoreMetric
	CategoryBreakdown []CategoryMetric
	ReviewerBreakdown []ReviewerMetric
}

type StoreMetric struct {
	StoreID string `json:"store_id"`
	Total   int    `json:"total"`
	Visible int    `json:"visible"`
}

type CategoryMetric struct {
	Category string `json:"category"`
	Total    int    `json:"total"`
}

type ReviewerMetric struct {
	Reviewer string `json:"reviewer"`
	Approved int    `json:"approved"`
	Rejected int    `json:"rejected"`
}

func Build(records []model.Record) Dashboard {
	dashboard := Dashboard{}
	stores := make(map[string]*StoreMetric)
	categories := make(map[string]*CategoryMetric)
	reviewers := make(map[string]*ReviewerMetric)
	for _, r := range records {
		dashboard.TotalRecords++
		storeMetric := stores[r.StoreID]
		if storeMetric == nil {
			storeMetric = &StoreMetric{StoreID: r.StoreID}
			stores[r.StoreID] = storeMetric
		}
		storeMetric.Total++
		if r.Status == model.StatusApproved || r.Status == model.StatusArchived {
			storeMetric.Visible++
			dashboard.VisibleRecords++
		}
		if r.Status == model.StatusPending {
			dashboard.PendingRecords++
		}
		if r.Status == model.StatusArchived {
			dashboard.ArchivedRecords++
		}
		categoryMetric := categories[r.Category]
		if categoryMetric == nil {
			categoryMetric = &CategoryMetric{Category: r.Category}
			categories[r.Category] = categoryMetric
		}
		categoryMetric.Total++
		if r.Reviewer != "" {
			reviewerMetric := reviewers[r.Reviewer]
			if reviewerMetric == nil {
				reviewerMetric = &ReviewerMetric{Reviewer: r.Reviewer}
				reviewers[r.Reviewer] = reviewerMetric
			}
			if r.Status == model.StatusApproved || r.Status == model.StatusArchived {
				reviewerMetric.Approved++
			}
			if r.Status == model.StatusRejected {
				reviewerMetric.Rejected++
			}
		}
	}
	if dashboard.TotalRecords > 0 {
		dashboard.CompletionRate = dashboard.VisibleRecords * 100 / dashboard.TotalRecords
	}
	for _, metric := range stores {
		dashboard.StoreBreakdown = append(dashboard.StoreBreakdown, *metric)
	}
	for _, metric := range categories {
		dashboard.CategoryBreakdown = append(dashboard.CategoryBreakdown, *metric)
	}
	for _, metric := range reviewers {
		dashboard.ReviewerBreakdown = append(dashboard.ReviewerBreakdown, *metric)
	}
	sort.Slice(dashboard.StoreBreakdown, func(i, j int) bool { return dashboard.StoreBreakdown[i].StoreID < dashboard.StoreBreakdown[j].StoreID })
	sort.Slice(dashboard.CategoryBreakdown, func(i, j int) bool {
		return dashboard.CategoryBreakdown[i].Category < dashboard.CategoryBreakdown[j].Category
	})
	sort.Slice(dashboard.ReviewerBreakdown, func(i, j int) bool {
		return dashboard.ReviewerBreakdown[i].Reviewer < dashboard.ReviewerBreakdown[j].Reviewer
	})
	return dashboard
}

func Headline(d Dashboard) string {
	parts := []string{
		"records=" + itoa(d.TotalRecords),
		"visible=" + itoa(d.VisibleRecords),
		"pending=" + itoa(d.PendingRecords),
		"completion=" + itoa(d.CompletionRate) + "%",
	}
	return strings.Join(parts, " ")
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	buf := make([]byte, 0, 12)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	if negative {
		return "-" + string(buf)
	}
	return string(buf)
}
