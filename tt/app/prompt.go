package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"varmijo/time-tracker/tt/domain"
)

func (kern *App) GetPrompt() domain.Prompt {
	var (
		wt, ct, pt, tt, dt float64
	)

	return func(pk domain.PromptType) string {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		date := kern.date.Get()

		if pk == domain.FULL_UPDATE {
			wt = domain.Must(kern.stats.GetHoursByDateStatus(ctx, date, domain.StatusPending))
			ct = domain.Must(kern.stats.GetHoursByDateStatus(ctx, date, domain.StatusCommited))
			pt = domain.Must(kern.stats.GetHoursByStatus(ctx, domain.StatusPool))
			tt = domain.Must(kern.stats.GetTrackedHours(ctx))
			dt = domain.Must(kern.stats.GetDebt(ctx, kern.config.GetWorkTime())).Total()
		}

		statusBar := ""

		if dt > 0 {
			statusBar = fmt.Sprintf("[Debt:%s]", domain.FormatDuration(dt))
		}

		if wt > 0 {
			statusBar = fmt.Sprintf("%s[Worked:%s]", statusBar, domain.FormatDuration(wt))
		}

		if ct > 0 {
			statusBar = fmt.Sprintf("%s[Commited:%s]", statusBar, domain.FormatDuration(ct))
		}

		if pt > 0 {
			statusBar = fmt.Sprintf("%s[Pool:%s]", statusBar, domain.FormatDuration(pt))
		}

		if kern.track.IsWorking(ctx) {
			statusBar = fmt.Sprintf("%s[Rec:%s][%s]", statusBar, domain.FormatDuration(tt), getClockEmoji())
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
