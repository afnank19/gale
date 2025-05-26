package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseArgs() {
	allArgs := os.Args
	fmt.Println(allArgs)

	args := allArgs[1:]
	fmt.Println("What we need -> ", args)

	for _, argStr := range args {
		fmt.Println(argStr)
		arg, value := splitArg(argStr)
		if value == "" { // means that argument does not have a value
			// handleNormalArg()
			fmt.Println("normal arg detected: ", arg)
		}

		if arg == "--worker" || arg == "-w" {
			value, err := strconv.Atoi(value)
			if err != nil {
				panic("Arg value should be an integer")
			}
			fmt.Println("Value -> ", value)
			continue
		}
		if arg == "--concurrency" || arg == "-c" {
			value, err := strconv.Atoi(value)
			if err != nil {
				panic("Arg value should be an integer")
			}
			fmt.Println("Value -> ", value)
			continue
		}
		if arg == "--duration" || arg == "-d" {
			fmt.Println("Value -> ", value)
			parseDurationValue(value)
			continue
		}
	}
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

	if value > 999999 {
		value = 999999
	}

	var duration time.Duration
	if lastChar == "s" {
		duration += time.Duration(value) * time.Second
	}

	if lastChar == "m" {
		duration += time.Duration(value) * time.Minute
	}

	if lastChar == "h" {
		duration += time.Duration(value) * time.Hour
	}

	fmt.Println("Final Duration: ", duration)
	return duration
}
