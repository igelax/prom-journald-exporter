# journald parser and Prometheus exporter

Export Prometheus metrics from journald events using Prometheus Go client library.  For demonstration purposes, journald is filtered for the `sudo` syslog identifier and a basic Prometheus counter metric is incremented.

## Build

```bash
go get github.com/msgarbossa/prom-journald-exporter
cd $GOPATH/src/github.com/msgarbossa/prom-journald-exporter
go build
```

### Makefile

The Makefile has the following targets:
- clean
- build_amd64
- build_arm64
- all

To cross-compile ARM64 on AMD64, install the following Debian package (Ubuntu 20.04):
- gcc-aarch64-linux-gnu

To cross-compile ARM on AMD64, install the following Debian package (Ubuntu 20.04):
- gcc-arm-linux-gnueabihf


## Download

The linux-amd64 build can be downloaded from [releases](https://github.com/msgarbossa/prom-journald-exporter/releases)

## Example

When started with the -debug option, matching journald entries are printed to stdout.

```bash
$ ./prom_journald_exporter -h
Usage of ./prom_journald_exporter:
  -verbose
    	Enable verbose output
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

## Testing

The TestProm function in main_test.go uses the logger shell command to send a journald test message and then queries the listener for the Prometheus exporter to look for the incremented counter.
```
go test
```

Debug journal logs (print filtered messages)
```
go run main.go -verbose
```
Then generate journald messages, which will be printed to stdout as they are parsed for Prometheus metrics.
