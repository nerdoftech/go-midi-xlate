package readhandlers

import (
	"testing"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/testdrv"
	"gitlab.com/gomidi/midi/writer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	drv     midi.Driver
	midiIn  midi.In
	midiOut midi.Out
	midiWr  *writer.Writer
)

func TestReadhandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Readhandlers Suite")
}

var _ = BeforeSuite(func() {
	drv = testdrv.New("test")

	ins, err := drv.Ins()
	Expect(err).Should(BeNil())
	outs, err := drv.Outs()
	Expect(err).Should(BeNil())

	midiIn, midiOut = ins[0], outs[0]

	err = midiIn.Open()
	Expect(err).Should(BeNil())
	err = midiOut.Open()
	Expect(err).Should(BeNil())

	midiWr = writer.New(midiOut)
})

var _ = AfterSuite(func() {
	midiIn.Close()
	midiOut.Close()
	drv.Close()
})
