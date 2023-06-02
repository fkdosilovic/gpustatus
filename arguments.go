package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Arguments struct {
	ShowFree bool
	ShowUsed bool
}

var usage = `
Usage: %s [options]
	
Options:
    -free	Show free GPUs.
    -used	Show used GPUs.
	
Examples:
    To show all servers and GPUs run:
    %[1]s

    To show only free GPUs run:
    %[1]s -free

    To show only used GPUs run:
    %[1]s -used
`

func ParseArguments() Arguments {
	var args Arguments

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		// flag.PrintDefaults()
	}

	flag.BoolVar(&args.ShowFree, "free", false, "Show free GPUs.")
	flag.BoolVar(&args.ShowUsed, "used", false, "Show used GPUs.")
	flag.Parse()

	checkArguments(args)

	return args
}

func checkArguments(args Arguments) {
	if args.ShowFree && args.ShowUsed {
		log.Fatalf("Cannot use both -free and -used flags.")
	}
}
