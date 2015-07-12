package main

/*
#include <bass.h>
#cgo LDFLAGS: -lbass
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
)

var DEBUG bool = false

func main() {
	// Parse flags
	flag.BoolVar(&DEBUG, "debug", false, "enable debugging mode")
	help := flag.Bool("help", false, "show this help")
	stations_file_path := flag.String("stations", "etc/stations.json", "load station file")

	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	// Do checks, load stations, begin playback
	checkSuperuser()
	initPlayback()

	stations := loadStations(*stations_file_path)

	_ = Tuner{
		Stations: stations,
		Display:  Display{},
	}
}

// loadStations loads the stations files
func loadStations(stations_file_path string) []Station {
	stations, err := parseStations(stations_file_path)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if DEBUG {
		fmt.Println("Loaded %d stations", len(stations))
	}

	return stations
}

// initPlayback loads BASS audio library
func initPlayback() {
	if uint16(C.BASS_GetVersion()>>16&0xffff) != C.BASSVERSION {
		fmt.Fprintln(os.Stderr, "BASS audio library version mismatch")
		os.Exit(1)
	}

	if C.BASS_Init(-1, 44100, 0, nil, nil) == C.TRUE {
		fmt.Fprintln(os.Stderr, "Can't initialize audio device")
		os.Exit(1)
	}

	C.BASS_SetVolume(1)
	C.BASS_SetConfig(C.BASS_CONFIG_NET_PLAYLIST, 1)
	C.BASS_SetConfig(C.BASS_CONFIG_NET_PREBUF, 0)
}

// checkSuperuser checks the user has root privileges, that are
// required on the Raspberry Pi to access GPIO pins
func checkSuperuser() {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "Must call pituner with superuser privileges")
		os.Exit(1)
	}
}

// showHelp prints help information to stderr
func showHelp() {
	fmt.Fprintf(os.Stderr, "Usage of pituner:\n")
	flag.PrintDefaults()
}
