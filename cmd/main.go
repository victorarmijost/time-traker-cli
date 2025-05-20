package main

import (
	"os"

	"varmijo/time-tracker/tt/app"
	"varmijo/time-tracker/tt/infrastructure/cmd/handlers"
	"varmijo/time-tracker/tt/infrastructure/cmd/repl"
	"varmijo/time-tracker/tt/infrastructure/cmd/repl/myterm"
	"varmijo/time-tracker/tt/infrastructure/config"
	"varmijo/time-tracker/tt/infrastructure/display"
	"varmijo/time-tracker/tt/infrastructure/repositories"
	"varmijo/time-tracker/tt/infrastructure/utils"

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
		logrus.Fatalf("Failed to create SQLite DB: %v", err)
	}

	records := repositories.NewSQLiteRecordRepository(db)

	track := repositories.NewSQLiteTrackRepository(db)

	stats := repositories.NewSQLiteStatsRepository(db)

	app := app.NewApp(cfg, records, track, stats)

	gui := display.NewGUI(app)

	mux := handlers.NewHandlers(app)

	term, closeTerm := myterm.NewTerm()
	defer closeTerm()

	term.PrintTitle("Welcome to Time Tracker CLI tool")

	cmds := repl.NewRepl(app.GetPromptData(), mux, term, "exit")

	go func() {
		cmds.Run()
		gui.Done()
	}()

	gui.Run()
}

// Set up the application logger
func setLogger(slevel string) *os.File {
	file, err := os.OpenFile(utils.GeAppPath(logFile), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logrus.Fatal(err)
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
