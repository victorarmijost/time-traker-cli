package repl

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type CloseTerm func()

func setupTerm() (*term.Terminal, CloseTerm) {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		panic("stdion/stout should be terminal")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	return term.NewTerminal(screen, ""), func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

func (h *Handler) setTermPrompt(p string) {
	if h.terminal == nil {
		panic("no term started")
	}

	t := h.terminal

	t.SetPrompt(h.colorString(p, YELLOW))
}

func (h *Handler) readTermLine() (string, error) {
	if h.terminal == nil {
		panic("no term started")
	}

	return h.terminal.ReadLine()
}

type TermColor uint16

const (
	BLACK TermColor = iota
	BLUE
	CYAN
	GREEN
	RED
	MAGENTA
	YELLOW
)

func (h *Handler) colorString(msg string, color TermColor) string {
	t := h.terminal

	c := ""
	switch color {
	case BLACK:
		c = string(t.Escape.Black)
	case BLUE:
		c = string(t.Escape.Blue)
	case CYAN:
		c = string(t.Escape.Cyan)
	case GREEN:
		c = string(t.Escape.Green)
	case MAGENTA:
		c = string(t.Escape.Magenta)
	case RED:
		c = string(t.Escape.Red)
	case YELLOW:
		c = string(t.Escape.Yellow)
	default:
		c = string(t.Escape.Reset)
	}

	return fmt.Sprintf("%s%s%s", c, msg, string(t.Escape.Reset))
}
