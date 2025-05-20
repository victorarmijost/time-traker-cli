package repl

func ParseArg[T any](r *Request, name string, p func(string) (T, error)) (T, error) {
	val, err := r.Arg(name)

	var zero T

	if err != nil {
		return zero, err
	}

	return p(val)
}
