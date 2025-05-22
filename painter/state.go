package painter

import (
	"golang.org/x/exp/shiny/screen"
	"image"
)

type FillFunc func(t screen.Texture)
type PainterState struct {
	bgOp     Operation
	bgRect   *BackgroundRectOp
	figures  []*FigureOp
	moves    []Operation
	isUpdate bool
}

func (p *PainterState) SetBackground(fillFunc FillFunc) {
	p.bgOp = OperationFunc(fillFunc)
}

func (p *PainterState) SetBackgroundRect(topLeft image.Point, bottomRight image.Point) {
	op := &BackgroundRectOp{
		TopLeft:     topLeft,
		BottomRight: bottomRight,
	}
	p.bgRect = op
}

func (p *PainterState) LoadOperations() []Operation {
	var opList []Operation

	if p.bgOp != nil {
		opList = append(opList, p.bgOp)
	}

	if p.bgRect != nil {
		opList = append(opList, p.bgRect)
	}

	if len(p.moves) != 0 {
		opList = append(opList, p.moves...)
	}

	if len(p.figures) != 0 {
		for _, tFigure := range p.figures {
			opList = append(opList, tFigure)
		}
	}
	if p.isUpdate {
		opList = append(opList, updateOp{})
	}
	return opList
}

func (p *PainterState) Update() {
	p.isUpdate = true
}

func (p *PainterState) AddFigure(center image.Point) {
	op := &FigureOp{
		Center: center,
	}
	p.figures = append(p.figures, op)
}

func (p *PainterState) AddMove(x, y int) {
	op := &MoveOp{
		X:       x,
		Y:       y,
		Figures: p.figures,
	}
	p.moves = append(p.moves, op)
}

func (p *PainterState) Reset() {
	p.bgOp = nil
	p.bgRect = nil
	p.figures = nil
	p.moves = nil
	p.isUpdate = true
	p.SetBackground(BlackFill)
}
