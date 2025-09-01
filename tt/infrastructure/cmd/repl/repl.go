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
	data    domain.PromptData
	exit    string
}

func NewRepl(data domain.PromptData, handler Handler, io IO, exit string) *Repl {
	h := &Repl{
		io:      io,
		handler: handler,
		data:    data,
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
	go func() {
		for {
			time.Sleep(1 * time.Second)

			c.promptUpdate(false)
		}
	}()
}

func (c *Repl) promptUpdate(forceRefresh bool) {
	statusBar := ""

	if forceRefresh {
		c.data.RefreshData()
	}

	if c.data.Dt() > 0 {
		statusBar = fmt.Sprintf("[Debt:%s]", domain.FormatDuration(c.data.Dt()))
	}

	if c.data.Wt() > 0 {
		statusBar = fmt.Sprintf("%s[Worked:%s]", statusBar, domain.FormatDuration(c.data.Wt()))
	}

	if c.data.Ct() > 0 {
		statusBar = fmt.Sprintf("%s[Commited:%s]", statusBar, domain.FormatDuration(c.data.Ct()))
	}

	if c.data.Pt() > 0 {
		statusBar = fmt.Sprintf("%s[Pool:%s]", statusBar, domain.FormatDuration(c.data.Pt()))
	}

	if c.data.IsWorking() {
		statusBar = fmt.Sprintf("%s[Rec:%s][%s]", statusBar, domain.FormatDuration(c.data.Tt()), getClockEmoji())
	}

	if !c.data.IsToday() {
		statusBar = fmt.Sprintf("%s[%s]", statusBar, c.data.GetDate().Format("06-01-02"))
	}

	if statusBar != "" {
		statusBar = fmt.Sprintf("%s tt", statusBar)
	} else {
		statusBar = "tt"
	}

	err := c.io.SetPrompt(fmt.Sprintf("%s > ", statusBar))
	if err != nil {
		panic(err)
	}
}

func getClockEmoji() string {
	clocks := []string{".", "'"}

	n := int64(len(clocks))

	return strings.TrimSpace(clocks[time.Now().Unix()%n])
}

func (c *Repl) replIter() string {
	c.promptUpdate(true)
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
