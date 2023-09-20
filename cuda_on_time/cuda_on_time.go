package cuda_on_time

import (
	"APT-Exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"strings"
	"time"
)

type CudaGauge struct {
	gauge prometheus.Gauge
	user  utils.User
}

var cudaGaugeMap = make(map[string]CudaGauge)

// init user cuda gauge from ./user.json
func initCudaUser() {
	var userList = utils.GetUser()
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

func collectCudaOnTime() {
	workingList := make(map[string]bool)
	res := utils.ExecAndGetRes("nvidia-smi")
	smiInfo := strings.Split(res, "\n")
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
		uid := utils.GetUidViaPid(pid)
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
