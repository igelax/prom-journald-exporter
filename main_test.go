package main

import (
	_ "bufio"
	_ "context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "strings"
	"sync"
	"testing"
	"time"
)

func testRoutine(t *testing.T) {
	time.Sleep(time.Second * 2)

	resp, err := http.Get("http://localhost:9101/metrics")
	if err != nil {
		t.Error("Expected prom http listener to be running")
	}

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Expected body to be valid")
	}

	// Convert the body to type string
	sb := string(body)

	// scanner := bufio.NewScanner(strings.NewReader(sb))
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	fmt.Print(sb)
	panic(0)

}

func TestProm(t *testing.T) {

	var (
		// ctx, cancel = context.WithCancel(context.Background())
		listenHTTP string
	)

	// command line (flag) variables and defaults
	flag.StringVar(&listenHTTP, "listenHTTP", ":9101", "ip:port to listen for http requests")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(3)
	go read_journal()        // go routine to follow/tail new journald logs and process metrics
	go prom_http(listenHTTP) // start Prometheus http endpoint
	go testRoutine(t)
	wg.Wait()

}
