package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/term"
)

type Handler struct {
	terminal *term.Terminal
	cmds     map[string]SubRutine
	reader   *bufio.Reader
	prompt   Prompt
	exit     string
	cmdsHelp map[string]string
}

func NewHandler(prompt Prompt, exit string) (*Handler, CloseTerm) {
	term, close := setupTerm()

	h := &Handler{
		terminal: term,
		cmds:     map[string]SubRutine{},
		reader:   bufio.NewReader(os.Stdin),
		prompt:   prompt,
		exit:     exit,
		cmdsHelp: map[string]string{},
	}

	go h.updatePromptBackground()

	return h, close
}

func (c *Handler) Handle(name string, sr SubRutine) {
	c.cmds[name] = sr
}

func (c *Handler) Help(name, help string) {
	if _, ok := c.cmds[name]; ok {
		c.cmdsHelp[name] = help
	}
}

func (c *Handler) getHelp(name string) string {
	if help, ok := c.cmdsHelp[name]; ok {
		return help
	}

	return "Help description is not added!"
}

func (c *Handler) get(name string) SubRutine {
	if f, ok := c.cmds[name]; ok {
		return f
	}

	return ActionFunc(func(ctx context.Context) (string, error) {
		return "", fmt.Errorf("command not found")
	})

}

func (c *Handler) runSubRutine(cmd string) {
	verb, argVals := parseCmd(cmd)
	sr := c.get(verb)

	c.Br()
	sr.Run(c, argVals...)
}

func (c *Handler) getInput() string {
	input, err := c.readTermLine()
	if err != nil {
		os.Exit(0)
	}

	return input
}

func (c *Handler) updatePromptBackground() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			c.printPrompt()
			c.terminal.Write([]byte{})
		}
	}()
}

func (c *Handler) printPrompt() {
	c.setTermPrompt(fmt.Sprintf("%s > ", c.prompt()))

}

func (c *Handler) replIter() string {
	c.printPrompt()
	return c.getInput()
}

func parseCmd(cmd string) (verb string, args []string) {
	parts := strings.Split(cmd, ";")

	verb = strings.TrimSpace(parts[0])

	if len(parts) > 1 {
		for _, s := range parts[1:] {
			args = append(args, strings.TrimSpace(s))
		}
	}

	return
}

func (c *Handler) shouldContinue(cmd string) bool {
	return !strings.EqualFold(c.exit, cmd)
}

func (c *Handler) Repl() {
	for cmd := c.replIter(); c.shouldContinue(cmd); cmd = c.replIter() {
		if cmd == "" {
			continue
		}

		if strings.ToLower(cmd) == "help" {
			c.help()
			continue
		}

		c.runSubRutine(cmd)
	}
}

func (c *Handler) help() {
	c.PrintHighightedMessage("Command list")
	c.Br()

	keys := make([]string, 0, len(c.cmds))
	for k := range c.cmds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	i := 0
	for _, cmd := range keys {
		i++
		fmt.Fprintf(c.terminal, "%d. {{ %s }} : %s\n\n", i, cmd, c.getHelp(cmd))
	}
}
