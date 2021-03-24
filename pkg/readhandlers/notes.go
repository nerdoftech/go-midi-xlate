package readhandlers

import (
	"fmt"

	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi/reader"
)

// ErrNoteNotFound error type
type ErrNoteNotFound struct {
	note core.Note
}

// Error .
func (e ErrNoteNotFound) Error() string {
	return fmt.Sprintf("Note not found: %+v", e.note)
}

// ErrNoteWrongChannel error type
type ErrNoteWrongChannel struct {
	recvChan uint8
	expChan  uint8
}

// Error .
func (e ErrNoteWrongChannel) Error() string {
	return fmt.Sprintf("Note on wrong channel, recieved: %d, expected: %d", e.recvChan, e.expChan)
}

// NoteReader listens to NoteOn messages
type NoteReader struct {
	handlers map[uint8]core.NoteHandler
	readers  []func(reader *reader.Reader)
	midiChan uint8
}

// NewNoteReader with NoteOn/NoteOff readers
func NewNoteReader(midiCh uint8) *NoteReader {
	log.Debug().Uint8("channel", midiCh).Msg("creating note reader")
	nr := &NoteReader{
		handlers: map[uint8]core.NoteHandler{},
		midiChan: midiCh - 1, // go-midi is 0 indexed
	}
	nr.readers = []func(reader *reader.Reader){
		reader.NoteOn(func(_ *reader.Position, channel, key, velocity uint8) {
			go nr.Dispatch(core.NewNote(channel, key, velocity, core.NoteOn))
		}),
		reader.NoteOff(func(_ *reader.Position, channel, key, velocity uint8) {
			go nr.Dispatch(core.NewNote(channel, key, velocity, core.NoteOff))
		}),
	}
	return nr
}

// GetMidiReaders returns the NoteOn/NoteOff readers
func (r *NoteReader) GetMidiReaders() []func(*reader.Reader) {
	return r.readers
}

// AddHandler adds a NoteHandler for a specific note message
func (r *NoteReader) AddHandler(note uint8, handler core.NoteHandler) {
	r.handlers[note] = handler
}

// Dispatch routes the note, most of the time you will not need this directly
func (r *NoteReader) Dispatch(note core.Note) error {
	dbgMsg := log.Debug().
		Uint8("channel", note.Channel+1).
		Uint8("note", note.Key).
		Uint8("velocity", note.Velocity).
		Int64("time", note.Time).
		Str("state", string(note.State))

	dbgMsg.Msg("new note")

	if note.Channel != r.midiChan {
		dbgMsg.Uint8("handler_channel", r.midiChan+1).Msg("note on other channel")
		return ErrNoteWrongChannel{note.Channel + 1, r.midiChan + 1}
	}

	handler, found := r.handlers[note.Key]
	if !found {
		dbgMsg.Msg("note not found")
		return ErrNoteNotFound{note}
	}

	handler.HandleNote(note)
	return nil
}
