package x32

import "fmt"

var (
	SysexoOverMIDIPrefix = []byte{0xF0, 0x00, 0x20, 0x32, 0x32}
	SysexoOverMIDISuffix = []byte{0xF7}
)

type X32Sysex struct {
	data     []byte
	overMidi bool
}

func (x X32Sysex) Raw() []byte {
	if !x.overMidi {
		return x.data
	}
	msg := SysexoOverMIDIPrefix
	msg = append(msg, x.data...)
	msg = append(msg, SysexoOverMIDISuffix...)
	return msg
}

func (x X32Sysex) String() string {
	return fmt.Sprintf("%T len: %v", x, len(x.data))
}
