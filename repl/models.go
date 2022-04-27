package repl

import (
	"context"
)

type ActionFunc func(ctx context.Context) (string, error)
type ActionFuncExt func(ctx context.Context, args map[string]string) (string, error)

type Prompt func() string

type withArgsFunc struct {
	Func      ActionFuncExt
	ArgNames  []string
	Templates *TemplateHandler
}

type InteractiveFunc func(ctx context.Context, r *Handler)

type SubRutine interface {
	Run(r *Handler, argVals ...string)
}
