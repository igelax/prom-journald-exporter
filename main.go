package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sync"

	"net/http"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// variables for Prometheus metrics
var (
	metricSudoCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sudo_count_total",
		Help: "The total number of sudo events",
	})
	debug bool
)

// sdjournal.JournalReader.Follow requires a custom type with a Write method (io.Write).
type JournalWriter struct{}

// Write method implementing the io.Write interface with JournalWriter pointer receiver
func (p *JournalWriter) Write(data []byte) (n int, err error) {
	// call JournalParser function with address of data to parse journald messages and process metrics
	JournalParser(&data)
	return len(data), nil
}

// JournalParser parses journal entry address and process Prometheus metrics (passed from Write method).
func JournalParser(entry *[]byte) {
	e := fmt.Sprintf("%s", *entry) // convert pointer value entry to string
	if debug {
		fmt.Printf("%s", e)
	}

	// Check entry using regexp and update Prometheus metrics
	r, _ := regexp.Compile("sudo:session")
	matched := r.MatchString(e)
	if matched {
		metricSudoCount.Inc() // increment Prometheus counter
	}
}

func main() {
	var (
		listenHTTP string
	)

	// command line (flag) variables and defaults
	flag.StringVar(&listenHTTP, "listenHTTP", ":9101", "ip:port to listen for http requests")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(2)
	go read_journal()        // go routine to follow/tail new journald logs and process metrics
	go prom_http(listenHTTP) // start Prometheus http endpoint
	wg.Wait()
}

// read_journal tails (follows) the journald log using the sdjournal package.
func read_journal() {

	// links for more information on sdjournal
	// https://pkg.go.dev/github.com/coreos/go-systemd/v22@v22.3.2/sdjournal#JournalReader
	// https://github.com/coreos/go-systemd/blob/v22.3.2/sdjournal/read.go

	// journal config
	jconf := sdjournal.JournalReaderConfig{
		Since: -1,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
				Value: "sudo", // ${APPNAME}.service
			},
		},
	}

	// journal reader
	jr, err := sdjournal.NewJournalReader(jconf)
	if err != nil {
		panic(err)
	}
	defer jr.Close()           // close JournalReader when done
	jrw := JournalWriter{}     // create variable of type JournalWriter (only implements Write method)
	err = jr.Follow(nil, &jrw) // follow journal and pass address of custom writer to parse new entries
	if err != nil {
		panic(err)
	}
}

// prom_http starts the Prometheus HTTP listener on the specified listen string at /metrics.
func prom_http(listen string) {
	fmt.Println("listening on", listen, "/metrics")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
