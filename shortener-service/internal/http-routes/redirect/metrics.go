package redirect

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

var redirectMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "redirect",
	Subsystem:  "http",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"code"})

func observeRedirect(d time.Duration, status int) {
	redirectMetrics.WithLabelValues(strconv.Itoa(status)).Observe(d.Seconds())
}
