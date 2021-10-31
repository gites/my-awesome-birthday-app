package service

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var Version = "latest"

var (
	bDayReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bday_app_processed_req_total",
		Help: "The total number of processed requsts",
	})
)

// runMetricsServer runs dedicated metric and health http server
func (b *BDayService) runMetricsServer() {

	log.WithField("context", "runMetricsServer").Infof("starting metrics server on port %s", b.MetricsPort)

	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/live", b.liveHandle)
	metricsMux.HandleFunc("/ready", b.readyHandle)
	metricsMux.HandleFunc("/version", b.versionHandle)
	metricsMux.Handle("/metrics", promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    ":" + b.MetricsPort,
		Handler: metricsMux,
	}

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("context", "runMetricsServer").WithError(err).Fatalf("metric server listen and serve")
		}
	}()
	b.waitForExitSignal(metricsServer, "runMetricsServer")

}

func (b *BDayService) versionHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", Version)
}

func (b *BDayService) liveHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (b *BDayService) readyHandle(w http.ResponseWriter, req *http.Request) {
	if b.isReady {
		fmt.Fprintf(w, "ok")
	} else {
		w.WriteHeader(http.StatusTooEarly)
		fmt.Fprintf(w, "not ready")
	}
}
