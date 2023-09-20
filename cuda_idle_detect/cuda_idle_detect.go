package cuda_idle_detect

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IdleCard struct {
	uid      string
	pid      string
	idleTime prometheus.Gauge
}

var idleCardSlice = make([]IdleCard, 6)

func detect() {

}
