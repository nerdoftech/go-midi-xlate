package main

import (
	"flag"
	"os"

	"github.com/nerdoftech/go-midi-xlate/pkg/cmd"
	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/nerdoftech/go-midi-xlate/pkg/outputs"
	"github.com/nerdoftech/go-midi-xlate/pkg/readhandlers"
	"github.com/nerdoftech/go-midi-xlate/pkg/x32"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	flgLogLvl    = flag.String("log", "info", "sets the log level")
	flgListPorts = flag.Bool("list", false, "shows MIDI ports")
	flgDelayNote = flag.Uint("note", 43, "MIDI note that will trigger beat")
	flgMidiChan  = flag.Uint("ch", 1, "MIDI channel to listen on")
	flgFxChan    = flag.Int("fxc", 1, "fx channel of delay")
	flgMidiIn    = flag.Int("in", 0, "input midi port index")
	flgMidiOut   = flag.Int("out", 0, "output midi port index")
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	flag.Parse()

	// Set log level
	lvl, err := zerolog.ParseLevel(*flgLogLvl)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse log level")
	}
	zerolog.SetGlobalLevel(lvl)

	mio, err := core.NewMidiIO()
	cmd.CheckFatalErr("", err)
	defer mio.Close()

	if *flgListPorts {
		cmd.ListPorts(mio, 0)
	}

	err = mio.OpenPorts(*flgMidiIn, *flgMidiOut)
	cmd.CheckFatalErr("", err)

	in, out := mio.GetIO()
	log.Info().Msgf("Input port: %s", in.String())
	log.Info().Msgf("Output port: %s", out.String())

	ntReader := readhandlers.NewNoteReader(uint8(*flgMidiChan))
	mOut := outputs.NewMidiOut(out)

	beatHdr := x32.NewBeatHandler(mOut, *flgFxChan, true)
	ntReader.AddHandler(uint8(*flgDelayNote), beatHdr)

	err = readhandlers.StartReaderListen(in, ntReader.GetMidiReaders()...)
	cmd.CheckFatalErr("could not get reader to listen", err)

	cmd.WaitForSignal()
}
