package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HeartbeatsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rcp_heartbeats_total",
		Help: "Total number of heartbeats received",
	})

	ActiveDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "rcp_active_devices",
		Help: "Number of currently online devices",
	})

	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rcp_http_requests_total",
		Help: "Total HTTP requests by method and path",
	}, []string{"method", "path", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "rcp_http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	RabbitMQMessagesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rcp_rabbitmq_messages_total",
		Help: "Total messages published to RabbitMQ",
	})

	WorkerJobsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rcp_worker_jobs_total",
		Help: "Total jobs processed by worker pool",
	})
)