package x32

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

const (
	TOTAL_beats = 3
	delayOSCStr = "/fx/%d/par/02 %d"
)

var (
	beatTimeout = time.Second * 5
	cacheClean  = time.Second * 30
)

type beats []int64

func (b beats) Len() int           { return len(b) }
func (b beats) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b beats) Less(i, j int) bool { return b[i] < b[j] }

type beatCache struct {
	ch *cache.Cache
}

func newBeatCache() *beatCache {
	return &beatCache{
		ch: cache.New(beatTimeout, cacheClean),
	}
}

func (b *beatCache) addTick(nt uint8, ts int64) int {
	key := strconv.FormatInt(int64(nt), 10) + "-" + strconv.FormatInt(ts, 10)
	b.ch.SetDefault(key, ts)
	return b.ch.ItemCount()
}

func (b *beatCache) getTicks() []int64 {
	tk := make(beats, b.ch.ItemCount())
	cnt := 0
	for _, t := range b.ch.Items() {
		tk[cnt] = t.Object.(int64)
		cnt++
	}
	sort.Sort(tk)
	return tk
}

func (b *beatCache) flush() {
	b.ch.Flush()
}

// BeatHandler collects noteOn until a tempo can be calculated
type BeatHandler struct {
	bc       *beatCache
	fxChan   int
	cb       core.SendMidi
	overMidi bool
}

// NewBeatHandler returns BeatHandler
func NewBeatHandler(cb core.SendMidi, fxChan int, overMidi bool) *BeatHandler {
	return &BeatHandler{
		bc:       newBeatCache(),
		fxChan:   fxChan,
		cb:       cb,
		overMidi: overMidi,
	}
}

// HandleNote processes incoming notes to get tempo
func (b *BeatHandler) HandleNote(nt core.Note) {
	// Only listen for NoteOn
	if nt.State != core.NoteOn {
		return
	}
	cnt := b.bc.addTick(nt.Key, nt.Time/1000)
	log.Debug().Interface("note", nt).Int("count", cnt).Msg("handling note")
	if cnt > TOTAL_beats {
		log.Debug().Int("count", cnt).Msg("got notes to calculate average")
		b.processBeats()
	}
}

func (b *BeatHandler) processBeats() {
	tm := avg(b.bc.getTicks())
	log.Debug().Int64("delay", tm).Msg("calculated delay average")
	b.bc.flush()

	oscCmd := fmt.Sprintf(delayOSCStr, b.fxChan, tm)
	msg := NewX32Sysex([]byte(oscCmd), b.overMidi)
	log.Debug().
		Str("osc_cmd", oscCmd).
		Msg("sending message")

	b.cb.Send(msg)
}

func avg(tks []int64) int64 {
	var sum int64 = 0
	for i, tk := range tks {
		if i > 0 {
			sum += tk - tks[i-1]
		}
	}
	return sum / int64(len(tks)-1)
}
