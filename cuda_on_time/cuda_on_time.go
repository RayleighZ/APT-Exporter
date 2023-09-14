package cuda_on_time

import (
	"APT-Exporter/utils"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"os"
	"strings"
	"time"
)

type User struct {
	Uid      string `json:"uid"`
	UserName string `json:"user_name"`
}

type CudaGauge struct {
	gauge prometheus.Gauge
	user  User
}

var cudaGaugeMap = make(map[string]CudaGauge)

// init user cuda gauge from ./user.json
func initCudaUser() {
	content, err := os.ReadFile("./cuda_on_time/user.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	jsonString := string(content)
	var userList []User
	err = json.Unmarshal([]byte(jsonString), &userList)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	log.Println("userlist")
	log.Println(userList)

	for _, user := range userList {
		gaugeName := "user_cuda_on_time_" + user.UserName
		gauge := promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: gaugeName,
				Help: "User's CUDA use time gauge, unit is minutes",
			})

		//prometheus.MustRegister(gauge)
		cudaGaugeMap[user.Uid] = CudaGauge{user: user, gauge: gauge}
	}
}

// add user cuda on time, unit is minutes
func addCountByUid(uid string, useTime float64) {
	cudaGauge, ok := cudaGaugeMap[uid]
	log.Println(ok)
	if ok {
		cudaGauge.gauge.Add(useTime)
	}
}

func getUidViaPid(pid string) string {
	return "1007"
	cmd := "ps -p " + pid + " -o uid"
	res := utils.ExecAndGetRes(cmd)
	lines := strings.Split(res, "\n")
	if len(lines) == 1 {
		return "none"
	}
	res = strings.Split(res, "\n")[1]
	uid := strings.Fields(res)[0]
	return uid
}

func collectCudaOnTime() {
	workingList := make(map[string]bool)
	res := utils.ExecAndGetRes("nvidia-smi")
	smiInfo := strings.Split(res, "\n")
	//log.Println(smiInfo)
	cudaUserInfoSplitLine := "ID   ID                                                             Usage"
	cudaInfoIndex := -1
	for i, s := range smiInfo {
		if strings.Contains(s, cudaUserInfoSplitLine) {
			cudaInfoIndex = i + 1
		}
	}
	cudaInfoLines := smiInfo[cudaInfoIndex+1 : len(smiInfo)-2]
	for _, cudaInfo := range cudaInfoLines {
		fields := strings.Fields(cudaInfo)
		fields = fields[1 : len(fields)-1]
		pid := fields[3]
		uid := getUidViaPid(pid)
		log.Println(uid)
		if uid != "none" {
			workingList[uid] = true
		}
	}
	for uid := range workingList {
		addCountByUid(uid, 15)
	}
}

func resetCudaGauge() {
	for _, gauge := range cudaGaugeMap {
		gauge.gauge.Set(0)
	}
}

func Start() {
	collectTicker := time.NewTicker(15 * time.Minute)
	resetTicker := time.NewTicker(7 * 24 * time.Hour)
	initCudaUser()
	go func() {
		for {
			select {
			case <-collectTicker.C:
				collectCudaOnTime()
			case <-resetTicker.C:
				resetCudaGauge()
			}
		}
	}()
}
