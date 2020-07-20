package metrics

import "regexp"

func init() {
	Register(collectReplicationMetrics)
}

var collectReplicationMetrics = map[string]MetricConfig{
	"master_slave_info": {
		Parser: &keyMatchParser{
			matches: map[string]string{
				"role": "master",
			},
			Parser: Parsers{
				&normalParser{
					MetricMeta: &MetaData{
						Name:      "connected_slaves",
						Help:      "the count of connected slaves, when pika serve instance's role is master",
						Type:      metricTypeGauge,
						Labels:    []string{LabelNameAddr, LabelNameAlias},
						ValueName: "connected_slaves",
					},
				},
				&regexParser{
					reg: regexp.MustCompile(`slave\d+:ip=(?P<slave_ip>[\d.]+),port=(?P<slave_port>[\d.]+),` +
						`state=(?P<slave_state>[a-z]+),sid=(?P<slave_sid>[\d]+),lag=(?P<slave_lag>[\d]+)`),
					MetricMeta: &MetaDatas{
						{
							Name:      "slave_state",
							Help:      "pika serve instance slave's state",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "slave_sid", "slave_ip", "slave_port"},
							ValueName: "slave_state",
						},
						{
							Name:      "slave_lag",
							Help:      "pika serve instance slave's binlog lag",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "slave_sid", "slave_ip", "slave_port"},
							ValueName: "slave_lag",
						},
					},
				},
			},
		},
	},

	"slave_info": {
		Parser: &keyMatchParser{
			matches: map[string]string{
				"role": "slave",
			},
			Parser: Parsers{
				&normalParser{
					MetricMeta: MetaDatas{
						{
							Name: "master_link_status",
							Help: "connection state between slave and master, when pika serve instance's " +
								"role is slave",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "master_host", "master_port"},
							ValueName: "master_link_status",
						},
						{
							Name: "repl_state",
							Help: "sync connection state between slave and master, when pika serve instance's " +
								"role is slave",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "master_host", "master_port"},
							ValueName: "repl_state",
						},
						{
							Name:      "slave_read_only",
							Help:      "is slave read only, when pika serve instance's role is slave",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "master_host", "master_port"},
							ValueName: "slave_read_only",
						},
					},
				},
				&versionMatchParser{
					verC: mustNewVersionConstraint(`>=3.0.0`),
					Parser: &normalParser{
						MetricMeta: &MetaData{
							Name:      "slave_priority",
							Help:      "slave priority, when pika serve instance's role is slave",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "master_host", "master_port"},
							ValueName: "slave_priority",
						},
					},
				},
			},
		},
	},

	"double_master_info": {
		Parser: &keyMatchParser{
			matches: map[string]string{
				"role":               "master",
				"double_master_mode": "true",
			},
			Parser: &regexParser{
				reg: regexp.MustCompile(`the peer-master host:(?P<the_peer_master_host>[^\n]*)[\s\S]*` +
					`the peer-master port:(?P<the_peer_master_port>[^\n]*)[\s\S]*` +
					`the peer-master server_id:(?P<the_peer_master_server_id>[^\n]*)[\s\S]*` +
					`repl_state:(?P<double_master_repl_state>[^\n]*)[\s\S]*` +
					`double_master_recv_info:\s*filenum\s*(?P<double_master_recv_info_binlog_filenum>[^\s]*)` +
					`\s*offset\s*(?P<double_master_recv_info_binlog_offset>[^\n]*)`),
				MetricMeta: MetaDatas{
					{
						Name: "double_master_info",
						Help: "the peer master info, when pika serve instance's role is master and " +
							"double_master_mode is true",
						Type: metricTypeGauge,
						Labels: []string{LabelNameAddr, LabelNameAlias, "the_peer_master_server_id",
							"the_peer_master_host", "the_peer_master_port"},
					},
					{
						Name: "double_master_repl_state",
						Help: "double master sync state, when pika serve instance's role is master and " +
							"double_master_mode is true",
						Type: metricTypeGauge,
						Labels: []string{LabelNameAddr, LabelNameAlias, "the_peer_master_server_id",
							"the_peer_master_host", "the_peer_master_port"},
						ValueName: "double_master_repl_state",
					},
					{
						Name: "double_master_recv_info_binlog_filenum",
						Help: "double master recv binlog file num, when pika serve instance's role is master and " +
							"double_master_mode is true",
						Type: metricTypeGauge,
						Labels: []string{LabelNameAddr, LabelNameAlias, "the_peer_master_server_id",
							"the_peer_master_host", "the_peer_master_port"},
						ValueName: "double_master_recv_info_binlog_filenum",
					},
					{
						Name: "double_master_recv_info_binlog_offset",
						Help: "double master recv binlog offset, when pika serve instance's role is master and " +
							"double_master_mode is true",
						Type: metricTypeGauge,
						Labels: []string{LabelNameAddr, LabelNameAlias, "the_peer_master_server_id",
							"the_peer_master_host", "the_peer_master_port"},
						ValueName: "double_master_recv_info_binlog_offset",
					},
				},
			},
		},
	},
}
