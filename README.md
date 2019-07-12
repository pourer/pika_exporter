# Pika Metric Exporter #

Prometheus expoter for nosql [Qihoo360/pika](https://github.com/Qihoo360/pika) metrics. Suppots Pika 2.x,3.x

Pika-Expoter is based on [Redis-Exporter](https://github.com/oliver006/redis_exporter)).

## Buiding ##

**Build and run locally:**  

To start using `pika_expoter`, install `Go` and run go get
```
$ go get github.com/pourer/pika_expoter
$ cd $GOPATH/src/github.com/pourer/pika_expoter
$ make
$ ./bin/pika_exporter <flags>
```


**Prometheus Configuration:**

Add a block to the scrape_configs of your prometheus.yml config file:
```
scrape_configs:

...

- job_name: pika
  scrape_interval: 15s
  static_configs:
  - targets: ['XXXXXX:9121']
    labels:
      group: 'test'
      
...
```

## Flags ##
Name | Environment Variables | Default | Description | Example
---|---|---|---|---
pika.host-file | PIKA_HOST_FILE | | Path to file containing one or more pika nodes, separated by newline. NOTE: mutually exclusive with pika.addr.Each line can optionally be comma-separated with the fields `<addr>`,`<password>`,`<alias>`. See [here](https://github.com/pourer/pika_exporter/raw/master/contrib/sample_pika_hosts_file.txt) for an example file.| --pika.host-file ./pika_hosts_file.txt
pika.addr | PIKA_ADDR | | Address of one or more pika nodes, separated by comma. | --pika.addr 192.168.1.2:9221,192.168.1.3:9221
pika.password | PIKA_PASSWORD | | Password for one or more pika nodes, separated by comma. | --pika.password 123.com,123.com
pika.alias | PIKA_ALIAS | | Pika instance alias for one or more pika nodes, separated by comma. | --pika.alias a,b
namespace | PIKA_EXPORTER_NAMESPACE | pika | Namespace for metrics | --namespace pika
metrics-file | PIKA_EXPORTER_METRICS_FILE | | Metrics definition file. Define the metrics info that need to collect. See details: [INFO Metrics Definition](#INFO-Metrics-Definition).| --metrics-file ./metrics_definition_file.ini
keyspace-stats-clock | PIKA_EXPORTER_KEYSPACE_STATS_CLOCK | -1 | Stats the number of keys at keyspace-stats-clock o'clock every day, in the range [0, 23]. If < 0, not open this feature. | --keyspace-stats-clock 0
check.key-patterns | PIKA_EXPORTER_CHECK_KEY_PARTTERNS | | Comma separated list of key-patterns to export value and length/size, searched for with SCAN. | --check.key-patterns db0=test*,db0=*abc*
check.keys | PIKA_EXPORTER_CHECK_KEYS | | Comma separated list of keys to export value and length/size. | --check.keys abc,test,wasd
check.scan-count | PIKA_EXPORTER_CHECK_SCAN_COUNT | 100 | When check keys and executing SCAN command, scan-count assigned to COUNT. | --check.scan-count 200
web.listen-address | PIKA_EXPOTER_WEB_LISTEN_ADDRESS | :9121 | Address to listen on for web interface and telemetry. | --web.listen-address ":9121"
web.telemetry-path | PIKA_EXPORTER_WEB_TELEMETRY_PATH | /metrics | Path under which to expose metrics. | --web.telemetry-path "/metrics"
log.level | PIKA_EXPORTER_LOG_LEVEL | info | Log level, valid options: `panic` `fatal` `error` `warn` `warning` `info` `debug`. | --log.level "debug"
log.format | PIKA_EXPORTER_LOG_FORMAT | json | Log format, valid options: `txt` `json`. | --log.format "json"
version | | false | Show version information and exit. | --version

## INFO Metrics Definition ##
The `keys-list` obtained by the **`INFO`** Command are as follows:

Key Name | Corresponding to the key in INFO
--- | ---
addr | By configuration: --pika.host-file or --pika.addr
alias | By configuration: --pika.host-file or --pika.alias
pika_version | pika_version
pika_git_sha | pika_git_sha
pika_build_compile_date | pika_build_compile_date
os | os
arch_bits | arch_bits
process_id | process_id
tcp_port | tcp_port
thread_num | thread_num
sync_thread_num | sync_thread_num
uptime_in_seconds | uptime_in_seconds
config_file | config_file
server_id | server_id
db_size | db_size
compression | compression
used_memory | used_memory
db_memtable_usage | db_memtable_usage
db_tablereader_usage | db_tablereader_usage
log_size | log_size
safety_purge | safety_purge
expire_logs_days | expire_logs_days
expire_logs_nums | expire_logs_nums
binlog_offset_filenum | binlog_offset format is: 0 388(file-number offset), pika_exporter separated it, the key of "file-number" is: binlog_offset_filenum.
binlog_offset_value | binlog_offset format is: 0 388(file-number offset), pika_exporter separated it, the key of "offset" is: binlog_offset_value.
connected_clients | connected_clients
total_connections_received | total_connections_received
instantaneous_ops_per_sec | instantaneous_ops_per_sec
total_commands_processed | total_commands_processed
is_bgsaving | is_bgsaving format is: No, ,0(whether in backup, latest backup start time, how long has been backed up), pika_exporter separated it, the key of "whether in backup" is: is_bgsaving.
bgsave_start_time | is_bgsaving format is: No, ,0(whether in backup, latest backup start time, how long has been backed up), pika_exporter separated it, the key of "latest backup start time" is: bgsave_start_time.
is_slots_reloading | is_slots_reloading format is: No, ,0(whether in slots reload, latest reload start time, how long has been reload), pika_exporter separated it, the key of "whether in slots reload" is: is_slots_reloading.
slots_reload_start_time | is_slots_reloading format is: No, ,0(whether in slots reload, latest reload start time, how long has been reload), pika_exporter separated it, the key of "latest reload start time" is: slots_reload_start_time.
is_slots_cleaning | is_slots_reloading format is: No, ,0(whether in slots clean, latest clean start time, how long has been clean), pika_exporter separated it, the key of "whether in slots clean" is: is_slots_cleaning.
slots_clean_start_time | is_slots_reloading format is: No, ,0(whether in slots clean, latest clean start time, how long has been clean), pika_exporter separated it, the key of "latest clean start time" is: slots_clean_start_time.
is_scaning_keyspace | is_scaning_keyspace
is_compact | is_compact
compact_cron | compact_cron
compact_interval | compact_interval
used_cpu_sys | used_cpu_sys
used_cpu_user | used_cpu_user
used_cpu_sys_children | used_cpu_sys_children
used_cpu_user_children | used_cpu_user_children
role | role
connected_slaves | connected_slaves
slave_ip | The information about the slave of the master obtained by the INFO command is: slave0:ip=192.168.1.1,port=57765,state=online,sid=2,lag=0, pika_exporter treats it as: slave_ip slave_port slave_state slave_sid slave_lag
slave_port | The information about the slave of the master obtained by the INFO command is: slave0:ip=192.168.1.1,port=57765,state=online,sid=2,lag=0, pika_exporter treats it as: slave_ip slave_port slave_state slave_sid slave_lag
slave_sid | The information about the slave of the master obtained by the INFO command is: slave0:ip=192.168.1.1,port=57765,state=online,sid=2,lag=0, pika_exporter treats it as: slave_ip slave_port slave_state slave_sid slave_lag
slave_state | The information about the slave of the master obtained by the INFO command is: slave0:ip=192.168.1.1,port=57765,state=online,sid=2,lag=0, pika_exporter treats it as: slave_ip slave_port slave_state slave_sid slave_lag
slave_lag | The information about the slave of the master obtained by the INFO command is: slave0:ip=192.168.1.1,port=57765,state=online,sid=2,lag=0, pika_exporter treats it as: slave_ip slave_port slave_state slave_sid slave_lag
master_host | master_host
master_port | master_port
master_link_status | master_link_status
slave_priority | slave_priority
slave_read_only | slave_read_only
repl_state | repl_state
the_peer_master_host | the peer-master host
the_peer_master_port | the peer-master port
the_peer_master_server_id | the peer-master server_id
double_master_recv_info_binlog_filenum | double_master_recv_info format is: 0 0(file-number offset), pika_exporter separated it, the key of "file-number" is: double_master_recv_info_binlog_filenum.
double_master_recv_info_binlog_offset | double_master_recv_info format is: 0 0(file-number offset), pika_exporter separated it, the key of "offset" is: double_master_recv_info_binlog_filenum.
keyspace_time | The latest statistical time of the keyspace obtained by the INFO command is: # Time:1970-01-01 00:00:00, pika_exporter treats it as a Unix integer timestamp: keyspace_time.
kv_keys | The number of kv of the keyspace obtained by the INFO command is: kv keys: 0, pika_exporter treats it as: kv_keys:0
hash_keys | The number of hash-kv of the keyspace obtained by the INFO command is: hash keys: 0, pika_exporter treats it as: hash_keys:0
list_keys | The number of list-kv of the keyspace obtained by the INFO command is: list keys: 0, pika_exporter treats it as: list_keys:0
set_keys | The number of set-kv of the keyspace obtained by the INFO command is: set keys: 0, pika_exporter treats it as: set_keys:0
zset_keys | The number of zset-kv of the keyspace obtained by the INFO command is: zset keys: 0, pika_exporter treats it as: zset_keys:0


You can customize the metrics definition file as as needed. The format of the metrics definition file: `.ini`. See [here](https://github.com/pourer/pika_exporter/raw/master/contrib/default_pika_metrics_file.ini)

For example:
```
[uptime_in_seconds]
labels = addr,alias
value = uptime_in_seconds
```
Described as follows:
> `Section` in the configuration file is the name of the collection indicator. The view in prometheus is: `namesapce`_section, for example(the `default namespace` is `pika`): pika_build_info

> Each collection indicator needs to be configured with `labels` and `value`,if `value` is not configured, the collection indicator value defaults to `0`.

> The name, labels, and indicator value of the collection indicator are as follows:
1. The configured labels or value must exist in the `keys-list` obtained by the `INFO` command (subject to the key list after disassembly). If not, the collection indicator is not collected;
2. If labels or value exists in the key-list, but its corresponding value is empty, the assignment is:`null`;
3. All collection indicator names, collection label names, and label values ​​are converted to: lowercase letters;
4. The label value is not changed except for conversion to lowercase letters;
5. Numeric format conversion of indicator values(not case sensitive):
    1) yes => 1 no => 0
    2) up => 1 down => 0
    3) online => 1 offline => 0
    4) null => 0
    5) The rest of the situation directly converts the string to float
