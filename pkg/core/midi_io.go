package core

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	driver "gitlab.com/gomidi/rtmididrv"
)

type MidiIO struct {
	drv  midi.Driver
	ins  []midi.In
	outs []midi.Out
	in   midi.In
	out  midi.Out
}

func NewMidiIO() (*MidiIO, error) {
	drv, err := driver.New()
	if err != nil {
		return nil, errors.Wrap(err, "could not create driver")
	}

	ins, err := drv.Ins()
	if err != nil {
		return nil, errors.Wrap(err, "could not get inputs")
	}

	outs, err := drv.Outs()
	if err != nil {
		return nil, errors.Wrap(err, "could not get outputs")
	}

	mio := &MidiIO{
		drv:  drv,
		ins:  ins,
		outs: outs,
	}
	return mio, nil
}

func (m *MidiIO) Close() {
	m.drv.Close()
}

// ListPort get a list of ports by idx=string
func (m *MidiIO) GetPortList() (ins, outs map[int]string) {
	ins, outs = map[int]string{}, map[int]string{}
	for _, in := range m.ins {
		ins[in.Number()] = in.String()
	}
	for _, out := range m.outs {
		outs[out.Number()] = out.String()
	}
	return ins, outs
}

func (m *MidiIO) OpenPorts(in, out int) error {
	log.Debug().Int("in", in).Int("out", out).Msg("opening io ports")
	if in > len(m.ins)-1 {
		return errors.New(fmt.Sprintf("Error - in port out of range: %d\n\n", in))
	}

	if out > len(m.outs)-1 {
		return errors.New(fmt.Sprintf("Error - in port out of range: %d\n\n", out))
	}

	m.in = m.ins[in]
	m.out = m.outs[out]

	err := m.in.Open()
	if err != nil {
		return errors.Wrap(err, "could not open 'in' port")
	}

	err = m.out.Open()
	if err != nil {
		return errors.Wrap(err, "could not open 'out' port")
	}

	return nil
}

func (m *MidiIO) GetIO() (midi.In, midi.Out) {
	return m.in, m.out
}
