package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AnswerButton struct {
	widget.BaseWidget
	fillColor   color.Color
	hoverColor  color.Color
	isHovered   bool
	Text        *canvas.Text
	Clicked     func()
	RightAnswer bool
	WrongAnswer bool
}

func NewAnswerButton(fill color.Color, hover color.Color, text string) *AnswerButton {
	r := &AnswerButton{
		fillColor:  fill,
		hoverColor: hover,
		Text:       canvas.NewText(text, color.White),
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *AnswerButton) CreateRenderer() fyne.WidgetRenderer {

	rect := canvas.NewRectangle(r.fillColor)
	rect.CornerRadius = 4

	r.Text = canvas.NewText(r.Text.Text, color.White)

	return &hoverableRectangleRenderer{rect: rect, hoverableRectangle: r, text: r.Text, rightAnswer: &r.RightAnswer, wrongAnswer: &r.WrongAnswer}
}

type hoverableRectangleRenderer struct {
	rect               *canvas.Rectangle
	text               *canvas.Text
	hoverableRectangle *AnswerButton
	rightAnswer        *bool
	wrongAnswer        *bool
}

func (r *hoverableRectangleRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
	r.text.TextSize = 15
	r.text.TextStyle.Bold = true
	r.text.Color = theme.ForegroundColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos((size.Width-r.text.MinSize().Width)/2, (size.Height-r.text.MinSize().Height)/2))
}

func (r *hoverableRectangleRenderer) MinSize() fyne.Size {
	return fyne.NewSize(10, 40)
}

func (r *hoverableRectangleRenderer) Refresh() {
	if r.hoverableRectangle.isHovered {
		r.rect.FillColor = r.hoverableRectangle.hoverColor
	} else {
		r.rect.FillColor = r.hoverableRectangle.fillColor
	}

	if *r.rightAnswer {
		r.rect.StrokeColor = color.RGBA{0, 128, 0, 255}
		r.rect.StrokeWidth = 1
		r.text.Color = color.RGBA{0, 128, 0, 255}
	} else if *r.wrongAnswer {
		r.rect.StrokeColor = color.RGBA{255, 0, 0, 255}
		r.rect.StrokeWidth = 1
		r.text.Color = color.RGBA{255, 0, 0, 255}
	}

	canvas.Refresh(r.rect)
	canvas.Refresh(r.text)
}

func (r *hoverableRectangleRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *hoverableRectangleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect, r.text}
}

func (r *hoverableRectangleRenderer) Destroy() {}

func (r *AnswerButton) MouseIn(*desktop.MouseEvent) {
	r.isHovered = true
	r.Refresh()
}

func (r *AnswerButton) MouseOut() {
	r.isHovered = false
	r.Refresh()
}

func (r *AnswerButton) MouseMoved(*desktop.MouseEvent) {}

func (r *AnswerButton) Tapped(e *fyne.PointEvent) {
	r.Clicked()
}
