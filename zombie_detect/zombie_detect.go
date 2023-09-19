package zombie_detect

import (
	"APT-Exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type ZombieProcess struct {
	pid        string
	uid        string
	zombieTime prometheus.Gauge
}

func initZombieGauge() {
	// init zombie gauge counter from user.json
	userSlice = utils.GetUser()

}

var Uid2UserNameMap = make(map[string]string)
