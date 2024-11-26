package repl

type Flag struct {
	name  string
	value string
}

func NewFlag(key, value string) Flag {
	return Flag{name: key, value: value}
}

func GetFlagValue(flags []Flag, name string) string {
	for _, flag := range flags {
		if flag.name == name {
			return flag.value
		}
	}
	return ""
}

type IO interface {
	Write(string, ...Flag) error
	SetPrompt(string) error
	Read() (string, error)
	ReadWithPrompt(string) (string, error)
}
