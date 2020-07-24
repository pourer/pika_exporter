package metrics

import "regexp"

func init() {
	Register(collectKeySpaceMetrics)
}

var collectKeySpaceMetrics = map[string]MetricConfig{
	"keyspace_info": {
		Parser: Parsers{
			&versionMatchParser{
				verC: mustNewVersionConstraint(`<3.0.5`),
				Parser: &regexParser{
					name:   "keyspace_info_<3.0.5",
					reg:    regexp.MustCompile(`(?P<type>[^\s]*)\s*keys:(?P<keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`~3.0.5`),
				Parser: &regexParser{
					name: "keyspace_info_~3.0.5",
					reg: regexp.MustCompile(`(?P<type>\w*):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`~3.1.0`),
				Parser: &regexParser{
					name: "keyspace_info_~3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)_\s*(?P<type>[^:]+):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`3.2.0 - 3.3.2`),
				Parser: &regexParser{
					name: "keyspace_info_3.1.0-3.3.2",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`>=3.3.3`),
				Parser: &regexParser{
					name: "keyspace_info_>=3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invalid_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
		},
		MetricMeta: MetaDatas{
			{
				Name:      "keys",
				Help:      "pika serve instance total count of the db's key-type keys",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias, "db", "type"},
				ValueName: "keys",
			},
		},
	},

	"keyspace_info_all": {
		Parser: Parsers{
			&versionMatchParser{
				verC: mustNewVersionConstraint(`~3.0.5`),
				Parser: &regexParser{
					name: "keyspace_info_all_~3.0.5",
					reg: regexp.MustCompile(`(?P<type>\w*):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`~3.1.0`),
				Parser: &regexParser{
					name: "keyspace_info_all_~3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)_\s*(?P<type>[^:]+):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`3.2.0 - 3.3.2`),
				Parser: &regexParser{
					name: "keyspace_info_all_3.1.0-3.3.2",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`>=3.3.3`),
				Parser: &regexParser{
					name: "keyspace_info_all_>=3.3.3",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invalid_keys=(?P<invalid_keys>[\d]+)`),
					Parser: &normalParser{},
				},
			},
		},
		MetricMeta: MetaDatas{
			{
				Name:      "expire_keys",
				Help:      "pika serve instance total count of the db's key-type expire keys",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias, "db", "type"},
				ValueName: "expire_keys",
			},
			{
				Name:      "invalid_keys",
				Help:      "pika serve instance total count of the db's key-type invalid keys",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias, "db", "type"},
				ValueName: "invalid_keys",
			},
		},
	},
}
