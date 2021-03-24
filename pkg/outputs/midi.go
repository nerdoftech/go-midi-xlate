package outputs

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
)

// MidiOutput is SendMidi that contains midi writer
type MidiOutput struct {
	wr midi.Writer
}

// NewMidiOut creates MidiOutput with write that has ConsolidateNotes disabled
func NewMidiOut(out midi.Out) *MidiOutput {
	wr := writer.New(out)
	wr.ConsolidateNotes(false)
	return &MidiOutput{
		wr: wr,
	}
}

// Send midi message
func (m *MidiOutput) Send(msg midi.Message) {
	err := m.wr.Write(msg)
	if err != nil {
		log.Err(err).Msg("could not write message to output")
	}
}
