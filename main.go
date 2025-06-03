package main

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"sync"

	"github.com/afnank19/gale/requester"
)

func setMaxProcs(procs int) {
	runtime.GOMAXPROCS(procs)
}

func main() {
	// For a 4 core 8 thread cpu, the default by Go will be 8, but best perf comes from 4 so we default to that
	setMaxProcs(runtime.NumCPU() / 2)
	fmt.Println("GALE! Hammer your server.")

	argsState := ParseArgs()

	if argsState.threads != -1 {
		setMaxProcs(argsState.threads)
	}

	host, target := resolveUrlPort(argsState.url)
	// fmt.Println("Host: "+host+" | Target: ", target+" | Path: ", argsState.url.Path)
	var rd requester.ReqData = requester.ReqData{
		Host: host,
		Scheme: argsState.url.Scheme,
		Target: target,
		Path: argsState.url.Path,
		Duration: argsState.duration,
	}

	var wg sync.WaitGroup
	var r requester.Result
	r.StatusCodes = make(map[int]int)
	r.TestDuration = argsState.duration

	// range is connections
	for range argsState.connections {
		wg.Add(1)
		go r.MakeRequest(&wg, &rd)
	}

	wg.Wait()

	report := requester.GenerateReport(&r)
	requester.DisplayReport(report)
}

func resolveUrlPort(url *url.URL) (host string, target string) {
	if url.Scheme == "http" {
		// handleHTTPUrl()
		host = url.Host

		// host already has a port specified
		if strings.Contains(host, ":") {
			target = host
		} else {
			target = host + ":80"
		}
	}

	if url.Scheme == "https" {
		// hanldeHTTPSUrl()
		host = url.Host
		target = host + ":443"
	}

	return
}