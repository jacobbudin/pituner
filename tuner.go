package main

/*
#include <bass.h>
#include <stdlib.h>
#cgo LDFLAGS: -lbass
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

type Tuner struct {
	Stations       []Station
	CurrentStation Station
	Display        Display
	stream         C.HSTREAM
	stream_url     *C.char
}

// playIndex begins playback of a `Station`, as determined by its index in `Stations`
func (t *Tuner) playIndex(index int) {
	if index < 0 || index >= len(t.Stations) {
		if DEBUG {
			fmt.Printf("Could not find station with index %d\n", index)
		}

		return
	}

	t.playStation(&t.Stations[index])
}

// playChannel begins playback of a `Station`, as determined by its channel number
func (t *Tuner) playChannel(channel int) {
	for _, station := range t.Stations {
		if station.Channel == channel {
			t.playStation(&station)
			return
		}
	}

	if DEBUG {
		fmt.Printf("Could not find channel #%d\n", channel)
	}
}

// playStation begins playback of a `Station`
func (t *Tuner) playStation(station *Station) {
	t.stream_url = C.CString((*station).Url)
	t.stream = C.BASS_StreamCreateURL(t.stream_url, 0, C.BASS_STREAM_BLOCK|C.BASS_STREAM_STATUS|C.BASS_STREAM_AUTOFREE, nil, nil)
	for {
		progress := (C.BASS_StreamGetFilePosition(t.stream, C.BASS_FILEPOS_BUFFER) * 100) / C.BASS_StreamGetFilePosition(t.stream, C.BASS_FILEPOS_END)

		if progress >= 75 {
			C.BASS_ChannelPlay((C.DWORD)(t.stream), C.FALSE)
			return
		}
	}
}

// stop stops playback
func (t *Tuner) stop() {
	C.BASS_ChannelStop((C.DWORD)(t.stream))
	C.BASS_StreamFree(t.stream)
	C.free((unsafe.Pointer)(t.stream_url))
}

// info gets playback info, if playing
func (t *Tuner) info() string {
	if t.stream == 0 {
		return ""
	}

	info := C.GoString(C.BASS_ChannelGetTags((C.DWORD)(t.stream), C.BASS_TAG_META))
	trigger := "StreamTitle='"

	if i := strings.Index(info, trigger); i != -1 {
		b := i + len(trigger)
		e := len(info)

		// find next unescaped single quote
		for j, r := range info {
			if j > b && r == '\'' && info[j-1] != '\\' {
				e = j
				break
			}
		}

		return info[b:e]
	}

	return info
}
