package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/afnank19/gale/requester"
)

var UNITS = map[string]time.Duration{
	"s": time.Second,
	"m": time.Minute,
	"h": time.Hour,
}

const MARKED = "marked"

type Arguments struct {
	threads     int
	connections int
	duration    time.Duration
	url         *url.URL
	urlStr 		string
}

// The code here is not well written.
// It will be reworked, but since this is not the projects goal, i don't want to waste
// time on it
func ParseArgs() Arguments {
	allArgs := os.Args
	args := allArgs[1:]

	var a Arguments = Arguments{
		threads: -1,
		connections: 10,
		duration: 10 * time.Second,
		urlStr: "",
	}

	var argMap map[string]string = make(map[string]string)

	for _, argStr := range args {
		arg, value := splitArg(argStr)

		if arg == "--threads" || arg == "-t" {
			checkIfArgParsed(argMap, arg)
			argMap["--threads"] = MARKED
			argMap["-t"] = MARKED

			value, err := strconv.Atoi(value)
			if err != nil {
				requester.ShowUsage()
			}

			if value < 1 {
				requester.ShowUsage()
			}
			a.threads = value
			continue
		} else if arg == "--connections" || arg == "-c" {
			checkIfArgParsed(argMap, arg)
			argMap["--connections"] = MARKED
			argMap["-c"] = MARKED

			value, err := strconv.Atoi(value)
			if err != nil {
				requester.ShowUsage()
			}

			if value < 1 {
				requester.ShowUsage()
			}
			a.connections = value
			continue
		} else if arg == "--duration" || arg == "-d" {
			checkIfArgParsed(argMap, arg)
			argMap["--duration"] = MARKED
			argMap["-d"] = MARKED

			d := parseDurationValue(value)
			a.duration = d
			continue
		} else if arg == "--url" || arg == "-u" {
			checkIfArgParsed(argMap, arg)
			argMap["--url"] = MARKED
			argMap["-u"] = MARKED

			a.urlStr = value
			url := parseUrl(value)
			a.url = url
			continue
		} else {
			fmt.Printf("ERROR: Unknown flag?\n")
			requester.ShowUsage()
		}
	}

	if a.urlStr == "" {
		requester.ShowUsage()
	}

	return a
}

func checkIfArgParsed(argMap map[string]string, currArg string) {
	_, exists := argMap[currArg]
	if exists {
		fmt.Printf("ERROR: repeating arguments\n")
		requester.ShowUsage()
	}
}

func splitArg(arg string) (string, string) {
	argTokens := strings.Split(arg, "=")

	if len(argTokens) < 2 || len(argTokens) > 2 {
		requester.ShowUsage()
	}
	
	return argTokens[0], argTokens[1]
}

func parseUrl(rawUrl string) *url.URL {
	url, err := url.Parse(rawUrl)
	if err != nil {
		requester.ShowUsage()
	}
	return url
}

func parseDurationValue(val string) time.Duration {
	if len(val) < 1 {
		requester.ShowUsage()
	}

	lastChar := val[len(val)-1:]

	numVal := val[:len(val)-1]
	value, err := strconv.Atoi(numVal)
	if err != nil {
		requester.ShowUsage()
	}

	if value < 1 {
		requester.ShowUsage()
	}

	if value > 999999 {
		value = 999999
	}

	var duration time.Duration
	dMultiplier, ok := UNITS[lastChar]
	if !ok {
		requester.ShowUsage()
	}

	duration = time.Duration(value) * dMultiplier

	return duration
}
