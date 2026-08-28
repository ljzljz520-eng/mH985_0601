package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"trainingdesk/internal/store"
)

func TestServerHealthAndRecords(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	h := New(s).Handler()
	health := httptest.NewRecorder()
	h.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), "ok") {
		t.Fatalf("health=%d %s", health.Code, health.Body.String())
	}
	req := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(`{"id":"r","store_id":"s","title":"Guide","content":"Body"}`))
	created := httptest.NewRecorder()
	h.ServeHTTP(created, req)
	if created.Code != http.StatusCreated {
		t.Fatalf("created=%d %s", created.Code, created.Body.String())
	}
	got := httptest.NewRecorder()
	h.ServeHTTP(got, httptest.NewRequest(http.MethodGet, "/records/r", nil))
	if got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "Guide") {
		t.Fatalf("record=%d %s", got.Code, got.Body.String())
	}
}
