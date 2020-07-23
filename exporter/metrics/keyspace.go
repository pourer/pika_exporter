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
					name: "keyspace_info_<3.0.5",
					reg:  regexp.MustCompile(`(?P<type>[^\s]*)\s*keys:(?P<keys>[\d]+)`),
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`~3.0.5`),
				Parser: &regexParser{
					name: "keyspace_info_~3.0.5",
					reg: regexp.MustCompile(`(?P<type>\w*):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`3.1.0 - 3.3.2`),
				Parser: &regexParser{
					name: "keyspace_info_>=3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`>=3.3.3`),
				Parser: &regexParser{
					name: "keyspace_info_>=3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invalid_keys=(?P<invalid_keys>[\d]+)`),
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
					name: "keyspace_info_~3.0.5",
					reg: regexp.MustCompile(`(?P<type>\w*):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`3.1.0 - 3.3.2`),
				Parser: &regexParser{
					name: "keyspace_info_>=3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
				},
			},
			&versionMatchParser{
				verC: mustNewVersionConstraint(`>=3.3.3`),
				Parser: &regexParser{
					name: "keyspace_info_>=3.1.0",
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invalid_keys=(?P<invalid_keys>[\d]+)`),
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
