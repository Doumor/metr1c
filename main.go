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
    "os"
    "os/exec"
    //"strings"
    "log"
    "time"
    "regexp"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
    //См. пояснения в rac
    cluster := "--cluster="+ os.Getenv("platform1c_admin_cluster")
    // 07593cfe-64c2-4656-be5f-61c3226286d5

    // Если админа нет, то создайте, это не очень безопасно.
    adminusr := "--cluster-user=" + os.Getenv("platform1c_admin_user")
    // Администратор

    adminpass := "--cluster-pwd=" + os.Getenv("platform1c_admin_pw")
    // 1234

    args := []string{"session", "list", cluster, adminusr, adminpass}
    // hidepid=1 иначе другие пользователи на сервере могут через htop увидеть
    // пользователя и пароль кластера 1С.

    progrun := "/opt/1cv8/x86_64/" + os.Getenv("platform1c_version") + "/rac"
    // 8.3.24.1467
    // Путь до исполняемого файла

    go func() {
        for{
            // ! Вывод из rac session list
            out, err := exec.Command(progrun, args...).Output()
            if err != nil {
                log.Fatal(err)
            }

            // Количество сессий
            re := regexp.MustCompile(`session-id *:.\d+\n`)
            sessionCount.Set(float64(len(re.FindAllString(string(out), -1))))


            // Таймер
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
    port := ":" + os.Getenv("metr1c_port") // Например 1599
    http.ListenAndServe(port, nil)
    // Использую порт как 1c (i.e. 1545, 1540, 1541, 1560-1591)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
