package readhandlers

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
)

// StartReaderListen creates and starts midi reader with reader callbacks
func StartReaderListen(in midi.In, rcb ...func(*reader.Reader)) error {
	rd := reader.New(append([]func(*reader.Reader){reader.NoLogger()}, rcb...)...)
	err := rd.ListenTo(in)
	return err
}
