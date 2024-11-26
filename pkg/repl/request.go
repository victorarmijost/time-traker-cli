package repl

import (
	"context"
	"errors"
)

type Request struct {
	ctx      context.Context
	verb     string
	argsVals []string
	args     map[string]string
}

func NewRequest(ctx context.Context, verb string, argsVals []string) *Request {
	return &Request{
		ctx:      ctx,
		verb:     verb,
		argsVals: argsVals,
		args:     map[string]string{},
	}
}

func (r *Request) Verb() string {
	return r.verb
}

var ErrArgNotFound = errors.New("arg not found")

func (r *Request) Arg(name string) (string, error) {
	val, ok := r.args[name]
	if !ok {
		return "", ErrArgNotFound
	}
	return val, nil
}

func (r *Request) Ctx() context.Context {
	return r.ctx
}

func (r *Request) ArgVals() []string {
	return r.argsVals
}

func (r *Request) SetNewCtx(ctx context.Context) {
	r.ctx = ctx
}

func (r *Request) SetArg(name, value string) {
	r.args[name] = value
}
