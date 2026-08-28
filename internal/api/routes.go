package api

import (
	"encoding/json"
	"net/http"

	"trainingdesk/internal/analytics"
	"trainingdesk/internal/exporter"

	"trainingdesk/internal/model"
	"trainingdesk/internal/report"
)

func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	items, err := s.catalog.Search(model.SearchQuery{Limit: 500})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, report.Build(items))
}

func (s *Server) analytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	items, err := s.catalog.Search(model.SearchQuery{Limit: 500})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, analytics.Build(items))
}

func (s *Server) export(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	items, err := s.catalog.Search(model.SearchQuery{Limit: 500})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	format := r.URL.Query().Get("format")
	if format == "json" {
		data, err := exporter.JSON(items, exporter.Options{IncludeContent: r.URL.Query().Get("content") == "1"})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	if err := exporter.CSV(w, items, exporter.Options{IncludeContent: r.URL.Query().Get("content") == "1"}); err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}

func (s *Server) workflows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		RecordID string `json:"record_id"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	wf, err := s.workflow.Start(req.RecordID, req.Name)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, wf)
}