6. If the index value is converted to float, the value of the corresponding collection indicator is: 0.

The metrics information defined in the [default_pika_metrics_file.ini](https://github.com/pourer/pika_exporter/raw/master/contrib/default_pika_metrics_file.ini) file is used as the **`default collection standard`**. 

## Keys Metrics Definition ##
You can export values of keys if they're in numeric format by using the --check.key-patterns or --check.keys flag. The pika_exporter will export the size (or, depending on the data type, the length) and the value of the key.

The name of the collection indicator: 
- **`namespace_key_value`**  
  Only the value of the key obtained by the GET command
  
- **`namespace_key_size`**  
  When the PFCOUNT command is used to obtain the size of a key from Pika, even if the key is in `KV-Structure`, the return value can be obtained normally and no error message is received. 
Since `Hyperloglog` is not commonly used, then the key in `Hyperloglog-Structure` is not supported. The key in `Hyperloglog-Structure` will be treated as a key in `KV-Structure`.

**Please note**: 
> Since pika allows duplicate names five times, `SCAN` Command has a priority output order, followed by: string -> hash -> list -> zset -> set;

> Since pika allows the name to be renamed five times, the `TYPE` Command has the priority output order, which is: string -> hash -> list -> zset -> set. If the key exists in the string, then only the string is output. If it does not exist, Then output the hash, and so on.

## Grafana Dashboard ##

See [here](https://github.com/pourer/pika_exporter/raw/master/contrib/grafana_prometheus_pika_dashboard.json)

Screenshots:  
![Overview](https://github.com/pourer/pika_exporter/raw/master/contrib/overview.png)

![BaseInfo](https://github.com/pourer/pika_exporter/raw/master/contrib/base_info.png)

![Replication](https://github.com/pourer/pika_exporter/raw/master/contrib/replication.png)

![TimeConsumingOperation](https://github.com/pourer/pika_exporter/raw/master/contrib/time_consuming_operation.png)

![KeysMetrics](https://github.com/pourer/pika_exporter/raw/master/contrib/keys_metrics.png)