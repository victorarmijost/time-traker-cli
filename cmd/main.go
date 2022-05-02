package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"
	"varmijo/time-tracker/bairestt"
	"varmijo/time-tracker/config"
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
	state := initState()
	defer saveState(state)

	cmds, closeTerm := initCmds(state)
	defer closeTerm()

	cmds.PrintTitle("Welcome to BairedDev Time Tracker CLI tool")

	config := initConfig(cmds)

	file := setLogger(config.LogLevel)
	defer file.Close()

	tt := login(config, cmds)

	recTemp := repl.NewTemplateHandler("rec")
	err := recTemp.Load()

	if err != nil {
		cmds.PrintErrorMsg("Warning: there is no record tamplete created")
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

func initCmds(state *state.State) (*repl.Handler, repl.CloseTerm) {
	cmds, close := repl.NewHandler(getPrompt(state), "exit")

	return cmds, close
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

func runConfig(r *repl.Handler, kern *Kernel) {
	ctx := context.Background()

	SetFocalPoint(kern)(ctx, r)
	SetProject(kern)(ctx, r)
	SetWorkingTime(kern)(ctx, r)
}

func setLogger(slevel string) *os.File {
	file, err := os.OpenFile(utils.GeAppPath(logFile), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetOutput(file)

	level, err := logrus.ParseLevel(slevel)

	if err != nil {
		logrus.SetLevel(logrus.ErrorLevel)
		return file
	}

	logrus.SetLevel(level)

	return file
}
