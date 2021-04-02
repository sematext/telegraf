package processors

import (
	"github.com/influxdata/telegraf"
)

var (
	measurementReplaces = map[string]string{
		"phpfpm":  "php",
		"mongodb": "mongo",
		"apache":  "apache",
		"nginx":   "nginx",
	}
	fieldReplaces = map[string]string{
		// apache
		"apache.BusyWorkers":          "apache.workers.busy",
		"apache.BytesPerReq":          "apache.bytes",
		"apache.ReqPerSec":            "apache.requests",
		"apache.ConnsAsyncClosing":    "apache.connections.async.closing",
		"apache.ConnsAsyncKeepAlive":  "apache.connections.async.keepAlive",
		"apache.ConnsAsyncWriting":    "apache.connections.async.writing",
		"apache.ConnsTotal":           "apache.connections",
		"apache.IdleWorkers":          "apache.workers.idle",
		"apache.scboard_closing":      "apache.workers.closing",
		"apache.scboard_dnslookup":    "apache.workers.dns",
		"apache.scboard_finishing":    "apache.workers.finishing",
		"apache.scboard_idle_cleanup": "apache.workers.cleanup",
		"apache.scboard_keepalive":    "apache.workers.keepalive",
		"apache.scboard_logging":      "apache.workers.logging",
		"apache.scboard_open":         "apache.workers.open",
		"apache.scboard_reading":      "apache.workers.reading",
		"apache.scboard_sending":      "apache.workers.sending",
		"apache.scboard_starting":     "apache.workers.starting",
		"apache.scboard_waiting":      "apache.workers.waiting",
		"php.accepted_conn":           "php.fpm.requests.accepted.conns",
		"php.listen_queue":            "php.fpm.queue.listen",
		"php.max_listen_queue":        "php.fpm.queue.listen.max",
		"php.listen_queue_len":        "php.fpm.queue.listen.len",
		"php.idle_processes":          "php.fpm.process.idle",
		"php.active_processes":        "php.fpm.process.active",
		"php.total_processes":         "php.fpm.process.total",
		"php.max_active_processes":    "php.fpm.process.active.max",
		"php.max_children_reached":    "php.fpm.process.childrenReached.max",
		"php.slow_requests":           "php.fpm.requests.slow",
		// nginx
		"nginx.accepts":  "nginx.requests.connections.accepted",
		"nginx.handled":  "nginx.requests.connections.handled",
		"nginx.active":   "nginx.requests.connections.active",
		"nginx.reading":  "nginx.requests.connections.reading",
		"nginx.writing":  "nginx.requests.connections.writing",
		"nginx.waiting":  "nginx.requests.connections.waiting",
		"nginx.requests": "nginx.request.count",
		// mongodb
		"mongo.flushes":                   "mongo.flushes",
		"mongo.flushes_total_time_ns":     "mongo.flushes.time",
		"mongo.document_inserted":         "mongo.documents.inserted",
		"mongo.document_updated":          "mongo.documents.updated",
		"mongo.document_deleted":          "mongo.documents.deleted",
		"mongo.document_returned":         "mongo.documents.returned",
		"mongo.resident_megabytes":        "mongo.memory.resident",
		"mongo.vsize_megabytes":           "mongo.memory.virtual",
		"mongo.mapped_megabytes":          "mongo.memory.mapped",
		"mongo.inserts":                   "mongo.ops.insert",
		"mongo.queries":                   "mongo.ops.query",
		"mongo.updates":                   "mongo.ops.update",
		"mongo.getmores":                  "mongo.ops.getmore",
		"mongo.commands":                  "mongo.ops.command",
		"mongo.repl_inserts":              "mongo.replica.ops.insert",
		"mongo.repl_queries":              "mongo.replica.ops.query",
		"mongo.repl_updates":              "mongo.replica.ops.update",
		"mongo.repl_deletes":              "mongo.replica.ops.delete",
		"mongo.repl_getmores":             "mongo.replica.ops.getmore",
		"mongo.repl_commands":             "mongo.replica.ops.command",
		"mongo.count_command_failed":      "mongo.commands.failed",
		"mongo.count_command_total":       "mongo.commands.total",
		"mongo.data_size":                 "mongo.database.data.size",
		"mongo.storage_size":              "mongo.database.storage.size",
		"mongo.index_size":                "mongo.database.index.size",
		"mongo.collections":               "mongo.database.collections",
		"mongo.objects":                   "mongo.database.objects",
		"mongo.connections_current":       "mongo.network.connections",
		"mongo.connections_total_created": "mongo.network.connections.total",
		"net_in_bytes":                    "mongo.network.transfer.rx.rate",
		"net_out_bytes":                   "mongo.network.transfer.tx.rate",
	}
)

// Rename processor renames the measurement (metric) names
// to match the existing metric names sent by Node.js agents
type Rename struct{}

// NewRename builds a new rename processor.
func NewRename() BatchProcessor { return &Rename{} }

// Process performs a lookup in the local maps of metric/field names
// and replaces the metric name with the new name.
func (r *Rename) Process(points []telegraf.Metric) ([]telegraf.Metric, error) {
	for _, point := range points {
		replace, ok := measurementReplaces[point.Name()]
		if !ok {
			continue
		}
		point.SetName(replace)
		for _, field := range point.FieldList() {
			key := point.Name() + "." + field.Key
			replace, ok := fieldReplaces[key]
			if !ok {
				continue
			}
			point.RemoveField(field.Key)
			point.AddField(replace, field.Value)
		}
	}
	return points, nil
}

func (Rename) Close() {}
