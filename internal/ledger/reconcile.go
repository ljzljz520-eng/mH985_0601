package ledger

import (
	"fmt"

	"trainingdesk/internal/model"
)

type Reconciliation struct {
	RecordID       string
	ExpectedStatus model.RecordStatus
	LedgerStatus   string
	Balanced       bool
	Message        string
}

func (l *Ledger) Reconcile(record model.Record) Reconciliation {
	result := Reconciliation{RecordID: record.ID, ExpectedStatus: record.Status, Balanced: true}
	latest, err := l.Latest(record.ID)
	if err != nil {
		if record.Status == model.StatusDraft {
			result.LedgerStatus = string(model.StatusDraft)
			result.Message = "draft has no transitions"
			return result
		}
		result.Balanced = false
		result.Message = "record has no ledger transition"
		return result
	}
	result.LedgerStatus = latest.To
	if latest.To != string(record.Status) {
		result.Balanced = false
		result.Message = fmt.Sprintf("record is %s but ledger is %s", record.Status, latest.To)
		return result
	}
	result.Message = "record and ledger agree"
	return result
}

func (l *Ledger) TransitionPath(recordID string) []string {
	path := make([]string, 0)
	for _, entry := range l.ForRecord(recordID) {
		if len(path) == 0 {
			path = append(path, entry.From)
		}
		path = append(path, entry.To)
	}
	return path
}

func ValidTransition(from, to model.RecordStatus) bool {
	switch from {
	case model.StatusDraft:
		return to == model.StatusPending
	case model.StatusPending:
		return to == model.StatusApproved || to == model.StatusRejected
	case model.StatusRejected:
		return to == model.StatusPending
	case model.StatusApproved:
		return to == model.StatusArchived
	case model.StatusArchived:
		return to == model.StatusApproved
	default:
		return false
	}
}
