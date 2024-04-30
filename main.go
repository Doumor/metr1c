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
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "regexp"
    "time"


    // prometheus exporter
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
    // see in "rac" help
    cluster := "--cluster="+ os.Getenv("platform1c_admin_cluster")
    // 07593cfe-64c2-4656-be5f-61c3226286d5

    // There are configurations without an administrator, but
    // this is insecure and metr1c only works with configurations
    // that have an administrator.
    adminusr := "--cluster-user=" + os.Getenv("platform1c_admin_user")
    // Examples: Администратор, Admin

    adminpass := "--cluster-pwd=" + os.Getenv("platform1c_admin_pw")
    // Examples: 1234, superpass, orsomethingsecure

    // Path to executable file
    progrun := "/opt/1cv8/x86_64/" + os.Getenv("platform1c_version") + "/rac"
    // Examples: 8.3.24.1467



    sessionListArgs := []string{"session", "list", cluster, adminusr, adminpass}
    // hidepid (Linux) must be equal 1 or it's unsecure.
    // rac accepts password and admin user as argument so any server user
    // may see it on htop if hidepid equals 0.

    // getting and parsing rac session list output
    go func() {
        for{
            // ! Output from rac session list
            out, err := exec.Command(progrun, sessionListArgs...).Output()
            if err != nil {
                log.Fatal(err)
            }

            // Session count
            re := regexp.MustCompile(`session-id *:.\d+\n`)
            sessionCount.Set(float64(len(re.FindAllString(string(out), -1))))


            // Timer
            time.Sleep(60 * time.Second) // 1 min
        }
    }()
}

var (
    sessionCount = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "platform1c_sessions_count",
        Help: "The total number of 1c user licenses",
    })
)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "-h" {
        fmt.Printf("v0.1.0\n")
        os.Exit(0)
    }

    recordMetrics()

    http.Handle("/metrics", promhttp.Handler())
    port := ":" + os.Getenv("metr1c_port") // Example: 1599
    http.ListenAndServe(port, nil)
    // We use port like other 1C products (i.e. 1545, 1540, 1541, 1560-1591)
}
