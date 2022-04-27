package repl

import "context"

func (i InteractiveFunc) Run(r *Handler, argVals ...string) {
	//Let the user handle the context
	ctx := context.Background()

	i(ctx, r)
}

func (i InteractiveFunc) WithArgs(t *TemplateHandler, ArgsNames ...string) SubRutine {
	return i
}
