package readhandlers

import (
	"time"

	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi/reader"
)

const (
	NOTE_HANDLER_CHAN_SIZE = 4
)

// NoteReader listens to NoteOn messagews
type NoteReader struct {
	handlers map[uint8]core.NoteHandler
	ch       chan core.Note
	reader   func(reader *reader.Reader)
	midiChan uint8
}

// NewNoteReader new NoteReader
func NewNoteReader(midiCh uint8) NoteReader {
	log.Debug().Uint8("channel", midiCh).Msg("creating note reader")
	midiCh-- // go-midi is 0 indexed
	ch := make(chan core.Note, NOTE_HANDLER_CHAN_SIZE)
	return NoteReader{
		handlers: map[uint8]core.NoteHandler{},
		ch:       ch,
		midiChan: midiCh,
		reader: reader.NoteOn(func(_ *reader.Position, channel, key, velocity uint8) {
			if channel != midiCh {
				log.Debug().Uint8("channel", channel+1).Uint8("note", key).Msg("note on other channel")
				return
			}
			ch <- core.Note{
				Channel:  channel,
				Key:      key,
				Velocity: velocity,
				Time:     time.Now().UnixNano() / 1000,
			}
		}),
	}
}

func (r *NoteReader) GetMidiReader() func(*reader.Reader) {
	return r.reader
}

func (r *NoteReader) AddHandler(note uint8, handler core.NoteHandler) {
	r.handlers[note] = handler
}

func (r *NoteReader) Start() {
	go func() {
		for note := range r.ch {
			log.Debug().Msgf("New note: %+v", note)

			handler, found := r.handlers[note.Key]
			if !found {
				log.Debug().Int("note", int(note.Key)).Msg("note not found")
				continue
			}
			go handler.HandleNote(note)
		}
	}()
}
