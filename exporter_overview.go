package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	RegisterExporter("overview", newExporterOverview)
}

var overviewMetricDescription = map[string]prometheus.Gauge{
	"object_totals.channels":                    newGauge("channelsTotal", "Total number of open channels."),
	"object_totals.connections":                 newGauge("connectionsTotal", "Total number of open connections."),
	"object_totals.consumers":                   newGauge("consumersTotal", "Total number of message consumers."),
	"object_totals.queues":                      newGauge("queuesTotal", "Total number of queues in use."),
	"object_totals.exchanges":                   newGauge("exchangesTotal", "Total number of exchanges in use."),
	"queue_totals.messages":                     newGauge("queue_messages_total", "Total number ready and unacknowledged messages in cluster."),
	"queue_totals.messages_ready":               newGauge("queue_messages_ready_total", "Total number of messages ready to be delivered to clients."),
	"queue_totals.messages_unacknowledged":      newGauge("queue_messages_unacknowledged_total", "Total number of messages delivered to clients but not yet acknowledged."),
	"message_stats.ack_details.rate":            newGauge("message_stats_ack_details_rate", "Total number of messages ack rate per second."),
	"message_stats.confirm_details.rate":        newGauge("message_stats_confirm_details_rate", "Confirm rate."),
	"message_stats.deliver_details.rate":        newGauge("message_stats_deliver_details_rate", "Deliver rate."),
	"message_stats.deliver_get_details.rate":    newGauge("message_stats_deliver_get_details_rate", "Deliver get rate."),
	"message_stats.deliver_no_ack_details.rate": newGauge("message_stats_deliver_no_ack_details_rate", "Deliver no ack rate."),
	"message_stats.disk_reads_details.rate":     newGauge("message_stats_disk_reads_details_rate", "Disk reads rate."),
	"message_stats.disk_writes_details.rate":    newGauge("message_stats_disk_writes_details_rate", "Disk writes rate."),
	"message_stats.get_details.rate":            newGauge("message_stats_get_details_rate", "Get rate."),
	"message_stats.get_no_ack_details.rate":     newGauge("message_stats_get_no_ack_details_rate", "Get no ack rate."),
	"message_stats.publish_details.rate":        newGauge("message_stats_publish_details_rate", "Publish rate."),
	"message_stats.redeliver_details.rate":      newGauge("message_stats_redeliver_details_rate", "Redeliver rate."),
}

type exporterOverview struct {
	overviewMetrics map[string]prometheus.Gauge
}

func newExporterOverview() Exporter {
	return exporterOverview{
		overviewMetrics: overviewMetricDescription,
	}
}

func (e exporterOverview) String() string {
	return "Exporter overview"
}

func (e exporterOverview) Collect(ch chan<- prometheus.Metric) error {
	rabbitMqOverviewData, err := getMetricMap(config, "overview")

	if err != nil {
		return err
	}

	log.WithField("overviewData", rabbitMqOverviewData).Debug("Overview data")
	for key, gauge := range e.overviewMetrics {
		if value, ok := rabbitMqOverviewData[key]; ok {
			log.WithFields(log.Fields{"key": key, "value": value}).Debug("Set overview metric for key")
			gauge.Set(value)
		}
	}

	for _, gauge := range e.overviewMetrics {
		gauge.Collect(ch)
	}
	return nil
}

func (e exporterOverview) Describe(ch chan<- *prometheus.Desc) {
	for _, gauge := range e.overviewMetrics {
		gauge.Describe(ch)
	}

}
