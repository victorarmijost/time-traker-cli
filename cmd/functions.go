package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"varmijo/time-tracker/localStore"
	"varmijo/time-tracker/repl"
)

//Add a new record
func AddRecord(kern *Kernel) repl.ActionFuncExt {
	return func(ctx context.Context, args map[string]string) (string, error) {
		phours, err := strconv.ParseFloat(args["Hours"], 32)

		if err != nil {
			return "", err
		}

		hours := kern.state.TaskTimeRounder(float32(phours))

		var recDate time.Time
		if kern.state.Date == nil {
			recDate = time.Now()
		} else {
			recDate = *kern.state.Date
		}

		if len(args["Comment"]) < 10 {
			return "", fmt.Errorf("comment must have more than 10 characters")
		}

		//Pending - Hardcoded record type 1
		record := &localStore.Record{
			TaskName: args["Task Name"],
			Date:     recDate,
			Comments: args["Comment"],
			Hours:    float32(hours),
		}

		err = localStore.Save(record)

		if err != nil {
			err := fmt.Errorf("new record can't be inserted, %w", err)
			return "", err
		}

		return fmt.Sprintf("%s - %0.2f hours inserted!", record.Comments, record.Hours), nil
	}
}

//Add a new record
func StartRecord(kern *Kernel) repl.ActionFuncExt {
	return func(ctx context.Context, args map[string]string) (string, error) {
		if kern.state.Date != nil {
			return "", fmt.Errorf("wrong date, change back to today")
		}

		state := kern.state
		defer state.Save()

		if len(args["Comment"]) < 10 {
			return "", fmt.Errorf("comment must have more than 10 characters")
		}

		err := state.StartRecord(args["Task Name"], args["Comment"], nil)

		if err != nil {
			return "", err
		}

		return "Record started!", nil
	}
}

func StopRecord(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		currentTask, err := kern.state.GetCurrentTask()

		if err != nil {
			return "", err
		}

		state := kern.state
		defer state.Save()

		comment := currentTask.Comment
		recDate := currentTask.StartTime
		taskName := currentTask.TaskName

		hours, err := state.EndRecord(nil)

		if err != nil {
			return "", err
		}

		//Pending - Hardcoded record type 1
		record := &localStore.Record{
			Date:     recDate,
			TaskName: taskName,
			Comments: comment,
			Hours:    float32(hours),
		}

		err = localStore.Save(record)

		if err != nil {
			err = fmt.Errorf("new rekernd can't be inserted, %w", err)
			return "", err
		}

		return fmt.Sprintf("%s - %0.2f hours inserted!", record.Comments, record.Hours), nil

	}
}

func StopRecordAt(kern *Kernel) repl.ActionFuncExt {
	return func(ctx context.Context, args map[string]string) (string, error) {
		state := kern.state
		defer state.Save()

		currentTask, err := state.GetCurrentTask()

		if err != nil {
			return "", err
		}

		comment := currentTask.Comment
		recDate := currentTask.StartTime

		endTime, err := parseHour(args["At"])

		if err != nil {
			return "", err
		}

		hours, err := state.EndRecord(endTime)

		if err != nil {
			return "", err
		}

		record := &localStore.Record{
			Date:     recDate,
			Comments: comment,
			Hours:    float32(hours),
		}

		err = localStore.Save(record)

		if err != nil {
			err = fmt.Errorf("error inserting new record, %w", err)
			return "", err
		}

		return fmt.Sprintf("%s - %0.2f hours inserted!", record.Comments, record.Hours), nil

	}
}

func CommitAll(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		state := kern.state

		files, err := localStore.ListByStatus(state.Date, localStore.StatusPending)

		if err != nil {
			return "", err
		}

		commitedTime := localStore.GetTimeByStatus(state.Date, localStore.StatusCommited)
		remTime := kern.config.WorkingTime - commitedTime
		for _, f := range files {
			record, err := localStore.Get(state.Date, f)
			if err != nil {
				return "", err
			}

			if remTime < record.Hours {
				newRecord := *record
				newRecord.Hours = record.Hours - remTime

				err = localStore.SaveToPool(&newRecord)
				if err != nil {
					return "", err
				}

				if remTime == 0 {
					err = localStore.DeleteRecord(&record.Date, record.Id)
					if err != nil {
						return "", err
					}

					//If there is not remaining time the time is not sent to the TT
					continue

				} else {
					record.Hours = remTime

					err = localStore.Save(record)
					if err != nil {
						return "", err
					}
				}
			}

			//If the record is 0 size, is marked as commited, but is not sent
			if record.Hours <= 0 {
				continue
			}

			err = localStore.SetCommit(state.Date, f)
			if err != nil {
				return "", err
			}

			remTime -= record.Hours
		}

		return "Records commited!", nil
	}
}

