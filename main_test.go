package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

func testRoutine(t *testing.T) {

	time.Sleep(time.Second * 2)

	cmd := exec.Command("logger", "sudo:session")

	err := cmd.Run()

	if err != nil {
		t.Error("Expected logger shell command to generate a log message, got", err)
	}

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

	// Scan results for regexp match indicating test entry was counted by Prometheus
	logMatch := false
	r, _ := regexp.Compile("^sudo_count_total 1$")
	scanner := bufio.NewScanner(strings.NewReader(sb))
	for scanner.Scan() {
		matched := r.MatchString(scanner.Text())
		if matched {
			logMatch = true
		}
		//fmt.Println(scanner.Text())
	}
	if !logMatch {
		t.Error("Expected to find test journald log message counted in Prometheus metrics")
	}

}

func TestProm(t *testing.T) {

	var (
		listenHTTP  string
		debug       bool
		ctx, cancel = context.WithCancel(context.Background())
	)

	defer func() {
		// Shutdown. Cancel application context will kill all attached tasks.
		fmt.Println("Running cancel()")
		cancel()
	}()

	// command line (flag) variables and defaults
	flag.StringVar(&listenHTTP, "listenHTTP", ":9101", "ip:port to listen for http requests")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.Parse()

	go read_journal(ctx, "")      // go routine to follow/tail new journald logs and process metrics
	go prom_http(ctx, listenHTTP) // start Prometheus http endpoint
	testRoutine(t)

}
