package main

import "os/exec"

func calculateTime(t, duration float64) int {
	count := int(t / duration)
	return count
}

func getRemainingTime(count int, t, duration float64) float64 {
	return t - float64(count)*duration
}

func calculatePercentage(t, duration float64) int {
	return int((t / duration) * 100)
}

func pomodoro(t float64) (int, int, string) {
	t = t * 60
	workTime := float64(25)
	breakTime := float64(5)
	longBreakTime := float64(15)

	cicleDuration := 4*workTime + 3*breakTime + longBreakTime
	cicles := calculateTime(t, cicleDuration)
	t = getRemainingTime(cicles, t, cicleDuration)

	pomCount := cicles * 4
	pairDuration := workTime + breakTime
	pairs := calculateTime(t, pairDuration)

	if pairs > 3 {
		pairs = 3
	}

	pomCount += pairs
	t = getRemainingTime(pairs, t, pairDuration)

	if t < workTime {
		return pomCount, calculatePercentage(t, workTime), "w"
	}

	pomCount++
	t -= workTime

	if pomCount%4 == 0 {
		return pomCount, calculatePercentage(t, longBreakTime), "lb"
	}

	return pomCount, calculatePercentage(t, breakTime), "b"
}

func osascriptAlert(msg string) {
	cmd := exec.Command("osascript", "-e", "display notification \""+msg+"\" with title \"TT Pomodoro\" sound name \"Blow\"")
	_ = cmd.Run()
}

func pomodoroAlert(ps string) {
	switch ps {
	case "w":
		osascriptAlert("Time to continue working!")
	case "b":
		osascriptAlert("Take a short break!")
	case "lb":
		osascriptAlert("Long break time!!")
	}
}
