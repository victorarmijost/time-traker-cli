package myterm

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"varmijo/time-tracker/tt/infrastructure/cmd/repl"

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

type MyTerm struct {
	terminal   *term.Terminal
	promptLock bool
}

func NewTerm() (*MyTerm, CloseTerm) {
	term, close := setupTerm()

	return &MyTerm{term, false}, close
}

func (h *MyTerm) setTermPrompt(prompt string) error {
	if h.terminal == nil {
		return fmt.Errorf("no terminal started")
	}

	h.terminal.SetPrompt(h.colorString(prompt, YELLOW))

	return h.updatePrompt()
}

func (h *MyTerm) updatePrompt() error {
	if h.terminal == nil {
		return fmt.Errorf("no terminal started")
	}
	_, err := h.terminal.Write([]byte{})

	return err
}

func (h *MyTerm) readTermLine() (string, error) {
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

func (h *MyTerm) colorString(msg string, color TermColor) string {
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

func (h *MyTerm) PrintMsg(msg string) {
	fmt.Fprintln(h.terminal, msg)
	h.Br()
}

func (h *MyTerm) PrintError(err error) {
	h.PrintErrorMsg(err.Error())
}

func (h *MyTerm) PrintErrorMsg(msg string) {
	h.PrintMsg(h.colorString(fmt.Sprintf("<< ERROR: %s >>", msg), RED))
}

func (h *MyTerm) PrintInfoMsg(msg string) {
	h.PrintMsg(h.colorString(fmt.Sprintf("[%s] **** %s ****", time.Now().Format("2006-01-02 15:04:05"), msg), BLUE))
}

func (h *MyTerm) PrintTitle(msg string) {
	fmt.Fprintln(h.terminal, h.colorString(strings.Repeat("*", len(msg)+6), CYAN))
	fmt.Fprintf(h.terminal, h.colorString("** %s **\n", CYAN), msg)
	fmt.Fprintln(h.terminal, h.colorString(strings.Repeat("*", len(msg)+6), CYAN))
	h.Br()
}

func (h *MyTerm) ReadWithPrompt(msg string) (string, error) {
	h.promptLock = true
	defer func() {
		h.promptLock = false
	}()

	h.setTermPrompt(msg)
	r, err := h.Read()
	if err != nil {
		return "", err
	}

	return r, nil
}

func (h *MyTerm) PrintHighightedMsg(message string) {
	fmt.Fprintln(h.terminal, h.colorString(message, BLUE))
	fmt.Fprint(h.terminal, h.colorString(strings.Repeat("=", len(message)), BLUE))
	h.Br()
}

func (h *MyTerm) Br() {
	fmt.Fprintln(h.terminal)
}

func (h *MyTerm) Read() (string, error) {
	val, err := h.readTermLine()
	if err != nil {
		return "", err
	}

	if val == "" {
		return "", nil
	}

	h.Br()

	return val, nil
}

func (h *MyTerm) SetPrompt(prompt string) error {
	if h.promptLock {
		return nil
	}

	return h.setTermPrompt(prompt)
}

func (h *MyTerm) Write(p string, flags ...repl.Flag) error {
	if h.terminal == nil {
		return fmt.Errorf("no terminal started")
	}

	if p == "" {
		h.Br()
		return nil
	}

	responseType := repl.GetFlagValue(flags, "response-type")

	switch responseType {
	case repl.InfoMsgResponse:
		h.PrintInfoMsg(string(p))
	case repl.ErrorMsgResponse:
		h.PrintErrorMsg(string(p))
	case repl.HighlightedMsgResponse:
		h.PrintHighightedMsg(string(p))
	default:
		h.PrintMsg(string(p))
	}

	return nil
}
