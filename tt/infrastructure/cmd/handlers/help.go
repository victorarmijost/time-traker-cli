package handlers

func (h *Handlers) AddHelp() {
	//Records
	h.mux.AddHelp("add", "Adds a new task record.")
	h.mux.AddHelp("rec", "Starts a new time recorer.")
	h.mux.AddHelp("end", "End the current time recorder, base on the initial time calculates the spent time.")
	h.mux.AddHelp("end at", "Similar to End but you can set the hour when the time recorded ended.")
	h.mux.AddHelp("commit", "Send all the pending time on current date to the Time Tracker.")
	h.mux.AddHelp("send pool", "Sends all the pending time to the pool.")
	h.mux.AddHelp("drop", "Drops the current working time recorder, all the information will be lost.")
	h.mux.AddHelp("edit", "Allows to modify the current time recorder, keeping the start time of the original record.")
	h.mux.AddHelp("rec at", "Allows to start a time recorder at an specific hour.")
	h.mux.AddHelp("list", "List all the records on the current date.")
	h.mux.AddHelp("view", "Allow to view the current time recorder.")
	h.mux.AddHelp("edit stored", "Allows to edit a non commited recored.")
	h.mux.AddHelp("poure", "Poures all the time on the pool to the current date.")
	h.mux.AddHelp("delete", "Allows to delete a non commited recored.")

	//Navigate
	h.mux.AddHelp("change date", "Allow to change the current working date.")

	//Templates
	h.mux.AddHelp("temp add", "Adds a new record template.")
	h.mux.AddHelp("temp list", "List all the existing templates.")
}
