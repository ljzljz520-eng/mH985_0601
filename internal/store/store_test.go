package store

import (
	"path/filepath"
	"testing"

	"trainingdesk/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "training.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	r := model.Record{ID: "r-1", StoreID: "s-1", Title: "Safety", Content: "Steps", Status: model.StatusDraft, Version: 1}
	if err := s.PutRecord(r); err != nil {
		t.Fatal(err)
	}
	if err := s.PutAudit(model.AuditEvent{ID: "r-1:000000000001", RecordID: "r-1", Action: "created", Seq: 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutWorkflow(model.Workflow{ID: "w-1", RecordID: "r-1", Name: "onboarding", Stage: "draft"}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutAttachment(model.Attachment{ID: "a-1", RecordID: "r-1", Name: "guide.pdf"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	got, err := s.GetRecord("r-1")
	if err != nil || got.Title != r.Title {
		t.Fatalf("reopened record = %#v, err=%v", got, err)
	}
	audit, err := s.ListAudit("r-1")
	if err != nil || len(audit) != 1 {
		t.Fatalf("audit = %#v, err=%v", audit, err)
	}
	wf, err := s.GetWorkflow("w-1")
	if err != nil || wf.RecordID != "r-1" {
		t.Fatalf("workflow = %#v, err=%v", wf, err)
	}
	attachments, err := s.ListAttachments("r-1")
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments = %#v, err=%v", attachments, err)
	}
}

func TestStoreRejectsInvalidIdentity(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "x.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.PutRecord(model.Record{}); err == nil {
		t.Fatal("expected identity error")
	}
	if err := s.DeleteRecord(""); err == nil {
		t.Fatal("expected delete identity error")
	}
}
