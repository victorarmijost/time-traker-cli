package display

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
	"varmijo/time-tracker/tt/app"
	"varmijo/time-tracker/tt/domain"

	"github.com/getlantern/systray"
)

type runningStatus int

const (
	Focus runningStatus = iota
	Working
	Idle
)

type GUI struct {
	done          chan struct{}
	app           *app.App
	propmptData   domain.PromptData
	focusTimer    *time.Time
	pomodoroTimer *time.Time
	recordingMenu *systray.MenuItem
	focusMenu     *systray.MenuItem
	driftMenu     *systray.MenuItem
	status        *StatusHandler[runningStatus]
	sync.Mutex
}

func NewGUI(app *app.App) *GUI {
	g := &GUI{
		done:        make(chan struct{}),
		app:         app,
		propmptData: app.GetPromptData(),
		status:      NewStatusHandler[runningStatus](),
	}

	g.registerHandlers()

	return g
}

func (g *GUI) Done() {
	g.Lock()
	defer g.Unlock()

	if g.done != nil {
		close(g.done)
		g.done = nil
	}
}

func (g *GUI) Run() {
	systray.Run(g.onReady, g.onExit)
}

func (g *GUI) titleText() string {
	if !g.propmptData.IsWorking() {
		return "💤"
	}

	if remaining := g.getFocusTime(); remaining > 0 {
		return fmt.Sprintf("%s %s", "🔥", domain.FormatDuration(remaining))
	}

	if remaining := g.getPomodoroTime(); remaining > 0 {
		return fmt.Sprintf("%s %s", "🍅", domain.FormatDuration(remaining))
	}

	return fmt.Sprintf("%s %s", getClockEmoji(), domain.FormatDuration(g.propmptData.Tt()))
}

func (g *GUI) tooltipText() string {
	statusBar := ""

	if g.propmptData.Dt() > 0 {
		statusBar = fmt.Sprintf("💳 %s", domain.FormatDuration(g.propmptData.Dt()))
	} else if g.propmptData.Dt() < 0 {
		statusBar = fmt.Sprintf("🪣 %s", domain.FormatDuration(-g.propmptData.Dt()))
	}

	if g.propmptData.Wt() > 0 {
		statusBar = fmt.Sprintf("%s\n🔨 %s", statusBar, domain.FormatDuration(g.propmptData.Wt()))
	}

	if !g.propmptData.IsToday() {
		statusBar = fmt.Sprintf("%s\n📅 %s", statusBar, g.propmptData.GetDate().Format("02/Jan/06"))
	}

	return statusBar
}

func (g *GUI) getFocusTime() float64 {
	if g.focusTimer == nil {
		return 0
	}

	left := time.Until(*g.focusTimer).Minutes()

	if left < 0 {
		return 0
	}

	return left
}

func (g *GUI) getPomodoroTime() float64 {
	if g.pomodoroTimer == nil {
		return 0
	}
	left := time.Until(*g.pomodoroTimer).Minutes()

	if left < 0 {
		return 0
	}

	return left
}

func (g *GUI) areTimersDone() bool {
	if g.getFocusTime() == 0 && g.getPomodoroTime() == 0 {
		g.focusTimer = nil
		g.pomodoroTimer = nil
		return true
	}

	return false
}

func (g *GUI) updateTitle(forceRefresh bool) {
	if forceRefresh {
		g.propmptData.RefreshData()
	}

	systray.SetTitle(g.titleText())
	systray.SetTooltip(g.tooltipText())
}

func (g *GUI) startRecordingWithTimers(from, to runningStatus) error {
	err := g.startTimers(from, to)
	if err != nil {
		return err
	}

	if from == Idle {
		return g.startRecording(from, to)
	}

	return nil
}

func (g *GUI) startRecording(from, to runningStatus) error {
	g.stardRecord()
	g.setRecodingTitle(to)
	return nil
}

func (g *GUI) stopRecording(from, to runningStatus) error {
	g.stopRecord()
	g.setRecodingTitle(to)

	return g.clearTimers(from, to)
}

