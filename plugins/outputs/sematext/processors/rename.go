package processors

import (
	"github.com/influxdata/telegraf"
)

var (
	measurementReplaces = map[string]string{
		"phpfpm":           "php",
		"mongodb":          "mongo",
		"mongodb_db_stats": "mongo",
		"apache":           "apache",
		"nginx":            "nginx",
	}
	fieldReplaces = map[string]string{
		// apache
		"apache.BusyWorkers":          "workers.busy",
		"apache.BytesPerReq":          "bytes",
		"apache.ReqPerSec":            "requests",
		"apache.ConnsAsyncClosing":    "connections.async.closing",
		"apache.ConnsAsyncKeepAlive":  "connections.async.keepAlive",
		"apache.ConnsAsyncWriting":    "connections.async.writing",
		"apache.ConnsTotal":           "connections",
		"apache.IdleWorkers":          "workers.idle",
		"apache.scboard_closing":      "workers.closing",
		"apache.scboard_dnslookup":    "workers.dns",
		"apache.scboard_finishing":    "workers.finishing",
		"apache.scboard_idle_cleanup": "workers.cleanup",
		"apache.scboard_keepalive":    "workers.keepalive",
		"apache.scboard_logging":      "workers.logging",
		"apache.scboard_open":         "workers.open",
		"apache.scboard_reading":      "workers.reading",
		"apache.scboard_sending":      "workers.sending",
		"apache.scboard_starting":     "workers.starting",
		"apache.scboard_waiting":      "workers.waiting",
		"php.accepted_conn":           "fpm.requests.accepted.conns",
		"php.listen_queue":            "fpm.queue.listen",
		"php.max_listen_queue":        "fpm.queue.listen.max",
		"php.listen_queue_len":        "fpm.queue.listen.len",
		"php.idle_processes":          "fpm.process.idle",
		"php.active_processes":        "fpm.process.active",
		"php.total_processes":         "fpm.process.total",
		"php.max_active_processes":    "fpm.process.active.max",
		"php.max_children_reached":    "fpm.process.childrenReached.max",
		"php.slow_requests":           "fpm.requests.slow",
		// nginx
		"nginx.accepts":  "requests.connections.accepted",
		"nginx.handled":  "requests.connections.handled",
		"nginx.active":   "requests.connections.active",
		"nginx.reading":  "requests.connections.reading",
		"nginx.writing":  "requests.connections.writing",
		"nginx.waiting":  "requests.connections.waiting",
		"nginx.requests": "request.count",
		// mongodb
		"mongo.flushes":                   "flushes",
		"mongo.flushes_total_time_ns":     "flushes.time",
		"mongo.document_inserted":         "documents.inserted",
		"mongo.document_updated":          "documents.updated",
		"mongo.document_deleted":          "documents.deleted",
		"mongo.document_returned":         "documents.returned",
		"mongo.resident_megabytes":        "memory.resident",
		"mongo.vsize_megabytes":           "memory.virtual",
		"mongo.mapped_megabytes":          "memory.mapped",
		"mongo.inserts":                   "ops.insert",
		"mongo.queries":                   "ops.query",
		"mongo.updates":                   "ops.update",
		"mongo.getmores":                  "ops.getmore",
		"mongo.commands":                  "ops.command",
		"mongo.repl_inserts":              "replica.ops.insert",
		"mongo.repl_queries":              "replica.ops.query",
		"mongo.repl_updates":              "replica.ops.update",
		"mongo.repl_deletes":              "replica.ops.delete",
		"mongo.repl_getmores":             "replica.ops.getmore",
		"mongo.repl_commands":             "replica.ops.command",
		"mongo.count_command_failed":      "commands.failed",
		"mongo.count_command_total":       "commands.total",
		"mongo.data_size":                 "database.data.size",
		"mongo.storage_size":              "database.storage.size",
		"mongo.index_size":                "database.index.size",
		"mongo.collections":               "database.collections",
		"mongo.objects":                   "database.objects",
		"mongo.connections_current":       "network.connections",
		"mongo.connections_total_created": "network.connections.total",
		"mongo.net_in_bytes":              "network.transfer.rx.rate",
		"mongo.net_out_bytes":             "network.transfer.tx.rate",
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
		removedFields := make([]string, 0)
		for _, field := range point.FieldList() {
			key := point.Name() + "." + field.Key
			replace, ok := fieldReplaces[key]
			if !ok {
				continue
			}
			// we can't remove the fields
			// while iterating because it
			// produces unwanted effects
			// e.g. metrics that have the
			// mapping are not renamed.
			// That's why we have to remove
			// them in a separate loop
			removedFields = append(removedFields, field.Key)
			point.AddField(replace, field.Value)
		}
		for _, f := range removedFields {
			point.RemoveField(f)
		}
	}
	return points, nil
}

func (Rename) Close() {}
