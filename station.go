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

	if err := validate(&stations); err != nil {
		return nil, err
	}

	generateChannels(&stations)

	return stations, nil
}

// validate checks to see the array of `Station` meets logical conditions
func validate(stations *[]Station) error {
	channels := make(map[int]bool)

	for _, station := range *stations {
		if station.Channel == 0 {
			continue
		}

		// Check for negative `Channels`
		if station.Channel < 0 {
			return errors.New("negative channel numbers exist in station listing")
		}

		// Check for duplicate `Channel`s
		if _, exists := channels[station.Channel]; exists {
			return errors.New("duplicate channel numbers exist in station listing")
		}

		channels[station.Channel] = true
	}

	return nil
}

// generateChannels creates channel numbers for those stations that lack them
func generateChannels(stations *[]Station) {
	i := 0
	channels := make(map[int]bool)

	for _, station := range *stations {
		if station.Channel != 0 {
			channels[station.Channel] = true
		}
	}

	for j, station := range *stations {
		if station.Channel == 0 {
			for {
				i++

				if _, exists := channels[i]; exists == false {
					(*stations)[j].Channel = i
					channels[i] = true
					break
				}
			}
		}
	}
}
