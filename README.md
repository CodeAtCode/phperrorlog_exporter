# phperrorlog_exporter

Prometheus php error log exporter, watches php error logs end exports metrics to prometheus.  
Based on [foomo/phperrorlog_exporter](https://github.com/foomo/phperrorlog_exporter), updated to latest Go version and dependence with a better logging and readme.

![Screenshot_20210506_151755](https://user-images.githubusercontent.com/403283/117305468-1321e480-ae7f-11eb-97b8-c6de95a02ee5.png)  
Grafana dashboard sample available [there](https://grafana.com/grafana/dashboards/14368).

## Compile

```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
go get github.com/CodeAtCode/phperrorlog_exporter
cd phperrorlog_exporter
go build phperrorlog_exporter.go
```

### Output

```
$ phperrorlog_exporter test.com:path/to/log prova:/another/path/log
Listening on address:port =>  :9423
observing test.com stats map[fatal:1 notice:2110 warning:1]
observing prova stats map[fatal:1 notice:1877 warning:1]
```

## Grafana

```
  - job_name: 'phperrorlog-exporter'
    scrape_interval: 30m
    scrape_timeout: 30m
    metrics_path: /metrics
    static_configs:
    - targets:
      - default
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - target_label: __address__
        replacement: "localhost:9423"
      - source_labels: [__param_target]
        target_label: instance
```
