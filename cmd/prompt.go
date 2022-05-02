package main

import (
	"fmt"
	"strings"
	"time"
	"varmijo/time-tracker/localStore"
	"varmijo/time-tracker/repl"
	"varmijo/time-tracker/state"
)

var wt, ct, pt, tt float32

func getPrompt(state *state.State) repl.Prompt {
	return func(pk repl.PromptType) string {
		if pk == repl.FULL_UPDATE {
			wt = localStore.GetTimeByStatus(state.Date, localStore.StatusPending)
			ct = localStore.GetTimeByStatus(state.Date, localStore.StatusCommited)
			pt = localStore.GetTimeByStatus(nil, localStore.StatusPool)
			tt = state.GetTaskTime(nil)
		}

		statusBar := ""

		if wt > 0 {
			statusBar = fmt.Sprintf("[Worked:%.2f]", wt)
		}

		if ct > 0 {
			statusBar = fmt.Sprintf("%s[Commited:%.2f]", statusBar, ct)
		}

		if pt > 0 {
			statusBar = fmt.Sprintf("%s[Pool:%.2f]", statusBar, pt)
		}

		if state.IsWorking() {
			statusBar = fmt.Sprintf("%s[Tracking:%.2f][%s]", statusBar, tt, getClockEmoji())
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
