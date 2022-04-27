package repl

import (
	"context"
	"time"
)

func (aCmd *withArgsFunc) Run(r *Handler, argVals ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	remArgs := aCmd.ArgNames

	args := make(map[string]string)

	if aCmd.Templates != nil {
		var tname string

		if len(argVals) == 0 {
			tname = r.GetInput("#")
		} else if len(argVals) > 0 {
			tname = argVals[0]
			argVals = argVals[1:]
		}

		temp, err := aCmd.Templates.Get(tname)

		if err != nil {
			r.PrintError(err)
			return
		}

		newRemArgs := []string{}
		for _, arg_name := range remArgs {
			if val, ok := temp[arg_name]; ok {
				args[arg_name] = val
			} else {
				newRemArgs = append(newRemArgs, arg_name)
			}
		}

		remArgs = newRemArgs
	}

	i := 0
	newRemArgs := []string{}
	for _, arg_name := range remArgs {
		if i < len(argVals) {
			args[arg_name] = argVals[i]
		} else {
			newRemArgs = append(newRemArgs, arg_name)
		}
		i++
	}

	remArgs = newRemArgs

	if len(remArgs) > 0 {
		for _, arg_name := range remArgs {
			val := r.GetInput(arg_name)
			args[arg_name] = val
		}

	}

	res, err := aCmd.Func(ctx, args)

	if err != nil {
		r.PrintError(err)
		return
	}

	r.PrintInfoMessage(res)
}

func (f ActionFunc) Run(r *Handler, argVals ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := f(ctx)

	if err != nil {
		r.PrintError(err)
		return
	}

	r.PrintInfoMessage(res)
}

func (f ActionFuncExt) WithArgs(t *TemplateHandler, ArgsNames ...string) SubRutine {
	return &withArgsFunc{
		Func:      f,
		ArgNames:  ArgsNames,
		Templates: t,
	}
}
