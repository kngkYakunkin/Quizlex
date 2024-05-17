package widgets

import (
	"fyne.io/fyne/v2/widget"
)

type subtitle struct {
	widget.Label
	text string
}

func NewSubtitle(text string) *subtitle {
	subtitle := &subtitle{text: text}
	subtitle.Text = text

	subtitle.ExtendBaseWidget(subtitle)

	return subtitle

}
