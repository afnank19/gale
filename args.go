package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseArgs() {
	allArgs := os.Args
	fmt.Println(allArgs)

	args := allArgs[1:]
	fmt.Println("What we need -> ", args)

	for i, argStr := range args {
		fmt.Println(argStr, i)
		arg, value := SplitArg(argStr)
		if value == -1 { // means that argument does not have a value
			// handleNormalArg()
			fmt.Println("normal arg detected: ", arg)
		}

		if arg == "--worker" || arg == "-w" {
			fmt.Println("Value -> ", value)
		}
		if arg == "--concurrency" || arg == "-c" {
			fmt.Println("Value -> ", value)
		}
		if arg == "--duration" || arg == "-d" {
			fmt.Println("Value -> ", value)
		}
	}
}

func SplitArg(arg string) (string, int) {
	argTokens := strings.Split(arg, "=")

	if len(argTokens) < 2 {
		return arg, -1
	}

	if len(argTokens) > 2 {
		panic("Bad argument provided")
	}

	value, err := strconv.Atoi(argTokens[1])
	if err != nil {
		panic("Bro, value of args can only be ints bro what is wrong with you")
	}

	return argTokens[0], value
}
