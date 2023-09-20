package zombie_detect

import (
	"APT-Exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type ZombieProcess struct {
	pid        string
	zombieTime prometheus.Gauge
}

var zombieMap = make(map[string][]ZombieProcess)

func initZombieGauge() {
	// init uid to zombie process slice map
	// will fill slice map with ZombieProcess when zp is detected
	userSlice := utils.GetUser()
	for _, user := range userSlice {
		zombieMap[user.Uid] = make([]ZombieProcess, 0)
	}
}

func detectZombieProcess() {

}

func Start() {
	collectTicker := time.NewTicker(15 * time.Minute)
	initZombieGauge()
	go func() {
		for {
			select {
			case <-collectTicker.C:
				detectZombieProcess()
			}
		}
	}()
}
