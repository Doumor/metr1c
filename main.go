/*
	metr1c
	Copyright (C) 2023 Doumor (doumor@vk.com)

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"doumor/metr1c/api"
	"doumor/metr1c/rac"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Version string

func getRecords(query rac.Query, cmd, subcmd, option string) rac.Query {
	query.Command = cmd
	query.SubCommand = subcmd
	query.Option = option

	err := query.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = query.Parse()
	if err != nil {
		log.Fatal(err)
	}

	return query
}

func countSessionTypes(sessions rac.Query) (float64, float64) {
	var active, hibernated int
	for _, session := range sessions.Records {
		switch session["hibernate"] {
		case "no":
			active++
		case "yes":
			hibernated++
		default:
			log.Printf("'rac session list': unexpected 'hibernate' field value: '%s'", session["hibernate"])
		}
	}

	return float64(active), float64(hibernated)
}

func countLicenseTypes(licenses rac.Query) (float64, float64) {
	var soft, hasp int
	for _, license := range licenses.Records {
		switch license["license-type"] {
		case "soft":
			soft++
		case "HASP":
			hasp++
		default:
			log.Printf("'rac session list --licenses': unexpected 'license-type' field value: '%s'", license["license-type"])
		}
	}

	return float64(soft), float64(hasp)
}

func countTotalProcMem(processes rac.Query) (float64, error) {
	var total int
	for _, process := range processes.Records {
		memory, err := strconv.Atoi(process["memory-size"])
		if err != nil {
			return 0, fmt.Errorf("parsing process 'memory-size': %w", err)
		}
		total += memory
	}

	return float64(total), nil
}

func createInfobaseNameMap(infobases rac.Query) map[string]string {
	baseUUIDToName := make(map[string]string)
	for _, record := range infobases.Records {
		baseUUIDToName[record["infobase"]] = record["name"]
	}

	return baseUUIDToName
}

func countSessionsByInfobase(sessions rac.Query) map[string]float64 {
	mapSessionUUIDToBase := make(map[string][]string)
	for _, record := range sessions.Records {
		if record["hibernate"] == "yes" {
			continue
		}
		mapSessionUUIDToBase[record["infobase"]] = append(mapSessionUUIDToBase[record["infobase"]], record["session"])
	}
	numSessionByBase := make(map[string]float64)
	for baseUUID, sessionUUID := range mapSessionUUIDToBase {
		numSessionByBase[baseUUID] = float64(len(sessionUUID))
	}
	return numSessionByBase
}

func recordMetrics(server *api.Server) {
	cluster := "--cluster=" + os.Getenv("platform1c_admin_cluster")

	// There are configurations without an administrator, but
	// this is insecure and metr1c only works with configurations
	// that have an administrator.
	adminusr := "--cluster-user=" + os.Getenv("platform1c_admin_user")
	// Examples: Администратор, Admin

	adminpass := "--cluster-pwd=" + os.Getenv("platform1c_admin_pw")
	// Examples: 1234, superpass, orsomethingsecure

	// Path to the executable file
	execPath := "/opt/1cv8/x86_64/" + os.Getenv("platform1c_version") + "/rac"
	// Examples: 8.3.24.1467

	// hidepid (Linux) must be equal 1 or it's unsecure.
	// rac accepts password and admin user as argument so any server user
	// may see it on htop if hidepid equals 0.

	baseQuery := rac.Query{
		ExecPath: execPath,
		Cluster:  cluster,
		User:     adminusr,
		Password: adminpass,
	}

	go func() {
		for {
			// Sessions
			sessions := getRecords(baseQuery, "session", "list", "")

			sessionCount.Set(float64(sessions.CountRecords()))

			active, hibernated := countSessionTypes(sessions)
			activeSessionCount.Set(active)
			hibernatedSessionCount.Set(hibernated)

			// Session licenses
			sessionsLicenses := getRecords(baseQuery, "session", "list", `--licenses`)
			soft, hasp := countLicenseTypes(sessionsLicenses)
			softLicensesCount.Set(soft)
			haspLicensesCount.Set(hasp)

			// Connections
			connections := getRecords(baseQuery, "connection", "list", "")
			connectionCount.Set(float64(connections.CountRecords()))

			// Processes
			processes := getRecords(baseQuery, "process", "list", "")
			processCount.Set(float64(processes.CountRecords()))

			memory, err := countTotalProcMem(processes)
			if err != nil {
				log.Println(err)
			}
			processMemTotal.Set(memory)

			// Infobases
			infobases := getRecords(baseQuery, "infobase", "summary", "list")
			numLicensesByBase := countSessionsByInfobase(sessions)
			for baseUUID, baseName := range createInfobaseNameMap(infobases) {
				if numSessions, ok := numLicensesByBase[baseUUID]; ok {
					sessionsPerInfobase.WithLabelValues(baseUUID, baseName).Set(numSessions)
				} else {
					sessionsPerInfobase.WithLabelValues(baseUUID, baseName).Set(0)
				}
			}

			server.UpdateSummary(api.Summary{
				SessionCount:       sessions.CountRecords(),
				SessionsActive:     int(active),
				SessionsHybernated: int(hibernated),
				UsedLicensesSoft:   int(soft),
				UsedLicensesHASP:   int(hasp),
				ConnectionCount:    connections.CountRecords(),
				ProcessCount:       processes.CountRecords(),
				ProcessesMemoryKB:  int(memory),
			})

			server.UpdateSessions(sessions.Records)
			server.UpdateConnections(connections.Records)
			server.UpdateProcesses(processes.Records)
			server.UpdateInfobases(infobases.Records)

			// Set a timeout before the next metrics gathering
			time.Sleep(60 * time.Second)
		}
	}()
}

var (
	sessionCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_sessions_count",
		Help: "The total number of 1c user sessions",
	})

	activeSessionCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_active_sessions_count",
		Help: "The total number of 1c user hybernated sessions",
	})

	hibernatedSessionCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_hybernated_sessions_count",
		Help: "The total number of 1c user hybernated sessions",
	})

	softLicensesCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_soft_licenses_count",
		Help: "The total number of 1c user used soft licenses",
	})

	haspLicensesCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_hasp_licenses_count",
		Help: "The total number of 1c user used hasp licenses",
	})

	connectionCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_connection_count",
		Help: "The total number of connections",
	})

	processCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_process_count",
		Help: "The total number of processes",
	})

	processMemTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_processes_total_memory_kbytes",
		Help: "The total number of used memory by all processes (KB)",
	})

	sessionsPerInfobase = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "platform1c_session_count_per_infobase",
		Help: "The number of sessions assigned with infobase"},
		// The two label names by which to split the metric.
		[]string{"InfobaseUUID", "InfobaseName"})
)

func main() {
	var help bool
	var version bool

	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&version, "version", false, "display version")
	flag.Parse()

	if help {
		fmt.Printf("metr1c - prometheus exporter for platform 1C\n")
		os.Exit(0)
	}

	if version {
		fmt.Printf("%s", Version)
		os.Exit(0)
	}

	apiServer := api.NewServer()
	recordMetrics(apiServer)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/api/summary", http.HandlerFunc(server.ServeSummary))
	http.Handle("/api/sessions", http.HandlerFunc(server.ServeSessions))
	http.Handle("/api/connections", http.HandlerFunc(server.ServeConnections))
	http.Handle("/api/processes", http.HandlerFunc(server.ServeProcesses))
	http.Handle("/api/infobases", http.HandlerFunc(server.ServeInfobases))

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("metr1c_port")),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(httpServer.ListenAndServe())
}
