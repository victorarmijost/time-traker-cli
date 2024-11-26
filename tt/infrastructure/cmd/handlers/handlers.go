package handlers

import (
	"fmt"
	"strconv"
	"time"
	"varmijo/time-tracker/pkg/repl"
	"varmijo/time-tracker/tt/app"
	"varmijo/time-tracker/tt/domain"

	"varmijo/time-tracker/pkg/repl/mux"
)

type Handlers struct {
	kern *app.App
	mux  *mux.Mux
}

func NewHandlers(kern *app.App) repl.Handler {
	h := &Handlers{
		kern: kern,
		mux:  mux.NewMux(),
	}

	h.Register()
	h.AddHelp()

	return h.GetMux()
}

func (h *Handlers) GetMux() *mux.Mux {
	return h.mux
}

func (h *Handlers) AddRecord(r *repl.Request, w repl.IO) {
	phours, err := repl.ParseArg(r, "Hours", domain.ParseDuration)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	err = h.kern.AddRecord(r.Ctx(), phours)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, fmt.Sprintf("%0.2f hours inserted!", phours))
}

func (h *Handlers) StartRecord(r *repl.Request, w repl.IO) {
	err := h.kern.StartRecord(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, "Record started!")
}

func (h *Handlers) StartRecordAt(r *repl.Request, w repl.IO) {
	pat, err := repl.ParseArg(r, "At", domain.ParseHour)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	err = h.kern.StartRecordAt(r.Ctx(), pat)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, fmt.Sprintf("Record started at %s!", pat.Format("15:04")))

}

func (h *Handlers) StopRecord(r *repl.Request, w repl.IO) {
	hours, err := h.kern.StopRecord(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, fmt.Sprintf("%0.2f hours inserted!", hours))
}

func (h *Handlers) StopRecordAt(r *repl.Request, w repl.IO) {
	pat, err := repl.ParseArg(r, "At", domain.ParseHour)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	hours, err := h.kern.StopRecordAt(r.Ctx(), pat)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, fmt.Sprintf("%0.2f hours inserted!", hours))
}

func (h *Handlers) CommitAll(r *repl.Request, w repl.IO) {
	amount, err := repl.ParseArg(r, "Amount", func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})

	var pamount *float64
	if err != nil {
		pamount = nil
	} else {
		pamount = &amount
	}

	err = h.kern.CommitAll(r.Ctx(), pamount)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, "Records commited!")
}

func (h *Handlers) SendToPool(r *repl.Request, w repl.IO) {
	err := h.kern.SendToPool(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, "Records saved to pool!")
}

func (h *Handlers) DropRecord(r *repl.Request, w repl.IO) {
	hours, err := h.kern.DropRecord(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, fmt.Sprintf("%0.2f hours dropped!", hours))
}

func (h *Handlers) ListLocal(r *repl.Request, w repl.IO) {
	list, err := h.kern.ListLocal(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	if len(list) == 0 {
		repl.PrintInfoMsg(w, "No records found")
		return
	}

	repl.PrintHighightedMsg(w, "Result")
	for i, r := range list {
		repl.PrintPlain(w, fmt.Sprintf("%d. %s", i+1, r))
	}
}

func (h *Handlers) PourePool(r *repl.Request, w repl.IO) {
	err := h.kern.PourePool(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, "Pool poured!")
}

func (h *Handlers) ChangeDate(r *repl.Request, w repl.IO) {
	date, err := repl.ParseArg(r, "Date", domain.GetDateFromText)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	err = h.kern.ChangeDate(r.Ctx(), date)
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	repl.PrintInfoMsg(w, "Date change!")
}

func (h *Handlers) DeleteStoredRecord(r *repl.Request, w repl.IO) {
}

func (h *Handlers) EditStoredRecord(r *repl.Request, w repl.IO) {
}

func (h *Handlers) GetDebts(r *repl.Request, w repl.IO) {
	debt, err := h.kern.GetDebt(r.Ctx())
	if err != nil {
		repl.PrintError(w, err)
		return
	}

	if debt.Length() == 0 {
		repl.PrintInfoMsg(w, "No debt found")
		return
	}

	repl.PrintHighightedMsg(w, "Result")
	debt.Do(func(date time.Time, hours float64) {
		repl.PrintPlain(w, fmt.Sprintf("%s: %s", date.Format("2006-01-02"), domain.FormatDuration(hours)))
	})

	repl.PrintPlain(w, fmt.Sprintf("Total debt: %s", domain.FormatDuration(debt.Total())))
}
