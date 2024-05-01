package main

import (
	"fmt"
	"strings"
	"time"
	"varmijo/time-tracker/localStore"
	"varmijo/time-tracker/repl"
	"varmijo/time-tracker/state"
)

func formatDuration(d float64) string {
	h := int(d)
	d -= float64(h)
	m := int(d * 60)

	return fmt.Sprintf("%d:%02d", h, m)
}

func getPrompt(state *state.State) repl.Prompt {
	var (
		wt, ct, pt, tt float64
	)

	return func(pk repl.PromptType) string {
		if pk == repl.FULL_UPDATE {
			wt = localStore.GetTimeByStatus(state.Date, localStore.StatusPending)
			ct = localStore.GetTimeByStatus(state.Date, localStore.StatusCommited)
			pt = localStore.GetTimeByStatus(nil, localStore.StatusPool)
			tt = state.GetTaskTime(nil)
		}

		statusBar := ""

		if wt > 0 {
			statusBar = fmt.Sprintf("[Worked:%s]", formatDuration(wt))
		}

		if ct > 0 {
			statusBar = fmt.Sprintf("%s[Commited:%s]", statusBar, formatDuration(ct))
		}

		if pt > 0 {
			statusBar = fmt.Sprintf("%s[Pool:%s]", statusBar, formatDuration(pt))
		}

		if state.IsWorking() {
			statusBar = fmt.Sprintf("%s[Tracking:%s][%s]", statusBar, formatDuration(tt), getClockEmoji())
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
