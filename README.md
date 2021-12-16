# journald parser and Prometheus exporter

Export Prometheus metrics from journald events using Prometheus Go client library.  For demonstration purposes, journald is filtered for the `sudo` syslog identifier and a basic Prometheus counter metric is incremented.

## Build

```bash
go get github.com/msgarbossa/prom-journald-exporter
cd $GOPATH/src/github.com/msgarbossa/prom-journald-exporter
go build
```

## Download

The linux-amd64 build can be downloaded from [releases](https://github.com/msgarbossa/prom-journald-exporter/releases)

## Example

When started with the -debug option, matching journald entries are printed to stdout.

```bash
$ ./prom_journald_exporter -h
Usage of ./prom_journald_exporter:
  -debug
    	Enable debug
  -listenHTTP string
    	ip:port to listen for http requests (default ":9101")

$ ./prom_journald_exporter -debug
listening on :9101 /metrics

```

Run a sudo command from another session and verify metric counts.

```bash
$ curl -s http://localhost:9101/metrics | grep sudo
# HELP sudo_count_total The total number of sudo events
# TYPE sudo_count_total counter
sudo_count_total 0

$ sudo ls > /dev/null

$ curl -s http://localhost:9101/metrics | grep sudo
# HELP sudo_count_total The total number of sudo events
# TYPE sudo_count_total counter
sudo_count_total 2
```
