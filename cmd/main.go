package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"varmijo/time-tracker/bairestt"
	"varmijo/time-tracker/config"
	"varmijo/time-tracker/localStore"
	"varmijo/time-tracker/repl"
	"varmijo/time-tracker/state"
	"varmijo/time-tracker/utils"

	"github.com/sirupsen/logrus"
)

const logFile = "tt.log"

type Kernel struct {
	tt      *bairestt.Bairestt
	state   *state.State
	config  *config.Config
	recTemp *repl.TemplateHandler
}

func main() {
	file := setLogger()
	defer file.Close()

	state := initState()
	defer saveState(state)

	cmds := initCmds(state)

	cmds.PrintTitle("Welcome to BairedDev Time Tracker CLI tool")

	config := initConfig(cmds)

	tt := login(config, cmds)

	recTemp := repl.NewTemplateHandler("rec")
	err := recTemp.Load()

	if err != nil {
		fmt.Println("Warning: there is no record tamplete created")
		fmt.Println()
	}

	kern := &Kernel{
		state:   state,
		tt:      tt,
		config:  config,
		recTemp: recTemp,
	}

	if !config.IsComplete() {
		runConfig(cmds, kern)
	}

	registerFunctions(cmds, kern)

	cmds.Repl()
}

func login(config *config.Config, cmds *repl.Handler) *bairestt.Bairestt {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if config.Email == "" {
		log.Fatal("Missing email configuration")
	}

	tt := bairestt.NewService(config.Email)

	err := tt.Start(ctx)

	if err == nil {
		return tt
	}

	pass := config.Password
	if pass == "" {
		pass = cmds.GetPass("Google password")
	}

	ctx, cancelPass := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelPass()

	cmds.PrintInfoMessage("Performing Google login, it will take a while please wait...")
	err = tt.StartWithPass(ctx, pass)

	if err != nil {
		log.Fatal(err)
	}

	cmds.PrintInfoMessage("Login successfull!!")

	return tt
}

func initCmds(state *state.State) *repl.Handler {
	cmds := repl.NewHandler(getPrompt(state), "exit")

	return cmds
}

func initState() *state.State {
	state := state.NewState()

	err := state.Load()

	if errors.Is(err, os.ErrNotExist) {
		return state
	}

	if err != nil {
		log.Fatal(err)
	}

	return state
}

func initConfig(r *repl.Handler) *config.Config {
	c := config.NewConfig()

	err := c.Load()

	if err != nil {
		email := r.GetInput("Email address")

		if email == "" {
			log.Fatal("missing email configuration, application can't start")
		}
		c.Email = email
		err = c.Save()

		if err != nil {
			log.Fatal("config file can't be saved, application can't start")
		}
	}

	return c
}

func saveState(state *state.State) {
	err := state.Save()
	if err != nil {
		log.Fatal(err)
	}
}

func getPrompt(state *state.State) repl.Prompt {
	return func() string {
		statusBar := ""

		wt := localStore.GetTimeByStatus(state.Date, localStore.StatusPending)
		if wt > 0 {
			statusBar = fmt.Sprintf("[Worked:%.2f]", wt)
		}

		ct := localStore.GetTimeByStatus(state.Date, localStore.StatusCommited)
		if ct > 0 {
			statusBar = fmt.Sprintf("%s[Commited:%.2f]", statusBar, ct)
		}

		pt := localStore.GetTimeByStatus(nil, localStore.StatusPool)
		if pt > 0 {
			statusBar = fmt.Sprintf("%s[Pool:%.2f]", statusBar, pt)
		}

		if state.IsWorking() {
			statusBar = fmt.Sprintf("%s[Tracking:%.2f]", statusBar, state.GetTaskTime(nil))
		}

		if state.Date != nil {
			statusBar = fmt.Sprintf("%s[%s]", statusBar, state.Date.Format("06-01-02"))
		}

		if statusBar != "" {
			return fmt.Sprintf("%s tt", statusBar)
		}

		return "tt"
	}
}

func runConfig(r *repl.Handler, kern *Kernel) {
	ctx := context.Background()

	SetFocalPoint(kern)(ctx, r)
	SetProject(kern)(ctx, r)
	SetWorkingTime(kern)(ctx, r)
}

func setLogger() *os.File {
	file, err := os.OpenFile(utils.GeAppPath(logFile), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetOutput(file)
	logrus.SetLevel(logrus.DebugLevel)

	return file
}
