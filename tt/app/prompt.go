package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"varmijo/time-tracker/tt/domain"
)

func formatDuration(d float64) string {
	h := int(d)
	d -= float64(h)
	m := int(d * 60)

	return fmt.Sprintf("%d:%02d", h, m)
}

func (kern *App) GetPrompt() domain.Prompt {
	var (
		wt, ct, pt, tt float64
	)

	return func(pk domain.PromptType) string {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		date := kern.date.Get()

		if pk == domain.FULL_UPDATE {
			wt = domain.Must(kern.records.GetHoursByDateStatus(ctx, date, domain.StatusPending))
			ct = domain.Must(kern.records.GetHoursByDateStatus(ctx, date, domain.StatusCommited))
			pt = domain.Must(kern.records.GetHoursByStatus(ctx, domain.StatusPool))
			tt = domain.Must(kern.track.GetHours(ctx))
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

		if kern.track.IsWorking(ctx) {
			statusBar = fmt.Sprintf("%s[Tracking:%s][%s]", statusBar, formatDuration(tt), getClockEmoji())
		}

		if !kern.date.IsToday() {
			statusBar = fmt.Sprintf("%s[%s]", statusBar, kern.date.Get().Format("06-01-02"))
		}

		if kern.pomodoro.Has() {
			pomState := kern.pomodoro.GetState()
			pomProg := kern.pomodoro.GetProgress()

			statusBar = fmt.Sprintf("%s[%s:%d%%]", statusBar, pomState, pomProg)

			if pomState == "b" {
				statusBar = fmt.Sprintf("%s[%0.f]", statusBar, kern.pomodoro.GetBreakTime())
			}

			if pomProg >= 100 {
				statusBar = fmt.Sprintf("%s%s", statusBar, getAlertEmoji())
			}
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

func getAlertEmoji() string {
	alerts := []string{"{!}", "{ }"}

	n := int64(len(alerts))

	return alerts[time.Now().Unix()%n]
}
