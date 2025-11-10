package main

import (
	"flag"
	"fmt"
)

// flag variables
var (
	flagText    = flag.String("text", "", "The Go struct definition text")
	flagName    = flag.String("name", "", "The name of targeted Go struct")
	flagDevel   = flag.Bool("devel", false, "Enable developer mode for verbose logging")
	flagHelp    = flag.Bool("help", false, "Show help message")
	flagVersion = flag.Bool("version", false, "Show version information")
)

var Version = "unknown"

func main() {
	flag.Parse()

	routing()
}

func routing() {
	if *flagVersion {
		showVersion()
		return
	}

	if *flagHelp {
		showHelp()
		return
	}

	if isRunnale() {
		runProgram()
		return
	}

	showHelp()
}

func runProgram() {}

func showVersion() {
	fmt.Printf("Go STRUCT TO JSON CLI: %s\n", Version)
}

func showHelp() {
	flag.PrintDefaults()
}

func isRunnale() bool {
	return *flagText != "" && *flagName != ""
}
