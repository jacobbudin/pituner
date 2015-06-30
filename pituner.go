package main

import (
	"flag"
	"fmt"
	"os"
)

type display struct {
	primary   string
	secondary string
}

type station struct {
	name    string
	url     string
	channel int
}

func main() {
	// Parse flags
	debug := flag.Bool("debug", false, "enable debugging mode")
	help := flag.Bool("help", false, "show this help")
	stations := flag.String("stations", "~/.pituner/stations.json", "load station file")

	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	fmt.Printf("%t %t %s\n", *debug, *help, *stations)
}

func showHelp() {
	fmt.Fprintf(os.Stderr, "Usage of pituner:\n")
	flag.PrintDefaults()
}