func SendToPool(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		state := kern.state

		files, err := localStore.ListByStatus(state.Date, localStore.StatusPending)

		if err != nil {
			return "", err
		}

		for _, f := range files {
			record, err := localStore.Get(state.Date, f)
			if err != nil {
				return "", err
			}

			newRecord := *record
			newRecord.Hours = record.Hours

			err = localStore.SaveToPool(&newRecord)
			if err != nil {
				return "", err
			}

			err = localStore.DeleteRecord(&record.Date, record.Id)
			if err != nil {
				return "", err
			}
		}

		return "Records saved to pool!", nil
	}
}

func ListLocal(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		state := kern.state

		pendingRecords, err := localStore.GetAllByStatus(state.Date, localStore.StatusPending)

		if err != nil {
			err = fmt.Errorf("can't get records, %w", err)
			return "", err
		}

		list := []string{}
		for _, r := range pendingRecords {
			list = append(list, fmt.Sprintf("[%.2f] - %s", r.Hours, r.Comments))
		}

		commitedRecords, err := localStore.GetAllByStatus(state.Date, localStore.StatusCommited)

		if err != nil {
			err = fmt.Errorf("can't get records, %w", err)
			return "", err
		}

		for _, r := range commitedRecords {
			list = append(list, fmt.Sprintf("[%.2f] - %s ✔️", r.Hours, r.Comments))
		}

		return repl.SprintList(list), nil
	}
}

func ChangeDate(kern *Kernel) repl.ActionFuncExt {
	return func(ctx context.Context, args map[string]string) (string, error) {
		state := kern.state
		defer state.Save()

		date, err := getDateFromText(args["Date"])

		if err != nil {
			err = fmt.Errorf("wrong date")
			return "", err
		}

		state.Date = date

		return "Date change!", nil
	}
}

func DropRecord(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		state := kern.state
		defer state.Save()

		currentTask, err := state.GetCurrentTask()

		if err != nil {
			return "", err
		}

		comment := currentTask.Comment

		hours, err := state.EndRecord(nil)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s - %0.2f hours dropped!", comment, hours), nil
	}
}

func StartRecordAt(kern *Kernel) repl.ActionFuncExt {
	return func(ctx context.Context, args map[string]string) (string, error) {
		state := kern.state
		defer state.Save()

		recDate, err := parseHour(args["At"])

		if err != nil {
			return "", err
		}

		if len(args["Comment"]) < 10 {
			return "", fmt.Errorf("comment must have more than 10 characters")
		}

		err = state.StartRecord(args["Task Name"], args["Comment"], recDate)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Record started at %s!", args["At"]), nil
	}
}

func EditRecord(kern *Kernel) repl.ActionFuncExt {
	return func(ctx context.Context, args map[string]string) (string, error) {
		state := kern.state
		defer state.Save()

		currentTask, err := state.GetCurrentTask()

		if err != nil {
			return "", err
		}

		date := currentTask.StartTime

		_, err = state.EndRecord(nil)

		if err != nil {
			return "", err
		}

		if len(args["Comment"]) < 10 {
			return "", fmt.Errorf("comment must have more than 10 characters")
		}

		err = state.StartRecord(args["Task Name"], args["Comment"], &date)

		if err != nil {
			return "", err
		}

		return "Record edited!", nil
	}
}

func ViewRecord(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		currentTask, err := kern.state.GetCurrentTask()

		if err != nil {
			return "", err
		}

		res := map[string]string{
			"Comment": currentTask.Comment,
		}

		return repl.SprintMap(res), nil
	}
}

