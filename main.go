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
	"time"

	"doumor/metr1c/rac"

	// prometheus exporter
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	// see in "rac" help
	cluster := "--cluster=" + os.Getenv("platform1c_admin_cluster")

	// 07593cfe-64c2-4656-be5f-61c3226286d5

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

	go func() {
		for {
			// Examine current 1C session info
			sessions := rac.RACQuery{
				ExecPath:   execPath,
				Command:    "session",
				SubCommand: "list",
				Cluster:    cluster,
				User:       adminusr,
				Password:   adminpass,
			}
			// Get output from a rac session list query
			err := sessions.Run()
			if err != nil {
				log.Fatal(err)
			}
			err = sessions.Parse()
			if err != nil {
				log.Fatal(err)
			}

			// Count current 1C sessions
			sessionCount.Set(float64(sessions.CountRecords()))

			connections := rac.RACQuery{
				ExecPath:   execPath,
				Command:    "connection",
				SubCommand: "list",
				Cluster:    cluster,
				User:       adminusr,
				Password:   adminpass,
			}

			err = connections.Run()
			if err != nil {
				log.Fatal(err)
			}
			err = connections.Parse()
			if err != nil {
				log.Fatal(err)
			}

			connectionCount.Set(float64(sessions.CountRecords()))

			// Set a timeout before the next metrics gathering
			time.Sleep(60 * time.Second) // 1 min
		}
	}()
}

var (
	sessionCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_sessions_count",
		Help: "The total number of 1c user licenses",
	})

	connectionCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_connection_count",
		Help: "The total number of connection",
	})
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
		fmt.Printf("v0.1.0\n")
		os.Exit(0)
	}

	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())

	port := ":" + os.Getenv("metr1c_port") // Example: 1599

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
