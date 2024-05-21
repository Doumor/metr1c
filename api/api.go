package api

import (
	"encoding/json"
	"log"
	"net/http"
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

func ServeJSON(data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Fatalf("serving jsonified data via an API handle '%s': %s", r.URL.Path, err)
		}
	}
}
