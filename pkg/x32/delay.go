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
	TOTAL_BEATS = 3
	CHAN_OUT    = 1
	delayOSCStr = "/fx/%d/par/02 %d"
)

var (
	beatTimeout = time.Second * 5
	cacheClean  = time.Second * 30
)

type Beats []int64

func (b Beats) Len() int           { return len(b) }
func (b Beats) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b Beats) Less(i, j int) bool { return b[i] < b[j] }

type BeatCache struct {
	ch *cache.Cache
}

func NewBeatCache() *BeatCache {
	return &BeatCache{
		ch: cache.New(beatTimeout, cacheClean),
	}
}

func (b *BeatCache) AddTick(nt uint8, ts int64) int {
	key := strconv.FormatInt(int64(nt), 10) + "-" + strconv.FormatInt(ts, 10)
	b.ch.SetDefault(key, ts)
	return b.ch.ItemCount()
}

func (b *BeatCache) GetTicks() []int64 {
	tk := make(Beats, b.ch.ItemCount())
	cnt := 0
	for _, t := range b.ch.Items() {
		tk[cnt] = t.Object.(int64)
		cnt++
	}
	sort.Sort(tk)
	return tk
}

func (b *BeatCache) Flush() {
	b.ch.Flush()
}

type BeatHandler struct {
	bc       *BeatCache
	fxChan   int
	cb       core.SendMidi
	overMidi bool
}

func NewBeatHandler(cb core.SendMidi, fxChan int, overMidi bool) *BeatHandler {
	return &BeatHandler{
		bc:       NewBeatCache(),
		fxChan:   fxChan,
		cb:       cb,
		overMidi: overMidi,
	}
}

func (b *BeatHandler) HandleNote(nt core.Note) {
	cnt := b.bc.AddTick(nt.Key, nt.Time/1000)
	log.Debug().Interface("note", nt).Int("count", cnt).Msg("handle note beat")
	if cnt > TOTAL_BEATS {
		log.Debug().Msg("got beats, averaging")
		b.processBeats()
	}
}

func (b *BeatHandler) processBeats() {
	tm := avg(b.bc.GetTicks())
	b.bc.Flush()
	log.Debug().Int64("avg", tm).Msgf("calculated avg")
	msg := X32Sysex{
		data:     []byte(fmt.Sprintf(delayOSCStr, b.fxChan, tm)),
		overMidi: b.overMidi,
	}
	log.Debug().Str("msg", msg.String()).Msg("sending message")
	b.cb.Send(msg)
}

func avg(tks []int64) int64 {
	var sum int64 = 0
	for i, tk := range tks {
		if i > 0 {
			sum += tk - tks[i-1]
			log.Debug().Msgf("tick diff: %d", tk-tks[i-1])
		}
	}
	log.Debug().Msgf("avg sum: %d, cnt: %d", sum, len(tks)-1)

	return sum / int64(len(tks)-1)
}
