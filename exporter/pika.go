package exporter

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pourer/pika_exporter/discovery"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	infoSectionFlag       = "#"
	infoKeyValueSeparator = ":"
)

type dbKeyPair struct {
	db, key string
}

type pikaExproer struct {
	dis                 discovery.Discovery
	namespace           string
	metrics             Metrics
	keyPatterns, keys   []dbKeyPair
	scanCount           int
	keyValues, keySizes *prometheus.GaugeVec
	scrapeDuration      prometheus.Gauge
	scrapeErrors        prometheus.Gauge
	scrapeCount         prometheus.Counter
	mutex               *sync.Mutex
	wg                  sync.WaitGroup
	done                chan struct{}
}

func NewPikaExporter(dis discovery.Discovery, namespace string, metrics Metrics,
	keyPatterns, keys string, scanCount, statsClockHour int) (*pikaExproer, error) {
	e := &pikaExproer{
		dis:       dis,
		namespace: namespace,
		metrics:   metrics,
		mutex:     new(sync.Mutex),
		done:      make(chan struct{}),
	}

	var err error
	if e.keyPatterns, err = parseKeyArg(keyPatterns); err != nil {
		return nil, err
	}
	if e.keys, err = parseKeyArg(keys); err != nil {
		return nil, err
	}
	e.scanCount = scanCount

	e.initMetrics()
	e.wg.Add(1)
	go e.statsKeySpace(statsClockHour)
	return e, nil
}

func (e *pikaExproer) initMetrics() {
	for _, metric := range e.metrics {
		metric.GaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      metric.Name,
		}, metric.Labels)
	}
	e.keyValues = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "key_value",
	}, []string{"addr", "alias", "db", "key", "key_value"})
	e.keySizes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "key_size",
	}, []string{"addr", "alias", "db", "key", "key_type"})
	e.scrapeDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "exporter_last_scrape_duration_seconds",
	})
	e.scrapeErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "exporter_last_scrape_error",
	})
	e.scrapeCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: e.namespace,
		Name:      "exporter_scrape_count",
	})
}

func (e *pikaExproer) Close() error {
	close(e.done)
	e.wg.Wait()
	return nil
}

func (e *pikaExproer) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range e.metrics {
		metric.Describe(ch)
	}

	e.keyValues.Describe(ch)
	e.keySizes.Describe(ch)

	ch <- e.scrapeDuration.Desc()
	ch <- e.scrapeErrors.Desc()
	ch <- e.scrapeCount.Desc()
}

func (e *pikaExproer) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.initMetrics()
	e.keySizes.Reset()
	e.keyValues.Reset()

	e.scrape()

	e.keySizes.Collect(ch)
	e.keyValues.Collect(ch)
	for _, metric := range e.metrics {
		metric.Collect(ch)
	}

	ch <- e.scrapeDuration
	ch <- e.scrapeErrors
	ch <- e.scrapeCount
}

func (e *pikaExproer) scrape() {
	startTime := time.Now()
	errCount := 0

	e.scrapeCount.Inc()

	fut := newFuture()
	for _, instance := range e.dis.GetInstances() {
		fut.Add()
		go func(addr, password, alias string) {
			c, err := newClient(addr, password, alias)
			if err != nil {
				fut.Done(futureKey{addr: addr, alias: alias},
					fmt.Errorf("new pika client failed. err:%s", err.Error()))
			} else {
				defer c.Close()

				fut.Add()
				fut.Done(futureKey{addr: c.Addr(), alias: c.Alias()}, e.collectInfo(c))
				fut.Done(futureKey{addr: c.Addr(), alias: c.Alias()}, e.collectKeys(c))
			}
		}(instance.Addr, instance.Password, instance.Alias)
	}

	for k, v := range fut.Wait() {
		if v != nil {
			errCount++
			log.Errorf("collect pika failed. pika server:%#v err:%s", k, v.Error())
		}
	}

	e.scrapeErrors.Set(float64(errCount))
	e.scrapeDuration.Set(time.Now().Sub(startTime).Seconds())
}

func (e *pikaExproer) collectInfo(c *client) error {
	info, err := c.Info()
	if err != nil {
		return err
	}

	infoKeyValues, err := parseInfo(info)
	if err != nil {
		return err
	}
	infoKeyValues["addr"] = c.Addr()
	infoKeyValues["alias"] = c.Alias()

	for _, metric := range e.metrics {
		valid := true
		labelValues := make([]string, len(metric.Labels))
		for i, label := range metric.Labels {
			if v, ok := infoKeyValues[label]; ok {
				labelValues[i] = v
			} else {
				log.Debugf("no label value found. addr:%s metricName:%s label:%s", c.Addr(), metric.Name, label)
				valid = false
				break
			}
		}

		var value float64
		if valid && metric.ValueName != "" {
			if v, ok := infoKeyValues[metric.ValueName]; ok {
				if vv, err := convertValue(v); err != nil {
					log.Warnf("convert value to float64 failed. addr:%s metricName:%s valueName:%s value:%s",
						c.Addr(), metric.Name, metric.ValueName, v)
				} else {
					value = vv
				}
			} else {
				log.Debugf("no value found. addr:%s metricName:%s valueName:%s",
					c.Addr(), metric.Name, metric.ValueName)
				valid = false
			}
		}

		if valid {
			metric.WithLabelValues(labelValues...).Set(value)
		}
	}

	return nil
}

