package metrics

func init() {
	Register(collectCPUMetrics)
}

var collectCPUMetrics = map[string]MetricConfig{
	"used_cpu_sys": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "used_cpu_sys",
			Help:      "pika serve instance total count of used cpu sys",
			Type:      metricTypeCounter,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "used_cpu_sys",
		},
	},
	"used_cpu_user": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "used_cpu_user",
			Help:      "pika serve instance total count of used cpu user",
			Type:      metricTypeCounter,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "used_cpu_user",
		},
	},
	"used_cpu_sys_children": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "used_cpu_sys_children",
			Help:      "pika serve instance children total count of used cpu sys",
			Type:      metricTypeCounter,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "used_cpu_sys_children",
		},
	},
	"used_cpu_user_children": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "used_cpu_user_children",
			Help:      "pika serve instance children total count of used cpu user",
			Type:      metricTypeCounter,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "used_cpu_user_children",
		},
	},
}
