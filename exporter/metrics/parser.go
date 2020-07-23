package metrics

import (
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/Masterminds/semver"
)

const (
	defaultValue = 0
)

type ParseOption struct {
	Version  *semver.Version
	Extracts map[string]string
	Info     string
}

type Parser interface {
	Parse(m MetricMeta, c Collector, opt ParseOption)
}

type Parsers []Parser

func (ps Parsers) Parse(m MetricMeta, c Collector, opt ParseOption) {
	for _, p := range ps {
		p.Parse(m, c, opt)
	}
}

type versionMatchParser struct {
	verC *semver.Constraints
	Parser
}

func (p *versionMatchParser) Parse(m MetricMeta, c Collector, opt ParseOption) {
	if opt.Version == nil || !p.verC.Check(opt.Version) {
		return
	}
	p.Parser.Parse(m, c, opt)
}

type keyMatchParser struct {
	matches map[string]string
	Parser
}

func (p *keyMatchParser) Parse(m MetricMeta, c Collector, opt ParseOption) {
	for key, matchValue := range p.matches {
		if v, _ := opt.Extracts[key]; strings.ToLower(v) != strings.ToLower(matchValue) {
			return
		}
	}
	p.Parser.Parse(m, c, opt)
}

type normalParser struct{}

func (p *normalParser) Parse(m MetricMeta, c Collector, opt ParseOption) {
	m.Lookup(func(m MetaData) {
		metric := Metric{
			MetaData:    m,
			LabelValues: make([]string, len(m.Labels)),
			Value:       defaultValue,
		}

		for i, labelName := range m.Labels {
			labelValue, ok := findInMap(labelName, opt.Extracts)
			if !ok {
				log.Debugf("normalParser::Parse not found label value. metricName:%s labelName:%s",
					m.Name, labelName)
			}

			metric.LabelValues[i] = labelValue
		}

		if m.ValueName != "" {
			if v, ok := findInMap(m.ValueName, opt.Extracts); !ok {
				log.Warnf("normalParser::Parse not found value. metricName:%s valueName:%s", m.Name, m.ValueName)
				return
			} else {
				metric.Value = convertToFloat64(v)
			}
		}

		if err := c.Collect(metric); err != nil {
			log.Errorf("normalParser::Parse metric collect failed. metric:%#v err:%s",
				m, m.ValueName)
		}
	})
}

type regexParser struct {
	name string
	reg  *regexp.Regexp
}

func (p *regexParser) Parse(m MetricMeta, c Collector, opt ParseOption) {
	matchMaps := p.regMatchesToMap(opt.Info)

	m.Lookup(func(m MetaData) {
		for _, matches := range matchMaps {
			metric := Metric{
				MetaData:    m,
				LabelValues: make([]string, len(m.Labels)),
				Value:       defaultValue,
			}

			for i, labelName := range m.Labels {
				labelValue, ok := findInMap(labelName, matches, opt.Extracts)
				if !ok {
					log.Debugf("regexParser::Parse not found label value. metricName:%s labelName:%s",
						m.Name, labelName)
				}

				metric.LabelValues[i] = labelValue
			}

			if m.ValueName != "" {
				if v, ok := findInMap(m.ValueName, matches, opt.Extracts); !ok {
					log.Warnf("regexParser::Parse not found value. metricName:%s valueName:%s",
						m.Name, m.ValueName)
					return
				} else {
					metric.Value = convertToFloat64(v)
				}
			}

			if err := c.Collect(metric); err != nil {
				log.Errorf("regexParser::Parse metric collect failed. metric:%#v err:%s",
					m, m.ValueName)
			}
		}
	})
}

func (p *regexParser) regMatchesToMap(s string) []map[string]string {
	if s == "" {
		return nil
	}

	multiMatches := p.reg.FindAllStringSubmatch(s, -1)
	if len(multiMatches) == 0 {
		log.Errorf("regexParser::Parse reg find sub match nil. name:%s text:%s", p.name, s)
		return nil
	}

	ms := make([]map[string]string, len(multiMatches))
	for i, matches := range multiMatches {
		ms[i] = make(map[string]string)
		for j, name := range p.reg.SubexpNames() {
			ms[i][name] = trimSpace(matches[j])
		}
	}
	return ms
}

func findInMap(key string, ms ...map[string]string) (string, bool) {
	for _, m := range ms {
		if v, ok := m[key]; ok {
			return v, true
		}
	}
	return "", false
}

func trimSpace(s string) string {
	return strings.TrimRight(strings.TrimLeft(s, " "), " ")
}

func convertToFloat64(s string) float64 {
	s = strings.ToLower(s)

	switch s {
	case "yes", "up", "online":
		return 1
	case "no", "down", "offline", "null":
		return 0
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return n
}

func mustNewVersionConstraint(version string) *semver.Constraints {
	c, err := semver.NewConstraint(version)
	if err != nil {
		panic(err)
	}
	return c
}
