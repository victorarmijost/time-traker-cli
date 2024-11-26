package main

import (
	"log"
	"os"

	"varmijo/time-tracker/pkg/repl"
	"varmijo/time-tracker/pkg/repl/myterm"
	"varmijo/time-tracker/pkg/utils"
	"varmijo/time-tracker/tt/app"
	"varmijo/time-tracker/tt/infrastructure/cmd/handlers"
	"varmijo/time-tracker/tt/infrastructure/config"
	"varmijo/time-tracker/tt/infrastructure/repositories"

	"github.com/sirupsen/logrus"
)

const logFile = "tt.log"

func main() {
	cfg := config.MustNewConfig()

	file := setLogger(cfg.GetLogLevel())
	defer file.Close()

	// Create sqlite DB
	db, err := repositories.NewSQLiteDB("tt")
	if err != nil {
		log.Fatal(err)
	}

	records := repositories.NewSQLiteRecordRepository(db)

	track := repositories.NewSQLiteTrackRepository(db)

	app := app.NewApp(cfg, records, track)

	mux := handlers.NewHandlers(app)

	term, close := myterm.NewTerm()
	defer close()

	term.PrintTitle("Welcome to Time Tracker CLI tool")

	cmds := repl.NewRepl(app.GetPrompt(), mux, term, "exit")

	cmds.Run()
}

// Set up the application logger
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
