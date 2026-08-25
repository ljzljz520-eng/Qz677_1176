package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"subsidy11/audit"
	"subsidy11/domain"
	"subsidy11/intake"
	"subsidy11/query"
	"subsidy11/review"
	"subsidy11/storage"
)

type server struct {
	store   *storage.Store
	imports *intake.Importer
	review  *review.Service
	query   *query.Service
	audit   *audit.Service
}

func main() {
	path := os.Getenv("SUBSIDY_DB")
	if path == "" {
		path = "./subsidy11.db"
	}
	store, err := storage.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	s := &server{store: store, imports: intake.NewImporter(), review: review.NewService(store), query: query.NewService(store), audit: audit.NewService(store)}
	address := os.Getenv("SUBSIDY_ADDR")
	if address == "" {
		address = ":8080"
	}
	log.Printf("subsidy11 listening on %s", address)
	log.Fatal(http.ListenAndServe(address, s.routes()))
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	mux.HandleFunc("/summary", s.summary)
	return logging(mux)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started))
	})
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		filter := query.Filter{Region: r.URL.Query().Get("region"), Status: r.URL.Query().Get("status"), CitizenID: r.URL.Query().Get("citizen_id")}
		items, err := s.query.Search(filter)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var record domain.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		writeError(w, err)
		return
	}
	record = domain.NormalizeRecord(record)
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now().UTC()
	}
	record.UpdatedAt = record.CreatedAt
	if err := intake.ImportOne(s.store, record); err != nil {
		writeError(w, err)
		return
	}
	_ = s.audit.Record(record.ID, "import", "api", "record received")
	writeJSON(w, http.StatusCreated, record)
}

func (s *server) record(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[1] == "" {
		writeErrorStatus(w, http.StatusBadRequest, fmt.Errorf("record id required"))
		return
	}
	id := parts[1]
	if len(parts) == 2 && r.Method == http.MethodGet {
		record, err := s.query.Get(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
		return
	}
	if len(parts) == 3 && parts[2] == "review" && r.Method == http.MethodPost {
		var input struct {
			OfficerID string `json:"officer_id"`
			Note      string `json:"note"`
			Approved  bool   `json:"approved"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, err)
			return
		}
		if err := s.review.Confirm(id, input.OfficerID, input.Note, input.Approved); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
		return
	}
	if len(parts) == 3 && parts[2] == "events" && r.Method == http.MethodGet {
		events, err := s.query.Timeline(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, events)
		return
	}
	writeErrorStatus(w, http.StatusNotFound, domain.ErrNotFound)
}

func (s *server) summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	items, err := s.query.Search(query.Filter{})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, query.BuildSummary(items))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) { writeErrorStatus(w, http.StatusBadRequest, err) }

func writeErrorStatus(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
