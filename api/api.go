package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type apiSummary struct {
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

type APIServer struct {
	mutex   sync.RWMutex
	Summary apiSummary
}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func RequestHandler(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *APIServer) ServeSummary(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	RequestHandler(w, s.Summary)
}
