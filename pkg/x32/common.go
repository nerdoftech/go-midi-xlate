package x32

import "fmt"

var (
	SysexoOverMIDIPrefix = []byte{0xF0, 0x00, 0x20, 0x32, 0x32}
	SysexoOverMIDISuffix = []byte{0xF7}
)

// X32Sysex creates a SysEX message for X32 consoles
type X32Sysex struct {
	data     []byte
	overMidi bool
}

// NewX32Sysex overMidi will add header for SysEx over MIDI
func NewX32Sysex(data []byte, overMidi bool) X32Sysex {
	return X32Sysex{
		data:     data,
		overMidi: overMidi,
	}
}

// Raw bytes of Sysex message w or w/o MIDI header
func (x X32Sysex) Raw() []byte {
	if x.overMidi {
		msg := SysexoOverMIDIPrefix
		msg = append(msg, x.data...)
		msg = append(msg, SysexoOverMIDISuffix...)
		return msg
	}
	return x.data
}

func (x X32Sysex) String() string {
	return fmt.Sprintf("%T len: %v", x, len(x.data))
}