func (g *GUI) clearTimers(from, to runningStatus) error {
	g.focusTimer = nil
	g.pomodoroTimer = nil
	return nil
}

func (g *GUI) startTimers(_, _ runningStatus) error {
	focusTime := time.Now().Add(2 * time.Minute)
	pomodoroTime := time.Now().Add(25 * time.Minute)
	g.focusTimer = &focusTime
	g.pomodoroTimer = &pomodoroTime

	return nil
}

func (g *GUI) registerHandlers() {
	g.status.Register(Idle, Focus, g.startRecordingWithTimers)
	g.status.Register(Idle, Working, g.startRecording)
	g.status.Register(Working, Idle, g.stopRecording)
	g.status.Register(Working, Focus, g.startTimers)
	g.status.Register(Focus, Idle, g.stopRecording)
	g.status.Register(Focus, Focus, g.startRecordingWithTimers)
	g.status.Register(Focus, Working, g.clearTimers)
}

func (g *GUI) updateStatus(status runningStatus) {
	g.status.UpdateStatus(status)
}

func (g *GUI) updateStatusIfNewAllowed(newStatus runningStatus, newAllowed ...runningStatus) {
	if slices.Contains(newAllowed, newStatus) {
		g.updateStatus(newStatus)
	}
}

func (g *GUI) updateStatusIfCurrentAllowed(newStatus runningStatus, currentAllowed ...runningStatus) {
	if slices.Contains(currentAllowed, g.getCurrentStatus()) {
		g.updateStatus(newStatus)
	}
}

func (g *GUI) getCurrentStatus() runningStatus {
	if !g.propmptData.IsWorking() {
		return Idle
	}

	if !g.areTimersDone() {
		return Focus
	}

	return Working
}

func (g *GUI) toogleStatus(currentStatus runningStatus) runningStatus {
	if currentStatus == Idle {
		return Working
	}

	return Idle
}

func (g *GUI) setRecodingTitle(status runningStatus) {
	switch status {
	case Idle:
		g.recordingMenu.SetTitle("Start")
		g.recordingMenu.SetTooltip("Start recording")
	case Focus:
		g.recordingMenu.SetTitle("Stop")
		g.recordingMenu.SetTooltip("Stop recording")
	case Working:
		g.recordingMenu.SetTitle("Stop")
		g.recordingMenu.SetTooltip("Stop recording")
	}
}

func (g *GUI) onReady() {
	g.recordingMenu = systray.AddMenuItem("", "")
	g.setRecodingTitle(g.getCurrentStatus())

	systray.AddSeparator()

	g.focusMenu = systray.AddMenuItem("Focus", "Start focus timer")
	g.driftMenu = systray.AddMenuItem("Drift", "Stop focus timer")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	// Goroutine for the timer
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				g.updateStatusIfNewAllowed(g.getCurrentStatus(), Idle, Working)
				g.updateTitle(false)
			case <-g.recordingMenu.ClickedCh:
				status := g.getCurrentStatus()
				g.updateStatus(g.toogleStatus(status))
				g.updateTitle(true)
			case <-g.focusMenu.ClickedCh:
				g.updateStatus(Focus)
				g.updateTitle(true)
			case <-g.driftMenu.ClickedCh:
				g.updateStatusIfCurrentAllowed(Working, Focus)
				g.updateTitle(false)
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-g.done:
				systray.Quit()
				return
			}
		}
	}()
}

func (g *GUI) onExit() {
	// Cleanup actions
}

func (g *GUI) stardRecord() {
	if g.getCurrentStatus() == Idle {
		_ = g.app.StartRecord(context.Background())
	}
}

func (g *GUI) stopRecord() {
	if g.getCurrentStatus() != Idle {
		_, _ = g.app.StopRecord(context.Background())
	}
}

func getClockEmoji() string {
	flames := []string{"⌛", "⏳"}

	n := int64(len(flames))

	return strings.TrimSpace(flames[time.Now().Unix()%n])
}
