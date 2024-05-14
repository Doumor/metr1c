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
			// ! rac session list
			// Examine current 1C session info
			sessions := rac.RACQuery{
				ExecPath:   execPath,
				Command:    "session",
				SubCommand: "list",
				Option:     "",
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

			var activeSessions, hibernatedSessions int = 0, 0

			for _, session := range sessions.Records {
				switch session["hibernate"] {
					case "no": activeSessions++
					case "yes" : hibernatedSessions++
					default : log.Fatal("'rac session list' hibernate field error")
				}
			}

			activeSessionCount.Set(float64(activeSessions))
			hibernatedSessionCount.Set(float64(hibernatedSessions))

			// ! rac session list --licenses
			// Examine the current 1C session information in terms of licenses
			sessionsLicenses := rac.RACQuery{
				ExecPath:   execPath,
				Command:    "session",
				SubCommand: "list",
				Option:     `--licenses`,
				Cluster:    cluster,
				User:       adminusr,
				Password:   adminpass,
			}

			err = sessionsLicenses.Run()

			if err != nil {
				log.Fatal(err)
			}

			err = sessionsLicenses.Parse()

			if err != nil {
				log.Fatal(err)
			}

			var softLicenses, haspLicenses int = 0, 0

			// Count field "license-type" for soft and hasp licenses
			for _, sessionLicense := range sessionsLicenses.Records {
				switch sessionLicense["license-type"] {
					case "soft": softLicenses++
					case "HASP": haspLicenses++
				}
			}

			softLicensesCount.Set(float64(softLicenses))
			haspLicensesCount.Set(float64(haspLicenses))

			// Set a timeout before the next metrics gathering
			time.Sleep(60 * time.Second) // 1 min
		}
	}()
}

var (
	// session list
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

	// session list --licenses
	softLicensesCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_soft_licenses_count",
		Help: "The total number of 1c user used soft licenses",
	})

	haspLicensesCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform1c_hasp_licenses_count",
		Help: "The total number of 1c user used hasp licenses",
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
