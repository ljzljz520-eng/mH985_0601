package model

type RecordStatus string

const (
	StatusDraft    RecordStatus = "draft"
	StatusPending  RecordStatus = "pending"
	StatusApproved RecordStatus = "approved"
	StatusRejected RecordStatus = "rejected"
	StatusArchived RecordStatus = "archived"
)

type Record struct {
	ID           string       `json:"id"`
	StoreID      string       `json:"store_id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	Category     string       `json:"category"`
	Status       RecordStatus `json:"status"`
	Version      int          `json:"version"`
	SortKey      int          `json:"sort_key"`
	Owner        string       `json:"owner"`
	Reviewer     string       `json:"reviewer"`
	CreatedSeq   int64        `json:"created_seq"`
	UpdatedSeq   int64        `json:"updated_seq"`
	PublishedSeq int64        `json:"published_seq"`
	ArchivedSeq  int64        `json:"archived_seq"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Note     string `json:"note"`
	Seq      int64  `json:"seq"`
}

type Workflow struct {
	ID        string   `json:"id"`
	RecordID  string   `json:"record_id"`
	Name      string   `json:"name"`
	Stage     string   `json:"stage"`
	Steps     []string `json:"steps"`
	Completed []string `json:"completed"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Size     int    `json:"size"`
	Digest   string `json:"digest"`
}

type SearchQuery struct {
	StoreID  string       `json:"store_id"`
	Status   RecordStatus `json:"status"`
	Category string       `json:"category"`
	Text     string       `json:"text"`
	Limit    int          `json:"limit"`
}

type ChangeRequest struct {
	RecordID string `json:"record_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Actor    string `json:"actor"`
	Seq      int64  `json:"seq"`
}

type ReviewRequest struct {
	RecordID string `json:"record_id"`
	Reviewer string `json:"reviewer"`
	Approve  bool   `json:"approve"`
	Note     string `json:"note"`
	Seq      int64  `json:"seq"`
}

type ImportRow struct {
	ID       string `json:"id"`
	StoreID  string `json:"store_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Owner    string `json:"owner"`
	SortKey  int    `json:"sort_key"`
}

type ImportReport struct {
	Accepted int
	Rejected int
	Errors   []string
	IDs      []string
}

type PublishResult struct {
	Record   Record
	Workflow Workflow
}

type Detail struct {
	Record      Record
	Audit       []AuditEvent
	Attachments []Attachment
	Workflow    Workflow
}
