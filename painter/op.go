package painter

import (
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

func BlackFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, screen.Src)
}

type BackgroundRectOp struct {
	TopLeft     image.Point
	BottomRight image.Point
}

func (br *BackgroundRectOp) Do(t screen.Texture) bool {
	t.Fill(
		image.Rect(
			br.TopLeft.X,
			br.TopLeft.Y,
			br.BottomRight.X,
			br.BottomRight.Y,
		),
		color.Black,
		screen.Src,
	)
	return false
}

type FigureOp struct {
	Center image.Point
}

func (npf *FigureOp) Do(t screen.Texture) bool {
	const length = 400
	const thickness = 100

	plusColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	horizontal := image.Rect(
		npf.Center.X-length/2,
		npf.Center.Y-thickness/2,
		npf.Center.X+length/2,
		npf.Center.Y+thickness/2,
	)
	t.Fill(horizontal, plusColor, draw.Src)

	vertical := image.Rect(
		npf.Center.X-thickness/2,
		npf.Center.Y-length/2,
		npf.Center.X+thickness/2,
		npf.Center.Y+length/2,
	)
	t.Fill(vertical, plusColor, draw.Src)

	return false
}

type MoveOp struct {
	X, Y    int
	Figures []*FigureOp
}

func (m *MoveOp) Do(_ screen.Texture) bool {
	for _, f := range m.Figures {
		f.Center.X += m.X
		f.Center.Y += m.Y
	}
	return false
}
