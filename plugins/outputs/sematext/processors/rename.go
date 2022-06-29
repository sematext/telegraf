package processors

import (
	"github.com/influxdata/telegraf"
)

var (
	measurementReplaces = map[string]string{
		"phpfpm":                                      "php",
		"mongodb":                                     "mongo",
		"mongodb_db_stats":                            "mongo",
		"mongodb_col_stats":                           "mongo",
		"mongodb_shard_stats":                         "mongo",
		"apache":                                      "apache",
		"nginx":                                       "nginx",
		"win_cpu.Percent_DPC_Time":                    "os.cpu.dpc.time",
		"win_cpu.Percent_Idle_Time":                   "os.cpu.idle.time",
		"win_cpu.Percent_Interrupt_Time":              "os.cpu.interrupt.time",
		"win_cpu.Percent_Privileged_Time":             "os.cpu.privileged.time",
		"win_cpu.Percent_Processor_Time":              "os.cpu.processor.time",
		"win_cpu.Percent_User_Time":                   "os.cpu.user.time",
		"win_disk.Percent_Free_Space":                 "os.disk.free.bytes",
		"win_diskio.Disk_Read_Bytes_persec":           "os.disk.read.bytes",
		"win_diskio.Disk_Write_Bytes_persec":          "os.disk.write.bytes",
		"win_mem.Available_Bytes":                     "os.memory.free",
		"win_mem.Modified_Page_List_Bytes":            "os.memory.modified.page.list.bytes",
		"win_mem.Standby_Cache_Core_Bytes":            "os.memory.standby.cache.core.bytes",
		"win_mem.Standby_Cache_Normal_Priority_Bytes": "os.memory.standby.cache.normal.priority.bytes",
		"win_mem.Standby_Cache_Reserve_Bytes":         "os.memory.standby.cache.reserve.bytes",
		"win_net.Bytes_Received_persec":               "os.network.rx",
		"win_net.Bytes_Sent_persec":                   "os.network.tx",
		"win_swap.Percent_Usage":                      "os.swap.percentage.usage",
		"win_system.Processor_Queue_Length":           "os.processor.queue.length",
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
		"phpfpm.accepted_conn":        "fpm.requests.accepted.conns",
		"phpfpm.listen_queue":         "fpm.queue.listen",
		"phpfpm.max_listen_queue":     "fpm.queue.listen.max",
		"phpfpm.listen_queue_len":     "fpm.queue.listen.len",
		"phpfpm.idle_processes":       "fpm.process.idle",
		"phpfpm.active_processes":     "fpm.process.active",
		"phpfpm.total_processes":      "fpm.process.total",
		"phpfpm.max_active_processes": "fpm.process.active.max",
		"phpfpm.max_children_reached": "fpm.process.childrenReached.max",
		"phpfpm.slow_requests":        "fpm.requests.slow",
		// nginx
		"nginx.accepts":  "requests.connections.accepted",
		"nginx.handled":  "requests.connections.handled",
		"nginx.active":   "requests.connections.active",
		"nginx.reading":  "requests.connections.reading",
		"nginx.writing":  "requests.connections.writing",
		"nginx.waiting":  "requests.connections.waiting",
		"nginx.requests": "request.count",
		// mongodb
		"mongodb.flushes":                 "flushes",
		"mongodb.flushes_total_time_ns":   "flushes.time",
		"mongodb.document_inserted":       "documents.inserted",
		"mongodb.document_updated":        "documents.updated",
		"mongodb.document_deleted":        "documents.deleted",
		"mongodb.document_returned":       "documents.returned",
		"mongodb.resident_megabytes":      "memory.resident",
		"mongodb.vsize_megabytes":         "memory.virtual",
		"mongodb.mapped_megabytes":        "memory.mapped",
		"mongodb.inserts":                 "ops.insert",
		"mongodb.queries":                 "ops.query",
		"mongodb.updates":                 "ops.update",
		"mongodb.getmores":                "ops.getmore",
		"mongodb.commands":                "ops.command",
		"mongodb.repl_inserts":            "replica.ops.insert",
		"mongodb.repl_queries":            "replica.ops.query",
		"mongodb.repl_updates":            "replica.ops.update",
		"mongodb.repl_deletes":            "replica.ops.delete",
		"mongodb.repl_getmores":           "replica.ops.getmore",
		"mongodb.repl_commands":           "replica.ops.command",
		"mongodb.count_command_failed":    "commands.failed",
		"mongodb.count_command_total":     "commands.total",
		"mongodb_db_stats.data_size":      "database.data.size",
		"mongodb_db_stats.storage_size":   "database.storage.size",
		"mongodb_db_stats.index_size":     "database.index.size",
		"mongodb_db_stats.collections":    "database.collections",
		"mongodb_db_stats.objects":        "database.objects",
		"mongodb_db_stats.avg_obj_size":   "database.avg_obj_size",
		"mongodb_db_stats.indexes":        "database.indexes",
		"mongodb_db_stats.num_extents":    "database.num_extents",
		"mongodb_db_stats.ok":             "database.ok",
		"mongo.connections_current":       "network.connections",
		"mongo.connections_total_created": "network.connections.total",
		"mongo.net_in_bytes":              "network.transfer.rx.rate",
		"mongo.net_out_bytes":             "network.transfer.tx.rate",
		// mongodb_col_stats -> these appear like they map to the same thing, but "from" side is actually
		// "name.metricName" and "to" side is just the new "metricName"
		"mongodb_col_stats.avg_obj_size":     "mongodb_col_stats.avg_obj_size",
		"mongodb_col_stats.count":            "mongodb_col_stats.count",
		"mongodb_col_stats.ok":               "mongodb_col_stats.ok",
		"mongodb_col_stats.size":             "mongodb_col_stats.size",
		"mongodb_col_stats.storage_size":     "mongodb_col_stats.storage_size",
		"mongodb_col_stats.total_index_size": "mongodb_col_stats.total_index_size",
		// mongodb_shard_stats -> same logic as for mongodb_col_stats
		"mongodb_shard_stats.in_use":     "mongodb_shard_stats.in_use",
		"mongodb_shard_stats.available":  "mongodb_shard_stats.available",
		"mongodb_shard_stats.created":    "mongodb_shard_stats.created",
		"mongodb_shard_stats.refreshing": "mongodb_shard_stats.refreshing",

		"win_cpu.Percent_DPC_Time":                    "os.cpu.dpc.time",
		"win_cpu.Percent_Idle_Time":                   "os.cpu.idle.time",
		"win_cpu.Percent_Interrupt_Time":              "os.cpu.interrupt.time",
		"win_cpu.Percent_Privileged_Time":             "os.cpu.privileged.time",
		"win_cpu.Percent_Processor_Time":              "os.cpu.processor.time",
		"win_cpu.Percent_User_Time":                   "os.cpu.user.time",
		"win_disk.Percent_Free_Space":                 "os.disk.free.bytes",
		"win_diskio.Disk_Read_Bytes_persec":           "os.disk.read.bytes",
		"win_diskio.Disk_Write_Bytes_persec":          "os.disk.write.bytes",
		"win_mem.Available_Bytes":                     "os.memory.free",
		"win_mem.Modified_Page_List_Bytes":            "os.memory.modified.page.list.bytes",
		"win_mem.Standby_Cache_Core_Bytes":            "os.memory.standby.cache.core.bytes",
		"win_mem.Standby_Cache_Normal_Priority_Bytes": "os.memory.standby.cache.normal.priority.bytes",
		"win_mem.Standby_Cache_Reserve_Bytes":         "os.memory.standby.cache.reserve.bytes",
		"win_net.Bytes_Received_persec":               "os.network.rx",
		"win_net.Bytes_Sent_persec":                   "os.network.tx",
		"win_swap.Percent_Usage":                      "os.swap.percentage.usage",
		"win_system.Processor_Queue_Length":           "os.processor.queue.length",
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
		originalName := point.Name()
		replace, ok := measurementReplaces[originalName]
		if !ok {
			continue
		}
		point.SetName(replace)
		removedFields := make([]string, 0)
		for _, field := range point.FieldList() {
			key := originalName + "." + field.Key
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
