package outputs

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
)

type MidiOutput struct {
	wr midi.Writer
}

func NewMidiOut(out midi.Out) *MidiOutput {
	return &MidiOutput{
		wr: writer.New(out),
	}
}

func (m *MidiOutput) Send(msg midi.Message) {
	err := m.wr.Write(msg)
	if err != nil {
		log.Err(err).Msg("could not write message to output")
	}
}
