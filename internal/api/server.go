package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"trainingdesk/internal/archive"
	"trainingdesk/internal/catalog"
	"trainingdesk/internal/flow008"
	"trainingdesk/internal/importer"
	"trainingdesk/internal/model"
	"trainingdesk/internal/review"
	"trainingdesk/internal/store"
	"trainingdesk/internal/validation"
	"trainingdesk/internal/workflow"
)

type Server struct {
	store    *store.Store
	catalog  *catalog.Catalog
	review   *review.Service
	archive  *archive.Service
	importer *importer.Service
	flow     *flow008.Handler
	workflow *workflow.Engine
}

func New(s *store.Store) *Server {
	c := catalog.New(s)
	rv := review.New(c, s)
	return &Server{store: s, catalog: c, review: rv, archive: archive.New(c, s), importer: importer.New(c), flow: flow008.New(c, rv), workflow: workflow.New(s)}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	mux.HandleFunc("/report", s.summary)
	mux.HandleFunc("/analytics", s.analytics)
	mux.HandleFunc("/export", s.export)
	mux.HandleFunc("/workflows", s.workflows)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := model.SearchQuery{StoreID: r.URL.Query().Get("store"), Text: r.URL.Query().Get("q"), Category: r.URL.Query().Get("category")}
		if err := validation.SearchInput(q); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		items, err := s.catalog.Search(q)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if r.Method == http.MethodPost {
		var row model.ImportRow
		if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validation.RecordInput(row); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		created, err := s.catalog.Register(row, int64(row.SortKey+1))
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeError(w, http.StatusBadRequest, store.ErrNotFound)
		return
	}
	if r.Method == http.MethodGet {
		detail, err := s.catalog.Detail(id)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, detail)
		return
	}
	if r.Method == http.MethodPatch {
		var req model.ChangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		req.RecordID = id
		changed, err := s.catalog.Change(req)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, changed)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
