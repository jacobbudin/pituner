package main

/*
#include <bass.h>
#cgo LDFLAGS: -lbass
*/
import "C"

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type Display struct {
	Primary   string
	Secondary string
}

type Station struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Channel int    `json:"channel"`
}

func main() {
	// Parse flags
	debug := flag.Bool("debug", false, "enable debugging mode")
	help := flag.Bool("help", false, "show this help")
	stations_file_path := flag.String("stations", "etc/stations.json", "load station file")

	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	// Do error-checking
	if uint16(C.BASS_GetVersion()>>16&0xffff) != C.BASSVERSION {
		fmt.Fprintf(os.Stderr, "BASS audio library version mismatch")
		os.Exit(1)
	}

	// Load station info
	stations, err := parseStations(*stations_file_path)

	if err != nil {
		panic(err)
	}

	if *debug {
		fmt.Printf("Loaded %d stations\n", len(stations))
	}
}

func showHelp() {
	fmt.Fprintf(os.Stderr, "Usage of pituner:\n")
	flag.PrintDefaults()
}

func parseStations(stations_file_path string) ([]Station, error) {
	// Open file
	stations_data, err := ioutil.ReadFile(stations_file_path)

	if err != nil {
		return nil, err
	}

	// Create station data
	var stations []Station

	if err := json.Unmarshal(stations_data, &stations); err != nil {
		return nil, err
	}

	return stations, nil
}
