package main

import (
	"fmt"
	"net/http"
	"strings"
)

// parsePlaylist downloads a `.pls` file and returns an array
// of URL `strings` that contain the actual playback URLs
func parsePlaylist(playlist_url string) []string {
	urls := make([]string, 0)
	resp, err := http.Get(playlist_url)

	if err != nil {
		if DEBUG {
			fmt.Println("Playlist %s could not be reached: %s", playlist_url, err.Error())
		}

		return urls
	}

	if resp.ContentLength == -1 {
		if DEBUG {
			fmt.Println("Playlist %s response was malformed", playlist_url)
		}

		return urls
	}

	contents := make([]byte, resp.ContentLength)

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

	if len(urls) == 0 {
		if DEBUG {
			fmt.Println("Playlist %s contained no playback URLs", playlist_url)
		}
	}

	return urls
}
