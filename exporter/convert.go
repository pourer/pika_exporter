package exporter

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	splitters = map[*regexp.Regexp]func(s string, m map[string]string){
		regexp.MustCompile("binlog_offset"):           binlogOffsetSplitter,
		regexp.MustCompile("is_bgsaving"):             isBgSavingSplitter,
		regexp.MustCompile("is_slots_reloading"):      isSlotsReloadingSplitter,
		regexp.MustCompile("is_slots_cleanuping"):     isSlotCleaningSplitter,
		regexp.MustCompile("^slave\\d+$"):             slaveSplitter,
		regexp.MustCompile("^time$"):                  keySpaceTimeSplitter,
		regexp.MustCompile("double_master_recv_info"): doubleMasterRecvInfoSplitter,
	}
)

func checkSplit(k, v string, m map[string]string) bool {
	for re, f := range splitters {
		if re.Match([]byte(k)) {
			f(v, m)
			return true
		}
	}
	return false
}

func binlogOffsetSplitter(s string, m map[string]string) {
	vs := make([]string, 2)
	ss := strings.Split(s, " ")
	copy(vs, ss)

	m["binlog_offset_filenum"] = vs[0]
	m["binlog_offset_value"] = vs[1]
}

func isBgSavingSplitter(s string, m map[string]string) {
	vs := make([]string, 3)
	ss := strings.Split(s, ",")
	copy(vs, ss)

	m["is_bgsaving"] = vs[0]
	m["bgsave_start_time"] = vs[1]
}

func isSlotsReloadingSplitter(s string, m map[string]string) {
	vs := make([]string, 3)
	ss := strings.Split(s, ",")
	copy(vs, ss)

	m["is_slots_reloading"] = vs[0]
	m["slots_reload_start_time"] = vs[1]
}

func isSlotCleaningSplitter(s string, m map[string]string) {
	vs := make([]string, 3)
	ss := strings.Split(s, ",")
	copy(vs, ss)

	m["is_slots_cleaning"] = vs[0]
	m["slots_clean_start_time"] = vs[1]
}

func slaveSplitter(s string, m map[string]string) {
	vs := make([]string, 5)
	ss := strings.Split(s, ",")

	for i, sv := range ss {
		if i >= len(vs) {
			break
		}
		svs := strings.Split(sv, "=")
		if len(svs) > 1 {
			vs[i] = svs[1]
		}
	}

	m["slave_ip"] = vs[0]
	m["slave_port"] = vs[1]
	m["slave_state"] = vs[2]
	m["slave_sid"] = vs[3]
	m["slave_lag"] = vs[4]
}

func keySpaceTimeSplitter(s string, m map[string]string) {
	m["keyspace_time"] = s
}

func doubleMasterRecvInfoSplitter(s string, m map[string]string) {
	vs := make([]string, 4)
	ss := strings.Split(s, " ")
	copy(vs, ss)

	m["double_master_recv_info_binlog_filenum"] = vs[1]
	m["double_master_recv_info_binlog_offset"] = vs[3]
}

func convertValue(s string) (float64, error) {
	switch s {
	case "yes", "up", "online":
		return 1, nil
	case "no", "down", "offline", "null":
		return 0, nil
	default:
		return strconv.ParseFloat(s, 64)
	}
}
