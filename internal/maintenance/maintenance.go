package maintenance

import (
	"errors"
	"sort"

	"trainingdesk/internal/model"
)

type Check struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Detail  string `json:"detail"`
}

type Snapshot struct {
	Records      int
	Audited      int
	WithContent  int
	WithReviewer int
	Statuses     map[model.RecordStatus]int
}

func SnapshotRecords(records []model.Record, audited map[string]bool) Snapshot {
	snapshot := Snapshot{Statuses: make(map[model.RecordStatus]int)}
	for _, record := range records {
		snapshot.Records++
		snapshot.Statuses[record.Status]++
		if audited[record.ID] {
			snapshot.Audited++
		}
		if record.Content != "" {
			snapshot.WithContent++
		}
		if record.Reviewer != "" {
			snapshot.WithReviewer++
		}
	}
	return snapshot
}

func CheckRecords(records []model.Record) []Check {
	checks := []Check{
		checkIDs(records),
		checkContent(records),
		checkStatuses(records),
	}
	return checks
}

func checkIDs(records []model.Record) Check {
	seen := make(map[string]bool)
	for _, record := range records {
		if record.ID == "" {
			return Check{Name: "ids", Detail: "empty record id"}
		}
		if seen[record.ID] {
			return Check{Name: "ids", Detail: "duplicate record id"}
		}
		seen[record.ID] = true
	}
	return Check{Name: "ids", Healthy: true, Detail: "record ids are unique"}
}

func checkContent(records []model.Record) Check {
	for _, record := range records {
		if record.Status == model.StatusApproved && record.Content == "" {
			return Check{Name: "content", Detail: "approved record has no content"}
		}
	}
	return Check{Name: "content", Healthy: true, Detail: "approved records have content"}
}

func checkStatuses(records []model.Record) Check {
	for _, record := range records {
		switch record.Status {
		case model.StatusDraft, model.StatusPending, model.StatusApproved, model.StatusRejected, model.StatusArchived:
		default:
			return Check{Name: "statuses", Detail: "unknown record status"}
		}
	}
	return Check{Name: "statuses", Healthy: true, Detail: "statuses are recognized"}
}

func Healthy(checks []Check) bool {
	if len(checks) == 0 {
		return false
	}
	for _, check := range checks {
		if !check.Healthy {
			return false
		}
	}
	return true
}

func Normalize(records []model.Record) ([]model.Record, error) {
	copyRecords := append([]model.Record(nil), records...)
	for i := range copyRecords {
		if copyRecords[i].ID == "" {
			return nil, errors.New("cannot normalize empty id")
		}
		copyRecords[i].Title = trim(copyRecords[i].Title)
		copyRecords[i].Category = trim(copyRecords[i].Category)
		copyRecords[i].Owner = trim(copyRecords[i].Owner)
	}
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].ID < copyRecords[j].ID })
	return copyRecords, nil
}

func trim(value string) string {
	left := 0
	right := len(value)
	for left < right && value[left] <= ' ' {
		left++
	}
	for right > left && value[right-1] <= ' ' {
		right--
	}
	return value[left:right]
}
