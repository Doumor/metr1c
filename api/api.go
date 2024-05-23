package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type APISummary struct {
	Platform1CVersion  string `json:"platform1c_version"`
	SessionCount       int    `json:"sessions_total"`
	SessionsActive     int    `json:"sessions_active"`
	SessionsHybernated int    `json:"sessions_hybernated"`
	UsedLicensesSoft   int    `json:"licenses_used_soft"`
	UsedLicensesHASP   int    `json:"licenses_used_hasp"`
	ConnectionCount    int    `json:"connections_total"`
	ProcessCount       int    `json:"processes_total"`
	ProcessesMemoryKB  int    `json:"processes_mem_kb_total"`
}

func requestHandler(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type APIServer struct {
	mutex    sync.RWMutex
	summary  APISummary
	sessions []map[string]string
}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func (s *APIServer) ServeSummary(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	requestHandler(w, s.summary)
}

func (s *APIServer) UpdateSummary(update APISummary) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.summary = update
}

func (s *APIServer) ServeSessions(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	requestHandler(w, s.sessions)
}

func (s *APIServer) UpdateSessions(update []map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions = update
}
