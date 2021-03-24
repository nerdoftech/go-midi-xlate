package readhandlers

import (
	"github.com/nerdoftech/go-midi-xlate/pkg/core"
	"github.com/nerdoftech/go-midi-xlate/pkg/core/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/gomidi/midi/writer"
)

var _ = Describe("NoteReader", func() {
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
	Context("Dispatch", func() {
		It("should work", func() {
			nr := NewNoteReader(1)
			nr.AddHandler(14, mockNH)

			var note core.Note
			mockNH.EXPECT().HandleNote(gomock.AssignableToTypeOf(note)).Do(func(n core.Note) {
				note = n
			})

			err := nr.Dispatch(core.NewNote(0, 14, 1, core.NoteOn))
			Expect(err).Should(BeNil())
		})
		It("should not read unhandled notes", func() {
			nr := NewNoteReader(1)
			nr.AddHandler(99, mockNH)

			err := nr.Dispatch(core.NewNote(0, 14, 1, core.NoteOn))
			Expect(err).Should(BeAssignableToTypeOf(ErrNoteNotFound{}))
		})
		It("should not read notes from unread channels", func() {
			nr := NewNoteReader(1)
			nr.AddHandler(14, mockNH)

			err := nr.Dispatch(core.NewNote(9, 14, 1, core.NoteOn))
			Expect(err).Should(BeAssignableToTypeOf(ErrNoteWrongChannel{}))
		})
	})
	Context("Test with midi readers.Reader", func() {
		It("should read note on/off", func() {
			nr := NewNoteReader(1)
			nr.AddHandler(14, mockNH)
			err := StartReaderListen(midiIn, nr.GetMidiReaders()...)
			Expect(err).Should(BeNil())
			defer midiIn.StopListening()

			notes := make([]core.Note, 2)
			time := 0
			mockNH.EXPECT().HandleNote(gomock.AssignableToTypeOf(core.Note{})).Do(func(n core.Note) {
				notes[time] = n
				time++
			}).Times(2)

			writer.NoteOn(midiWr, 14, 1)
			Eventually(func() uint8 {
				return notes[0].Key
			}).Should(Equal(uint8(14)))

			writer.NoteOff(midiWr, 14)
			Eventually(func() uint8 {
				return notes[1].Key
			}).Should(Equal(uint8(14)))
		})
	})
})
