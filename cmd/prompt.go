package main

import (
	"fmt"
	"strings"
	"time"
	"varmijo/time-tracker/localStore"
	"varmijo/time-tracker/repl"
	"varmijo/time-tracker/state"
)

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
			statusBar = fmt.Sprintf("%s[Tracking:%.2f][%s]", statusBar, state.GetTaskTime(nil), getClockEmoji())
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

const clocksEmojis = `. '`

func getClockEmoji() string {
	clocks := strings.Split(clocksEmojis, " ")

	n := int64(len(clocks))

	return strings.TrimSpace(clocks[time.Now().Unix()%n])
}
