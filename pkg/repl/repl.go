package repl

import (
	"context"
	"fmt"
	"strings"
	"time"
	"varmijo/time-tracker/tt/domain"
)

type Repl struct {
	io      IO
	handler Handler
	prompt  domain.Prompt
	exit    string
}

func NewRepl(prompt domain.Prompt, handler Handler, io IO, exit string) *Repl {
	h := &Repl{
		io:      io,
		handler: handler,
		prompt:  prompt,
		exit:    exit,
	}

	go h.updatePromptBackground()

	return h
}

func (c *Repl) Run() {
	for cmd := c.replIter(); c.shouldContinue(cmd); cmd = c.replIter() {
		c.serveCmd(cmd)
	}
}

func (c *Repl) serveCmd(cmd string) {
	if cmd == "" {
		return
	}

	verb, argVals := parseCmd(cmd)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := NewRequest(ctx, verb, argVals)

	c.handler.ServeCmd(req, c.io)
}

func (c *Repl) updatePromptBackground() {
	const hardUpdateCount = 60

	count := hardUpdateCount
	go func() {
		for {
			time.Sleep(1 * time.Second)

			if count <= 0 {
				c.promptUpdate(domain.FULL_UPDATE)
				count = hardUpdateCount
			} else {
				c.promptUpdate(domain.SOFT_UPDATE)
				count--
			}
		}
	}()
}

func (c *Repl) promptUpdate(t domain.PromptType) {
	err := c.io.SetPrompt(fmt.Sprintf("%s > ", c.prompt(t)))
	if err != nil {
		panic(err)
	}
}

func (c *Repl) replIter() string {
	c.promptUpdate(domain.FULL_UPDATE)
	val, err := c.io.Read()
	if err != nil {
		panic(err)
	}

	return val
}

func (c *Repl) shouldContinue(cmd string) bool {
	return !strings.EqualFold(c.exit, cmd)
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
