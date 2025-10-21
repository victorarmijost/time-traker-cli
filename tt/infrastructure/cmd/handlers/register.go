package handlers

import (
	"varmijo/time-tracker/tt/infrastructure/cmd/repl"
)

func (h *Handlers) Register() {
	//Records
	h.mux.Handle("add", repl.HandleFunc(h.AddRecord), "Hours")
	h.mux.Handle("rec", repl.HandleFunc(h.StartRecord))
	h.mux.Handle("rec at", repl.HandleFunc(h.StartRecordAt), "At")

	h.mux.Handle("end", repl.HandleFunc(h.StopRecord))
	h.mux.Handle("end at", repl.HandleFunc(h.StopRecordAt), "At")
	h.mux.Handle("drop", repl.HandleFunc(h.DropRecord))

	//Navigate
	h.mux.Handle("change date", repl.HandleFunc(h.ChangeDate), "Date")
}
