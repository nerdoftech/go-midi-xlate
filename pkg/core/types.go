package core

//go:generate mockgen -destination=mocks/mocks.go -package=mocks -source=types.go

import "gitlab.com/gomidi/midi"

type SendMidi interface {
	Send(msg midi.Message)
}

type Note struct {
	Channel  uint8
	Key      uint8
	Velocity uint8
	Time     int64
}

type NoteHandler interface {
	HandleNote(note Note)
}
