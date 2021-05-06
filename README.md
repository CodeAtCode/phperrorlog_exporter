# phperrorlog_exporter

Prometheus php error log exporter, watches php error logs end exports metrics to prometheus.  
Based on [foomo/phperrorlog_exporter](https://github.com/foomo/phperrorlog_exporter), updated to latest Go version and dependence with a better logging and readme.

## Compile

```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
go get github.com/CodeAtCode/phperrorlog_exporter
cd phperrorlog_exporter
go build phperrorlog_exporter.go
```