func CreateTemplate(kern *Kernel) repl.InteractiveFunc {
	return func(ctx context.Context, r *repl.Handler) {

		tempName := r.GetInput("Template name")

		_, err := kern.recTemp.Get(tempName)

		if err == nil {
			r.PrintErrorMsg("Template already exist")
			return
		}

		temp := repl.Template{}

		taskName := r.GetInput("Task Name")

		if taskName != "" {
			temp["Task Name"] = taskName
		} else {
			r.PrintErrorMsg("Missing task name")
			return
		}

		if comment := r.GetInput("Comment (enter to skip)"); comment != "" {
			temp["Comment"] = comment
		}

		if hours := r.GetInput("Hours (enter to skip)"); hours != "" {
			temp["Hours"] = hours
		}

		if description := r.GetInput("Description (enter to skip)"); description != "" {
			temp["x-description"] = description
		}

		err = kern.recTemp.Save(tempName, temp)
		if err != nil {
			r.PrintError(err)
			return
		}

		r.PrintInfoMessage("Template added!")
	}
}

func ListTemplates(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		tmps := kern.recTemp.List()

		return repl.SprintList(tmps), nil
	}
}

func SetWorkingTime(kern *Kernel) repl.InteractiveFunc {
	return func(ctx context.Context, r *repl.Handler) {
		r.PrintTitle("Configure your working time")

		wts := r.GetInput("Working time per day (hours)")

		wt, err := strconv.ParseFloat(wts, 32)

		if err != nil {
			r.PrintError(err)
			return
		}

		kern.config.WorkingTime = float32(wt)

		err = kern.config.Save()
		if err != nil {
			r.PrintError(err)
			return
		}

		r.PrintInfoMessage("Working time added!")
	}
}

type RecordSearch []*localStore.Record

func (t RecordSearch) GetElement(i int) string {
	return fmt.Sprintf("[%.2f] - %s", t[i].Hours, t[i].Comments)
}

func (t RecordSearch) Size() int {
	return len(t)
}

func EditStoredRecord(kern *Kernel) repl.InteractiveFunc {
	return func(ctx context.Context, r *repl.Handler) {
		records, err := localStore.GetAllByStatus(kern.state.Date, localStore.StatusPending)

		if err != nil {
			r.PrintError(err)
			return
		}

		id := r.SelectFromList(RecordSearch(records))

		if id < 0 {
			return
		}

		record := records[id]

		r.PrintHighightedMessage("Current data")

		r.PrintMap(map[string]string{
			"Duration": fmt.Sprintf("%.2f", record.Hours),
			"Comment":  record.Comments,
		})

		cont := r.GetInput("Select an option (e: edit, q: cancel)")

		if strings.ToLower(cont) == "q" {
			r.PrintMessage("Canceled!")
			return
		}

		if strings.ToLower(cont) != "e" {
			r.PrintMessage("Wrong input!")
			return
		}

		r.PrintHighightedMessage("Editing data")

		AddRecord(kern).WithArgs(kern.recTemp, "Task Name", "Comment", "Hours").Run(r)

		err = localStore.DeleteRecord(kern.state.Date, record.Id)

		if err != nil {
			r.PrintError(err)
			return
		}
	}
}

func DeleteStoredRecord(kern *Kernel) repl.InteractiveFunc {
	return func(ctx context.Context, r *repl.Handler) {
		records, err := localStore.GetAllByStatus(kern.state.Date, localStore.StatusPending)

		if err != nil {
			r.PrintError(err)
			return
		}

		id := r.SelectFromList(RecordSearch(records))

		if id < 0 {
			return
		}

		record := records[id]

		r.PrintHighightedMessage("Data that will be deleted")

		r.PrintMap(map[string]string{
			"Duration": fmt.Sprintf("%.2f", record.Hours),
			"Comment":  record.Comments,
		})

		cont := r.GetInput("Select an option (d: delete, q: cancel)")

		if strings.ToLower(cont) == "q" {
			r.PrintMessage("Canceled!")
			return
		}

		if strings.ToLower(cont) != "d" {
			r.PrintMessage("Wrong input!")
			return
		}

		err = localStore.DeleteRecord(kern.state.Date, record.Id)

		if err != nil {
			r.PrintError(err)
			return
		}

		r.PrintMessage("Record deleted!")
	}
}

func PourePool(kern *Kernel) repl.ActionFunc {
	return func(ctx context.Context) (string, error) {
		err := localStore.PourePool(kern.state.Date)
		if err != nil {
			return "", err
		}

		return "Pool poured!", nil
	}
}
