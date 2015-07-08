package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Station struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Channel int    `json:"channel"`
}

// ParseStations opens a file path to a JSON-formatted file
// and returns an array of `Station`s
func ParseStations(stations_file_path string) ([]Station, error) {
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

	if err := validate(stations); err != nil {
		return nil, err
	}

	return stations, nil
}

// validate checks to see the array of `Station` meets logical conditions
func validate(stations []Station) error {
	// Check for duplicate `Channel`s
	channels := make(map[int]bool)

	for _, station := range stations {
		if station.Channel == 0 {
			continue
		}

		if _, exists := channels[station.Channel]; exists {
			return errors.New("duplicate channel numbers exist in station listing")
		}

		channels[station.Channel] = true
	}

	return nil
}
