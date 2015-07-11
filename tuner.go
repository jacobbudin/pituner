package main

/*
#include <bass.h>
#include <stdlib.h>
#cgo LDFLAGS: -lbass
*/
import "C"

import (
	"unsafe"
)

type Tuner struct {
	Stations       []Station
	CurrentStation Station
	Display        Display
	stream         C.HSTREAM
	stream_url     *C.char
}

// play begins playback
func (t *Tuner) play(station *Station) {
	t.stream_url = C.CString((*station).Url)
	t.stream = C.BASS_StreamCreateURL(t.stream_url, 0, C.BASS_STREAM_BLOCK|C.BASS_STREAM_STATUS|C.BASS_STREAM_AUTOFREE, nil, nil)
}

// stop stops playback
func (t *Tuner) stop() {
	C.BASS_ChannelStop((C.DWORD)(t.stream))
	C.BASS_StreamFree(t.stream)
	C.free((unsafe.Pointer)(t.stream_url))
}
