package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var UNITS = map[string]time.Duration{
	"s": time.Second,
	"m": time.Minute,
	"h": time.Hour,
}

type Arguments struct {
	threads     int
	connections int
	duration    time.Duration
	url         *url.URL
}

func ParseArgs() Arguments {
	allArgs := os.Args
	args := allArgs[1:]

	var a Arguments = Arguments{
		threads: -1,
	}

	for _, argStr := range args {
		arg, value := splitArg(argStr)

		if arg == "--threads" || arg == "-t" {
			value, err := strconv.Atoi(value)
			if err != nil {
				panic("Arg value should be an integer")
			}

			if value < 1 {
				panic("SHOW USAGE")
			}
			a.threads = value
			continue
		}
		if arg == "--connections" || arg == "-c" {
			value, err := strconv.Atoi(value)
			if err != nil {
				panic("Arg value should be an integer")
			}

			if value < 1 {
				panic("SHOW USAGE")
			}
			a.connections = value
			continue
		}
		if arg == "--duration" || arg == "-d" {
			fmt.Println("Value -> ", value)
			d := parseDurationValue(value)
			a.duration = d
			continue
		}
		if arg == "--url" || arg == "-u" {
			url := parseUrl(value)
			a.url = url
			continue
		}
	}

	return a
}

func splitArg(arg string) (string, string) {
	argTokens := strings.Split(arg, "=")

	if len(argTokens) < 2 {
		return arg, ""
	}

	if len(argTokens) > 2 {
		panic("Bad argument provided")
	}

	return argTokens[0], argTokens[1]
}

func parseUrl(rawUrl string) *url.URL {
	url, err := url.Parse(rawUrl)
	if err != nil {
		panic("Bad URL")
	}
	return url
}

func parseDurationValue(val string) time.Duration {
	if len(val) < 1 {
		panic("No value provided for duration")
	}

	lastChar := val[len(val)-1:]

	numVal := val[:len(val)-1]
	value, err := strconv.Atoi(numVal)
	if err != nil {
		log.Fatalln("ERROR: Arg value should be an integer")
	}

	if value < 1 {
		log.Fatalln("SHOW USAGE")
	}

	if value > 999999 {
		value = 999999
	}

	var duration time.Duration
	dMultiplier, ok := UNITS[lastChar]
	if !ok {
		log.Fatalln("SHOW USAGE, LAST CHAR")
	}

	duration = time.Duration(value) * dMultiplier

	fmt.Println("Final Duration: ", duration)
	return duration
}