func (e *pikaExproer) collectKeys(c *client) error {
	allKeys := append([]dbKeyPair{}, e.keys...)
	keys, err := getKeysFromPatterns(c, e.keyPatterns, e.scanCount)
	if err != nil {
		log.Errorf("get keys from patterns failed. addr:%s err:%s", c.Addr(), err.Error())
	} else {
		allKeys = append(allKeys, keys...)
	}

	log.Debugf("collectKeys allKeys:%#v", allKeys)
	for _, k := range allKeys {
		if err := c.Select(k.db); err != nil {
			log.Warnf("couldn't select database %#v when getting key info. addr:", k.db, c.Addr())
			continue
		}

		keyInfo, err := c.Type(k.key)
		if err != nil {
			log.Warnf("get key info failed. addr:%s key:%s err:%s", c.Addr(), k.key, err.Error())
			continue
		}

		e.keySizes.WithLabelValues(c.Addr(), c.Alias(), "db"+k.db, k.key, keyInfo.keyType).Set(keyInfo.size)
		if value, err := c.Get(k.key); err == nil {
			e.keyValues.WithLabelValues(c.Addr(), c.Alias(), "db"+k.db, k.key, value).Set(1)
		}
	}

	return nil
}

func getKeysFromPatterns(c *client, keyPatterns []dbKeyPair, scanCount int) ([]dbKeyPair, error) {
	var expandedKeys []dbKeyPair
	for _, kp := range keyPatterns {
		if regexp.MustCompile(`[\?*\[\]\^]+`).MatchString(kp.key) {
			if err := c.Select(kp.db); err != nil {
				return expandedKeys, err
			}
			keyNames, err := c.Scan(kp.key, scanCount)
			if err != nil {
				log.Errorln("get keys from patterns scan failed. pattern:", kp.key)
				continue
			}
			for _, keyName := range keyNames {
				expandedKeys = append(expandedKeys, dbKeyPair{db: kp.db, key: keyName})
			}
		} else {
			expandedKeys = append(expandedKeys, kp)
		}
	}
	return expandedKeys, nil
}

func (e *pikaExproer) statsKeySpace(hour int) {
	defer e.wg.Done()

	if hour < 0 {
		log.Infoln("stats KeySpace not open")
		return
	}

	timer := time.NewTimer(getClockDuration(hour))
	defer timer.Stop()

	for {
		select {
		case <-e.done:
			return
		case <-timer.C:
			timer.Reset(getClockDuration(hour))
		}

		for _, v := range e.dis.GetInstances() {
			c, err := newClient(v.Addr, v.Password, v.Alias)
			if err != nil {
				log.Warnln("stats KeySpace new pika client failed. err:", err)
				continue
			}
			if _, err := c.InfoKeySpaceOne(); err != nil {
				log.Warnln("stats KeySpace execute INFO KEYSPACE 1 failed. err:", err)
			}
		}
	}
}

func parseKeyArg(keysArgString string) ([]dbKeyPair, error) {
	if keysArgString == "" {
		return nil, nil
	}

	var (
		keys []dbKeyPair
		err  error
	)
	for _, k := range strings.Split(keysArgString, ",") {
		db := "0"
		key := ""
		frags := strings.Split(k, "=")
		switch len(frags) {
		case 1:
			db = "0"
			key, err = url.QueryUnescape(strings.TrimSpace(frags[0]))
		case 2:
			db = strings.Replace(strings.TrimSpace(frags[0]), "db", "", -1)
			key, err = url.QueryUnescape(strings.TrimSpace(frags[1]))
		default:
			return keys, fmt.Errorf("invalid key list argument: %s", k)
		}
		if err != nil {
			return keys, fmt.Errorf("couldn't parse db/key string: %s", k)
		}

		keys = append(keys, dbKeyPair{db, key})
	}
	return keys, err
}

func parseInfo(info string) (map[string]string, error) {
	keyValues := make(map[string]string)
	lines := strings.Split(info, "\r\n")

	for _, line := range lines {
		line = strings.ToLower(line)
		cleanLine := cleanString(line)
		if cleanLine == "" {
			continue
		}
		if isSection(line) {
			continue
		}

		k, v := fetchKV(cleanLine)
		if checkSplit(k, v, keyValues) {
			continue
		}
		if v == "" {
			v = "null"
		}
		keyValues[k] = v
	}

	return keyValues, nil
}

func isSection(s string) bool {
	return s != "" &&
		strings.Index(s, infoSectionFlag) == 0 &&
		!strings.Contains(s, infoKeyValueSeparator)
}

func fetchKV(s string) (string, string) {
	pos := strings.Index(s, infoKeyValueSeparator)
	if pos < 0 {
		return s, ""
	}
	k := strings.Replace(s[:pos], " ", "_", -1)
	k = strings.Replace(k, "-", "_", -1)
	v := s[pos+1:]
	if strings.Index(v, " ") == 0 {
		v = strings.Replace(v, " ", "", 1)
	}
	return k, v
}

func cleanString(s string) string {
	s = strings.Replace(s, infoSectionFlag, "", 1)
	if strings.Index(s, " ") == 0 {
		s = strings.Replace(s, " ", "", 1)
	}
	return s
}

func getClockDuration(hour int) time.Duration {
	timeNow, timeDst := time.Now(), time.Now()
	subHour := hour - timeNow.Hour()
	if subHour <= 0 {
		timeDst = timeNow.AddDate(0, 0, 1).Add(time.Duration(subHour) * time.Hour).Truncate(time.Hour)
	} else {
		timeDst = timeNow.Add(time.Duration(subHour) * time.Hour).Truncate(time.Hour)
	}

	return timeDst.Sub(timeNow)
}
