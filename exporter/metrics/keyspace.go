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
					reg: regexp.MustCompile(`(?P<type>[^\s]*)\s*keys:(?P<keys>[\d]+)`),
					MetricMeta: &MetaData{
						Name:      "keys",
						Help:      "pika serve instance total count of the key-type keys",
						Type:      metricTypeGauge,
						Labels:    []string{LabelNameAddr, LabelNameAlias, "type"},
						ValueName: "keys",
					},
				},
			},

			&versionMatchParser{
				verC: mustNewVersionConstraint(`~3.0.5`),
				Parser: &regexParser{
					reg: regexp.MustCompile(`(?P<type>\w*):\s*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					MetricMeta: MetaDatas{
						{
							Name:      "keys",
							Help:      "pika serve instance total count of the key-type keys",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "type"},
							ValueName: "keys",
						},
						{
							Name:      "expire_keys",
							Help:      "pika serve instance total count of the key-type expire keys",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "type"},
							ValueName: "expire_keys",
						},
						{
							Name:      "invalid_keys",
							Help:      "pika serve instance total count of the key-type invalid keys",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "type"},
							ValueName: "invalid_keys",
						},
					},
				},
			},

			&versionMatchParser{
				verC: mustNewVersionConstraint(`>=3.1.0`),
				Parser: &regexParser{
					reg: regexp.MustCompile(`(?P<db>db[\d]+)\s*(?P<type>[^_]+)\w*keys=(?P<keys>[\d]+)[,\s]*` +
						`expires=(?P<expire_keys>[\d]+)[,\s]*invaild_keys=(?P<invalid_keys>[\d]+)`),
					MetricMeta: MetaDatas{
						{
							Name:      "keys",
							Help:      "pika serve instance total count of the db's key-type keys",
							Type:      metricTypeGauge,
							Labels:    []string{LabelNameAddr, LabelNameAlias, "db", "type"},
							ValueName: "keys",
						},
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
			},
		},
	},
}
