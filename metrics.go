package main

import (
	"errors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lazyfrosch/dslmodem_exporter/pkg/zyxel"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	up = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "up",
		Help:      "If we are connected to the modem.",
	})
	collectionSeconds = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "collection_seconds",
		Help:      "Retrieval time for the DSL statistics from the modem.",
	})
	status = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "status",
		Name:      "info",
		Help:      "Metadata.",
	}, []string{"status", "mode", "profile", "traffic_type"})
	linkUptime = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "uptime_seconds",
		Help:      "Since when the link is established.",
	})
	linkRateUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "rate_up_bytes",
		Help:      "Rate of the upstream link in bits/s.",
	})
	linkRateDown = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "rate_down_bytes",
		Help:      "Rate of the downstream link in bits/s.",
	})
	snrMarginUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "snr_margin_up_db",
		Help:      "Signal to noise margin for upstream in Decibel.",
	})
	snrMarginDown = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "snr_margin_down_db",
		Help:      "Signal to noise margin for downstream in Decibel.",
	})
	transmitPowerUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "transmit_power_up_dbm",
		Help:      "Current transmitting power to upstream in Decibel milliwatt.",
	})
	transmitPowerDown = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "transmit_power_down_dbm",
		Help:      "Current transmitting power to downstream in Decibel milliwatt.",
	})
	receivePowerUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "receive_power_up_dbm",
		Help:      "Current receiving power on upstream in Decibel milliwatt.",
	})
	receivePowerDown = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "receive_power_down_dbm",
		Help:      "Current receiving power on downstream in Decibel milliwatt.",
	})
	attenuationUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "attenuation_up_dbm",
		Help:      "Total attenuation for upstream in Decibel.",
	})
	attenuationDown = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "attenuation_down_dbm",
		Help:      "Total attenuation for downstream in Decibel.",
	})
	attainableDataRateUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "attainable_rate_up_bytes",
		Help:      "Attainable data rate for upstream in bits/s.",
	})
	attainableDataRateDown = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "link",
		Name:      "attainable_rate_down_bytes",
		Help:      "Attainable data rate for downstream in bits/s.",
	})
)

func updateLoop(client *zyxel.Client, interval time.Duration, logger log.Logger) {
	var (
		loggedIn bool
		err      error
		stats    *zyxel.VDSLStatus
	)

	for {
		if !loggedIn {
			loggedIn = true
			_ = level.Info(logger).Log("msg", "Logging into DSL modem", "url", client.BaseURL.String())

			err = client.Login()
			if err != nil {
				_ = level.Error(logger).Log("msg", "Error during login", "err", err)
				time.Sleep(interval)
				continue
			}
		} else {
			time.Sleep(interval)
		}

		_ = level.Debug(logger).Log("msg", "Retrieving DSL statistics")

		begin := time.Now()
		stats, err = client.GetXDSLStatistics()
		if err != nil {
			if errors.Is(err, zyxel.ErrHTTPUnauthorized) {
				// Force re-login
				loggedIn = false
				continue
			}

			// TODO: clear other states?
			up.Set(0)

			_ = level.Error(logger).Log("msg", "Error during retrieving statistics", "err", err)
		} else {
			up.Set(1)
			collectionSeconds.Set(time.Since(begin).Seconds())

			status.With(prometheus.Labels{
				"status":       stats.Status,
				"mode":         stats.Mode,
				"profile":      stats.Profile,
				"traffic_type": stats.TrafficType,
			}).Set(1)

			linkUptime.Set(stats.LinkUptime.Seconds())
			linkRateUp.Set(stats.LineRateUp.BitsPerSec())
			linkRateDown.Set(stats.LineRateDown.BitsPerSec())

			snrMarginUp.Set(float64(stats.SNRMarginUp))
			snrMarginDown.Set(float64(stats.SNRMarginDown))

			transmitPowerUp.Set(float64(stats.TransmitPowerUp))
			transmitPowerDown.Set(float64(stats.TransmitPowerDown))

			receivePowerUp.Set(float64(stats.ReceivePowerUp))
			receivePowerDown.Set(float64(stats.ReceivePowerDown))

			attenuationUp.Set(float64(stats.AttenuationUp))
			attenuationDown.Set(float64(stats.AttenuationDown))

			attainableDataRateUp.Set(stats.AttainableDataRateUp.BitsPerSec())
			attainableDataRateDown.Set(stats.AttainableDataRateDown.BitsPerSec())
		}
	}
}
