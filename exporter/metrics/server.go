package metrics

import "regexp"

func init() {
	Register(collectServerMetrics)
}

var collectServerMetrics = map[string]MetricConfig{
	"build_info": {
		Parser: &regexParser{
			reg: regexp.MustCompile(`pika_version:(?P<pika_version>[^\n]*)[\s\S]*` +
				`pika_git_sha:(?P<pika_git_sha>[^\n]*)[\s\S]*` +
				`pika_build_compile_date:(?P<pika_build_compile_date>[^\n]*)[\s\S]*` +
				`os:(?P<os>[^\n]*)[\s\S]*` +
				`arch_bits:(?P<arch_bits>[^\n]*)`),
			MetricMeta: &MetaData{
				Name:   "build_info",
				Help:   "pika binary file build info",
				Type:   metricTypeGauge,
				Labels: []string{LabelNameAddr, LabelNameAlias, "os", "arch_bits", "pika_version", "pika_git_sha", "pika_build_compile_date"},
			},
		},
	},
	"server_info": {
		Parser: &regexParser{
			reg: regexp.MustCompile(`process_id:(?P<process_id>[^\n]*)[\s\S]*` +
				`tcp_port:(?P<tcp_port>[^\n]*)[\s\S]*` +
				`config_file:(?P<config_file>[^\n]*)[\s\S]*` +
				`server_id:(?P<server_id>[^\n]*)[\s\S]*` +
				`role:(?P<role>[^\n]*)`),
			MetricMeta: &MetaData{
				Name:   "server_info",
				Help:   "pika serve instance info",
				Type:   metricTypeGauge,
				Labels: []string{LabelNameAddr, LabelNameAlias, "process_id", "tcp_port", "config_file", "server_id", "role"},
			},
		},
	},
	"uptime_in_seconds": {
		Parser: &normalParser{
			MetricMeta: &MetaData{
				Name:      "uptime_in_seconds",
				Help:      "pika serve instance uptime in seconds",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias},
				ValueName: "uptime_in_seconds",
			},
		},
	},
	"thread_num": {
		Parser: &normalParser{
			MetricMeta: &MetaData{
				Name:      "thread_num",
				Help:      "pika serve instance thread num",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias},
				ValueName: "thread_num",
			},
		},
	},
	"sync_thread_num": {
		Parser: &normalParser{
			MetricMeta: &MetaData{
				Name:      "sync_thread_num",
				Help:      "pika serve instance sync thread num",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias},
				ValueName: "sync_thread_num",
			},
		},
	},
}
