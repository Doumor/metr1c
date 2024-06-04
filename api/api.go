package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Summary struct {
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
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("api error: encoding data to JSON: %s", err)
	}
}

type Server struct {
	mutex       sync.RWMutex
	summary     Summary
	sessions    []map[string]string
	connections []map[string]string
	processes   []map[string]string
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ServeSummary(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	requestHandler(w, s.summary)
}

func (s *Server) UpdateSummary(update Summary) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.summary = update
}

func (s *Server) ServeSessions(w http.ResponseWriter, _ *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	requestHandler(w, s.sessions)
}

func (s *Server) UpdateSessions(update []map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions = update
}

func (s *Server) ServeConnections(w http.ResponseWriter, _ *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	requestHandler(w, s.connections)
}

func (s *Server) UpdateConnections(update []map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.connections = update
}

func (s *Server) ServeProcesses(w http.ResponseWriter, _ *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	requestHandler(w, s.processes)
}

func (s *Server) UpdateProcesses(update []map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.processes = update
}
