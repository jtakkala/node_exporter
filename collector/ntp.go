// +build !nontp

package collector

import (
	"flag"
	"fmt"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

var (
	ntpServer = flag.String("collector.ntp.server", "", "NTP server to use for ntp collector.")
)

type ntpCollector struct {
	drift prometheus.Gauge
}

func init() {
	Factories["ntp"] = NewNtpCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// the offset between ntp and the current system time.
func NewNtpCollector() (Collector, error) {
	if *ntpServer == "" {
		return nil, fmt.Errorf("No NTP server specifies, see --ntpServer")
	}

	return &ntpCollector{
		drift: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "ntp_drift_seconds",
			Help:      "Time between system time and ntp time.",
		}),
	}, nil
}

func (c *ntpCollector) Update(ch chan<- prometheus.Metric) (err error) {
	t, err := ntp.Time(*ntpServer)
	if err != nil {
		return fmt.Errorf("Couldn't get ntp drift: %s", err)
	}
	drift := t.Sub(time.Now())
	log.Debugf("Set ntp_drift_seconds: %f", drift.Seconds())
	c.drift.Set(drift.Seconds())
	c.drift.Collect(ch)
	return err
}
