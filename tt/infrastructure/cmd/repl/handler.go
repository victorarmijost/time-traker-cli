package repl

type Handler interface {
	ServeCmd(*Request, IO)
}

type HandleFunc func(*Request, IO)

func (f HandleFunc) ServeCmd(req *Request, w IO) {
	f(req, w)
}
