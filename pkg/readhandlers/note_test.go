package readhandlers

import (
	"github.com/golang/mock/gomock"
	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/nerdoftech/go-midi-xlate/pkg/core/mocks"
	"gitlab.com/gomidi/midi/writer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notes", func() {
	Context("NoteReader", func() {
		var (
			mockNH *mocks.MockNoteHandler
			ctrl   *gomock.Controller
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockNH = mocks.NewMockNoteHandler(ctrl)
		})
		AfterEach(func() {
			ctrl.Finish()
			midiIn.StopListening()
		})
		It("should read note", func() {
			nr := NewNoteReader(1)
			nr.AddHandler(14, mockNH)
			err := StartReaderListen(midiIn, nr.GetMidiReader())
			Expect(err).Should(BeNil())

			var note core.Note
			mockNH.EXPECT().HandleNote(gomock.AssignableToTypeOf(note)).Do(func(n core.Note) {
				note = n
			})

			nr.Start()
			writer.NoteOn(midiWr, 14, 1)

			Eventually(func() uint8 {
				return note.Key
			}).Should(Equal(uint8(14)))

		})
		It("should not read other notes or channels", func() {
			nr := NewNoteReader(1)
			nr.AddHandler(99, mockNH)
			err := StartReaderListen(midiIn, nr.GetMidiReader())
			Expect(err).Should(BeNil())

			nr.Start()
			writer.NoteOn(midiWr, 14, 1)
			// Gomock ensures NoteHandler is not called
		})

	})
})
