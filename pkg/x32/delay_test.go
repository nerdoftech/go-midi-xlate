package x32

import (
	"time"

	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/gomidi/midi"

	"github.com/golang/mock/gomock"
	"github.com/nerdoftech/go-midi-xlate/pkg/core/mocks"
)

var _ = Describe("x32 Delay", func() {
	Context("BeatHandler", func() {
		var (
			mockSM *mocks.MockSendMidi
			ctrl   *gomock.Controller
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockSM = mocks.NewMockSendMidi(ctrl)
		})
		AfterEach(func() {
			ctrl.Finish()
		})
		It("BeatHandler should be NoteHandler interface", func() {
			bh := &BeatHandler{}
			var nhI core.NoteHandler = bh
			Expect(nhI).Should(BeAssignableToTypeOf(nhI))
		})
		It("HandleNote should work", func() {
			bh := NewBeatHandler(mockSM, 4, true)

			note := core.Note{
				Channel:  1,
				Key:      10,
				Velocity: 5,
				Time:     0,
			}

			mockSM.EXPECT().Send(gomock.AssignableToTypeOf(X32Sysex{})).
				DoAndReturn(func(msg midi.Message) {
					Expect(msg).Should(BeAssignableToTypeOf(X32Sysex{}))
					raw := msg.Raw()
					Expect(raw).Should(HaveLen(22))
					return
				})

			for i := 0; i < 5; i++ {
				note.Time = time.Now().UnixNano() / 1000
				bh.HandleNote(note)
				time.Sleep(250 * time.Millisecond)
			}

		})
	})

})
