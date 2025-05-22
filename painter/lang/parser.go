package lang

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	ps painter.PainterState
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		commandLine := scanner.Text()
		if commandLine == "" {
			continue
		}
		err := p.parseCommand(commandLine)
		if err != nil {
			return nil, err
		}
	}

	return p.ps.LoadOperations(), nil
}

var mapCommand = map[string]func(*Parser, []string) error{
	"white":  setBgWhite,
	"green":  setBgGreen,
	"update": update,
	"bgrect": setBgRect,
	"figure": addFigure,
	"move":   addMove,
	"reset":  reset,
}

func (p *Parser) parseCommand(line string) error {
	line = strings.TrimSpace(line)
	tokens := strings.Split(line, " ")
	command := tokens[0]
	handler, ok := mapCommand[command]
	if !ok {
		return errors.New("unknown command: " + command)
	}
	return handler(p, tokens[1:])
}

func setBgWhite(p *Parser, args []string) error {
	if len(args) != 0 {
		return errors.New("white: invalid arguments")
	}
	p.ps.SetBackground(painter.WhiteFill)
	return nil
}

func setBgGreen(p *Parser, args []string) error {
	if len(args) != 0 {
		return errors.New("green: invalid arguments")
	}
	p.ps.SetBackground(painter.GreenFill)
	return nil
}

func update(p *Parser, args []string) error {
	if len(args) != 0 {
		return errors.New("update: invalid arguments")
	}
	p.ps.Update()
	return nil
}

func setBgRect(p *Parser, args []string) error {
	if len(args) != 4 {
		return errors.New("bgrect: invalid arguments")
	}
	x1, err := parseNormalizedCoord(args[0])
	if err != nil {
		return err
	}
	y1, err := parseNormalizedCoord(args[1])
	if err != nil {
		return err
	}
	x2, err := parseNormalizedCoord(args[2])
	if err != nil {
		return err
	}
	y2, err := parseNormalizedCoord(args[3])
	if err != nil {
		return err
	}
	p.ps.SetBackgroundRect(image.Point{X: x1, Y: y1}, image.Point{X: x2, Y: y2})
	return nil
}

func parseNormalizedCoord(s string) (int, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert coordinates %q: %w", s, err)
	}
	const scale = 800
	return int(f * scale), nil
}

func addFigure(p *Parser, args []string) error {
	if len(args) != 2 {
		return errors.New("figure: invalid arguments")
	}
	x, err := parseNormalizedCoord(args[0])
	if err != nil {
		return err
	}
	y, err := parseNormalizedCoord(args[1])
	if err != nil {
		return err
	}
	p.ps.AddFigure(image.Point{X: x, Y: y})
	return nil
}

func addMove(p *Parser, args []string) error {
	if len(args) != 2 {
		return errors.New("move: invalid arguments")
	}
	x, err := parseNormalizedCoord(args[0])
	if err != nil {
		return err
	}
	y, err := parseNormalizedCoord(args[1])
	if err != nil {
		return err
	}
	p.ps.AddMove(x, y)
	return nil
}

func reset(p *Parser, args []string) error {
	if len(args) != 0 {
		return errors.New("reset: invalid arguments")
	}
	p.ps.Reset()
	return nil
}
