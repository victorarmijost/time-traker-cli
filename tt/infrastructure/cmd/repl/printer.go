package repl

const (
	PlainResult            string = "plain-result"
	InfoMsgResponse        string = "info-msg"
	ErrorMsgResponse       string = "error-msg"
	HighlightedMsgResponse string = "highlighted-msg"
)

func Br(w IO) {
	err := w.Write("\n")
	if err != nil {
		panic(err)
	}
}

func printResponseWithType(w IO, msg string, responseType string) {
	err := w.Write(msg, NewFlag("response-type", responseType))
	if err != nil {
		panic(err)
	}
}

func PrintPlain(w IO, msg string) {
	printResponseWithType(w, msg, PlainResult)
}

func PrintInfoMsg(w IO, msg string) {
	printResponseWithType(w, msg, InfoMsgResponse)
}

func PrintErrorMsg(w IO, msg string) {
	printResponseWithType(w, msg, ErrorMsgResponse)
}

func PrintHighightedMsg(w IO, msg string) {
	printResponseWithType(w, msg, HighlightedMsgResponse)
}

func PrintError(w IO, err error) {
	if err == nil {
		return
	}
	PrintErrorMsg(w, err.Error())
}
