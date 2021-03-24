package core

//go:generate mockgen -destination=mocks/mocks.go -package=mocks -source=types.go

import (
	"time"

	"gitlab.com/gomidi/midi"
)

// SendMidi for sending messages via MIDI, OSC, etc.
type SendMidi interface {
	Send(msg midi.Message)
}

// Note with midi info
type Note struct {
	Channel  uint8
	Key      uint8
	Velocity uint8
	Time     int64
	State    NoteState
}

// NewNote creates new Note with unix milli timestamp
func NewNote(channel, key, velocity uint8, state NoteState) Note {
	return Note{
		Channel:  channel,
		Key:      key,
		Velocity: velocity,
		Time:     time.Now().UnixNano() / 1000,
		State:    state,
	}
}

// NoteState for note on/off
type NoteState string

const (
	NoteOn  NoteState = "on"
	NoteOff NoteState = "off"
)

// NoteHandler to register handlers for specific notes
type NoteHandler interface {
	HandleNote(note Note)
}
