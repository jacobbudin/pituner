package main

import (
	"fmt"
	"net/http"
	"strings"
)

// parsePlaylist downloads a `.pls` file and returns an array
// of URL `strings` that contain the actual playback URLs
func parsePlaylist(playlist_url string) []string {
	resp, err := http.Get(playlist_url)

	if err != nil {
		if DEBUG {
			fmt.Println("Playlist %s could not be reached", playlist_url)
		}

		return make([]string, 0)
	}

	urls := make([]string, 0)
	contents := make([]byte, 1000)

	resp.Body.Read(contents)
	resp.Body.Close()

	lines := strings.Split(string(contents[:]), "\n")

	// Add all lines starting with `File` and containing `=` to `urls`
	for _, line := range lines {
		if strings.HasPrefix(line, "File") {
			i := strings.Index(line, "=")
			if i != -1 {
				urls = append(urls, line[i+1:])
			}
		}
	}

	return urls
}
