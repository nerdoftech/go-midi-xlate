package x32

import (
	"time"

	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/nerdoftech/go-midi-xlate/pkg/core/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/gomidi/midi"
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
		It("should be NoteHandler interface", func() {
			bh := &BeatHandler{}
			var nhI core.NoteHandler = bh
			Expect(nhI).Should(BeAssignableToTypeOf(nhI))
		})
		It("HandleNote should work", func() {
			bh := NewBeatHandler(mockSM, 4, true)

			mockSM.EXPECT().Send(gomock.AssignableToTypeOf(X32Sysex{})).
				DoAndReturn(func(msg midi.Message) {
					Expect(msg).Should(BeAssignableToTypeOf(X32Sysex{}))
					raw := msg.Raw()
					Expect(raw).Should(HaveLen(21))
					Expect(raw[0:5]).Should(Equal(SysexoOverMIDIPrefix))
					Expect(string(raw[5 : len(raw)-3])).Should(Equal("/fx/4/par/02 "))
					Expect(raw[len(raw)-1 : len(raw)]).Should(Equal(SysexoOverMIDISuffix))
					return
				})

			for i := 0; i < 5; i++ {
				note := core.NewNote(1, 10, 5, core.NoteOn)
				bh.HandleNote(note)
				time.Sleep(50 * time.Millisecond)
			}
		})
		It("HandleNote should not read note off", func() {
			bh := NewBeatHandler(mockSM, 4, false)

			for i := 0; i < 5; i++ {
				note := core.NewNote(1, 10, 5, core.NoteOff)
				bh.HandleNote(note)
				time.Sleep(50 * time.Millisecond)
			}
			// mockSM.Send should not be called
		})
	})

})
