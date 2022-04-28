package main

import (
	"varmijo/time-tracker/repl"
)

func registerFunctions(cmds *repl.Handler, kern *Kernel) {
	rt := kern.recTemp

	//Tasks
	cmds.Handle("find task", SearchTask(kern).WithArgs(nil, "Name"))

	//Records
	cmds.Handle("add", AddRecord(kern).WithArgs(rt, "Id", "Comment", "Hours"))
	cmds.Handle("rec", StartRecord(kern).WithArgs(rt, "Id", "Comment"))
	cmds.Handle("end", StopRecord(kern))
	cmds.Handle("end at", StopRecordAt(kern).WithArgs(nil, "At"))
	cmds.Handle("commit", CommitAll(kern))
	cmds.Handle("drop", DropRecord(kern))
	cmds.Handle("edit", EditRecord(kern).WithArgs(rt, "Id", "Comment"))
	cmds.Handle("rec at", StartRecordAt(kern).WithArgs(rt, "Id", "Comment", "At"))
	cmds.Handle("list", ListLocal(kern))
	cmds.Handle("view", ViewRecord(kern))
	cmds.Handle("edit stored", EditStoredRecord(kern))
	cmds.Handle("poure", PourePool(kern))
	cmds.Handle("delete", DeleteStoredRecord(kern))

	//Navigate
	cmds.Handle("change date", ChangeDate(kern).WithArgs(nil, "Date"))

	//Templates
	cmds.Handle("temp add", CreateTemplate(kern))
	cmds.Handle("temp list", ListTemplates(kern))

	addHelp(cmds)
}

func addHelp(cmds *repl.Handler) {
	//Tasks
	cmds.Help("find task", "Search a task, by its name, on tasks master got from the Time Tracker.")

	//Records
	cmds.Help("add", "Adds a new task record.")
	cmds.Help("rec", "Starts a new time recorer.")
	cmds.Help("end", "End the current time recorder, base on the initial time calculates the spent time.")
	cmds.Help("end at", "Similar to End but you can set the hour when the time recorded ended.")
	cmds.Help("commit", "Send all the pending time on current date to the Time Tracker.")
	cmds.Help("drop", "Drops the current working time recorder, all the information will be lost.")
	cmds.Help("edit", "Allows to modify the current time recorder, keeping the start time of the original record.")
	cmds.Help("rec at", "Allows to start a time recorder at an specific hour.")
	cmds.Help("list", "List all the records on the current date.")
	cmds.Help("view", "Allow to view the current time recorder.")
	cmds.Help("edit stored", "Allows to edit a non commited recored.")
	cmds.Help("poure", "Poures all the time on the pool to the current date.")
	cmds.Help("delete", "Allows to delete a non commited recored.")

	//Navigate
	cmds.Help("change date", "Allow to change the current working date.")

	//Templates
	cmds.Help("temp add", "Adds a new record template.")
	cmds.Help("temp list", "List all the existing templates.")
}
