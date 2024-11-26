package handlers

import (
	"varmijo/time-tracker/pkg/repl"
)

func (h *Handlers) Register() {
	//Records
	h.mux.Handle("add", repl.HandleFunc(h.AddRecord), "Hours")
	h.mux.Handle("rec", repl.HandleFunc(h.StartRecord))
	h.mux.Handle("rec at", repl.HandleFunc(h.StartRecordAt), "At")

	h.mux.Handle("end", repl.HandleFunc(h.StopRecord))
	h.mux.Handle("end at", repl.HandleFunc(h.StopRecordAt), "At")
	h.mux.Handle("commit", repl.HandleFunc(h.CommitAll), "Amount")
	h.mux.Handle("send pool", repl.HandleFunc(h.SendToPool))
	h.mux.Handle("drop", repl.HandleFunc(h.DropRecord))
	h.mux.Handle("list", repl.HandleFunc(h.ListLocal))
	//h.mux.Handle("edit stored", repl.HandleFunc(h.EditStoredRecord))
	h.mux.Handle("poure", repl.HandleFunc(h.PourePool))
	//h.mux.Handle("delete", repl.HandleFunc(h.DeleteStoredRecord))

	//Navigate
	h.mux.Handle("change date", repl.HandleFunc(h.ChangeDate), "Date")

	//Stats
	h.mux.Handle("debt", repl.HandleFunc(h.GetDebts))
}
