package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/rs/zerolog/log"
)

func CheckFatalErr(desc string, err error) {
	if err != nil {
		log.Fatal().Err(err).Msg(desc)
	}
}

func WaitForSignal() {
	log.Debug().Msgf("server started, waiting for signal to shutdown")
	sdCh := make(chan os.Signal, 1)
	signal.Notify(sdCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sdCh
	log.Info().Interface("signal", sig).Msg("got signal, shutting down")
}

func ListPorts(mio *core.MidiIO, code int) {
	ins, outs := mio.GetPortList()
	fmt.Println("In ports:")
	printPortMap(ins)
	fmt.Println("Out ports:")
	printPortMap(outs)
	os.Exit(code)
}

func printPortMap(pm map[int]string) {
	for i, p := range pm {
		fmt.Printf("  %d - %s\n", i, p)
	}
}
