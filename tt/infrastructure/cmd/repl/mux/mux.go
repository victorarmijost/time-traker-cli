package mux

import (
	"context"
	"fmt"
	"sort"
	"time"
	"varmijo/time-tracker/tt/infrastructure/cmd/repl"
)

type Mux struct {
	handlers map[string]handlerWithArgs
	helps    map[string]string
}

type handlerWithArgs struct {
	handler  repl.Handler
	argNames []string
}

func NewMux() *Mux {
	m := &Mux{
		handlers: map[string]handlerWithArgs{},
		helps:    map[string]string{},
	}

	m.Handle("help", repl.HandleFunc(m.help))

	return m
}

func (m *Mux) Handle(verb string, h repl.Handler, args ...string) {
	m.handlers[verb] = handlerWithArgs{handler: h, argNames: args}
}

func (m *Mux) AddHelp(verb, help string) error {
	if _, ok := m.helps[verb]; ok {
		panic(fmt.Sprintf("help for verb %s already exists", verb))
	}

	m.helps[verb] = help

	return nil
}

func (m *Mux) getHelp(verb string) string {
	help, ok := m.helps[verb]
	if !ok {
		return fmt.Sprintf("help for cmd %s not found", verb)
	}

	return help
}

func (m *Mux) getVerbs() []string {
	verbs := make([]string, 0, len(m.handlers))
	for verb := range m.handlers {
		verbs = append(verbs, verb)
	}

	return verbs
}

func (m *Mux) get(verb string) (handlerWithArgs, bool) {
	h, ok := m.handlers[verb]
	if !ok {
		return h, false
	}

	return h, true
}

func (m *Mux) help(r *repl.Request, w repl.IO) {
	repl.PrintHighightedMsg(w, "Command list")

	verbs := m.getVerbs()
	sort.Strings(verbs)

	for i, verb := range verbs {
		repl.PrintPlain(w, fmt.Sprintf("%d. {{ %s }} : %s", i, verb, m.getHelp(verb)))
	}
}

func cmdNotFound(r *repl.Request, w repl.IO) {
	repl.PrintErrorMsg(w, fmt.Sprintf("command not found: %s", r.Verb()))
}

func (m *Mux) Handler(req *repl.Request) (repl.Handler, []string) {
	h, ok := m.get(req.Verb())
	if !ok {
		return repl.HandleFunc(cmdNotFound), nil
	}

	return h.handler, h.argNames
}

func (m *Mux) ServeCmd(r *repl.Request, w repl.IO) {
	h, argNames := m.Handler(r)

	m.handleArgs(r, w, argNames)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r.SetNewCtx(ctx)

	h.ServeCmd(r, w)
}

func (c *Mux) handleArgs(r *repl.Request, w repl.IO, argNames []string) {
	argVals := r.ArgVals()
	remArgs := argNames

	i := 0
	newRemArgs := []string{}
	for _, argName := range remArgs {
		if i < len(argVals) {
			r.SetArg(argName, argVals[i])
		} else {
			newRemArgs = append(newRemArgs, argName)
		}
		i++
	}

	remArgs = newRemArgs

	if len(remArgs) > 0 {
		for _, argName := range remArgs {
			val, err := w.ReadWithPrompt(fmt.Sprintf("- %s: ", argName))
			if err != nil {
				panic(err)
			}

			r.SetArg(argName, val)
		}

	}
}
