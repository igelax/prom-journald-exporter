package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

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
	verbose bool
)

// sdjournal.JournalReader.Follow requires a custom type with a Write method (io.Write).
type JournalWriter struct{}

// Write method implementing the io.Write interface with JournalWriter pointer receiver
func (p *JournalWriter) Write(data []byte) (n int, err error) {
	// call JournalParser function with address of data to parse journald messages and process metrics
	JournalParser(&data)
	return len(data), nil
}

// JournalParser parses journal entry address and processes Prometheus metrics (passed from Write method).
func JournalParser(entry *[]byte) {
	e := fmt.Sprintf("%s", *entry) // convert pointer value entry to string
	if verbose {
		fmt.Printf("%s", e)
	}

	// Check entry using regexp and update Prometheus metrics
	r, _ := regexp.Compile("sudo:session")
	matched := r.MatchString(e)
	if matched {
		metricSudoCount.Inc() // increment Prometheus counter
		if verbose {
			fmt.Println("incremented prometheus counter")
		}
	}
}

func main() {
	var (
		listenHTTP  string
		syslog_id   string
		ctx, cancel = context.WithCancel(context.Background())
	)

	defer func() {
		// Close database, redis, truncate message queues, etc.
		fmt.Println("Running cancel()")
		cancel()
	}()

	// command line (flag) variables and defaults
	flag.StringVar(&listenHTTP, "listenHTTP", ":9101", "ip:port to listen for http requests")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output with filtered log entries")
	flag.StringVar(&syslog_id, "syslogIdentifier", "sudo", "syslog identifier used to filter journald (empty string for no filtering")
	flag.Parse()

	go read_journal(ctx, syslog_id) // go routine to follow/tail new journald logs and process metrics
	go prom_http(ctx, listenHTTP)   // start Prometheus http endpoint

	// Wait for SIGINT.
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGHUP)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig

	// Shutdown. Cancel application context will kill all attached tasks.
	cancel()
}

// read_journal tails (follows) the journald log using the sdjournal package.
func read_journal(ctx context.Context, syslogIdentifier string) {

	defer func() {
		fmt.Println("running deferred ctx.Done")
		<-ctx.Done()
	}()

	// links for more information on sdjournal
	// https://pkg.go.dev/github.com/coreos/go-systemd/v22@v22.3.2/sdjournal#JournalReader
	// https://github.com/coreos/go-systemd/blob/v22.3.2/sdjournal/read.go

	// journal config
	jconf := sdjournal.JournalReaderConfig{
		Since: -1,
	}

	// Add Match rule
	if syslogIdentifier != "" {
		jconf.Matches = []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
				Value: syslogIdentifier, // ${APPNAME}.service
			},
		}
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
func prom_http(ctx context.Context, listen string) {

	defer func() {
		fmt.Println("running deferred ctx.Done")
		<-ctx.Done()
	}()

	fmt.Println("listening on", listen, "/metrics")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
