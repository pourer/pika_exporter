package exporter

import (
	"testing"
	"github.com/Masterminds/semver"
	"github.com/pourer/pika_exporter/exporter/test"
	"github.com/pourer/pika_exporter/exporter/metrics"
	"github.com/stretchr/testify/assert"
)

func mustNewVersionConstraint(version string) *semver.Constraints {
	c, err := semver.NewConstraint(version)
	if err != nil {
		panic(err)
	}
	return c
}

func Test_Parse_Info(t *testing.T) {
	for _, infoCase := range test.InfoCases {
		version, extracts, err := parseInfo(infoCase.Info)
		if err != nil {
			t.Errorf("%s parse info fialed. err:%s", infoCase.Name, err.Error())
		}

		extracts[metrics.LabelNameAddr] = "127.0.0.1"
		extracts[metrics.LabelNameAlias] = ""

		collector := metrics.CollectFunc(func(m metrics.Metric) error {
			t.Logf("metric:%#v", m)
			return nil
		})
		parseOpt := metrics.ParseOption{
			Version:  version,
			Extracts: extracts,
			Info:     infoCase.Info,
		}
		t.Logf("##########%s begin parse###########", infoCase.Name)
		for _, metricConfig := range metrics.MetricConfigs {
			metricConfig.Parse(collector, parseOpt)
		}
	}
}

func Test_Parse_Version_Error(t *testing.T) {
	assert := assert.New(t)

	info := `# Server
pika_version:aaa
pika_git_sha:b22b0561f9093057d2e2d5cc783ff630fb2c8884
pika_build_compile_date: Nov  7 2019
os:Linux 3.10.0-1062.9.1.el7.x86_64 x86_64`

	version, _, err := parseInfo(info)
	assert.Nil(version)
	assert.Error(err)
}

func Test_Parse_Extracts_Error(t *testing.T) {
	testCases := []struct {
		info         string
		checkVersion *semver.Constraints
		checkKeys    map[string]string
	}{
		{
			info: `# Server
pika_version:3.0.10
role-slave`,
			checkVersion: mustNewVersionConstraint(`~3.0.5`),
			checkKeys: map[string]string{
				"role": "slave",
			},
		},
		{
			info: `# Server
pika_version:3.0.10
:role:slave`,
			checkVersion: mustNewVersionConstraint(`~3.0.5`),
			checkKeys: map[string]string{
				"role": "slave",
			},
		},
	}

	for _, testCase := range testCases {
		version, extracts, err := parseInfo(testCase.info)
		if err != nil {
			t.Error(err)
		}

		if !testCase.checkVersion.Check(version) {
			t.Error("version not ok")
		}

		for k, v := range testCase.checkKeys {
			if vv, ok := extracts[k]; ok && vv == v{
				t.Errorf("not found key:%s", k)
			}
		}

		extracts[metrics.LabelNameAddr] = "127.0.0.1"
		extracts[metrics.LabelNameAlias] = ""

		collector := metrics.CollectFunc(func(m metrics.Metric) error {
			t.Errorf("metric:%#v shouldn't collect", m)
			return nil
		})
		parseOpt := metrics.ParseOption{
			Version:  version,
			Extracts: extracts,
			Info:     testCase.info,
		}
		for _, metricConfig := range metrics.MetricConfigs {
			metricConfig.Parse(collector, parseOpt)
		}
	}
}

func Benchmark_Parse(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			info := test.V236MasterInfo

			version, extracts, err := parseInfo(info)
			if err != nil {
				b.Error(err)
			}

			extracts[metrics.LabelNameAddr] = "127.0.0.1"
			extracts[metrics.LabelNameAlias] = ""

			collector := metrics.CollectFunc(func(m metrics.Metric) error {
				return nil
			})
			parseOpt := metrics.ParseOption{
				Version:  version,
				Extracts: extracts,
				Info:     info,
			}
			for _, metricConfig := range metrics.MetricConfigs {
				metricConfig.Parse(collector, parseOpt)
			}
		}
	})
}
