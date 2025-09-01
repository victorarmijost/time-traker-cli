package app

import (
	"context"
	"sync"
	"time"

	"varmijo/time-tracker/tt/domain"
)

type promptData struct {
	app                *App
	wt, ct, pt, tt, dt float64
	sync.RWMutex
}

var promptDataInstance *promptData
var promptDataOnce sync.Once

func (p *promptData) RefreshData() {
	p.Lock()
	defer p.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	date := p.app.date.Get()

	p.wt = domain.Must(p.app.stats.GetHoursByDateStatus(ctx, date, domain.StatusPending))
	p.ct = domain.Must(p.app.stats.GetHoursByDateStatus(ctx, date, domain.StatusCommitted))
	p.pt = domain.Must(p.app.stats.GetHoursByStatus(ctx, domain.StatusPool))
	p.tt = domain.Must(p.app.stats.GetTrackedHours(ctx))
	p.dt = domain.Must(p.app.stats.GetDebt(ctx, p.app.config.GetWorkTime())).Total()
}

func (p *promptData) Wt() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.wt
}
func (p *promptData) Ct() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.ct
}
func (p *promptData) Pt() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.pt
}
func (p *promptData) Tt() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.tt
}
func (p *promptData) Dt() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.dt
}
func (p *promptData) IsWorking() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return p.app.track.IsWorking(ctx)
}
func (p *promptData) IsToday() bool {
	return p.app.date.IsToday()
}

func (p *promptData) GetDate() time.Time {
	return p.app.date.Get()
}

func (p *promptData) keepRefreshing() {
	for range time.Tick(60 * time.Second) {
		p.RefreshData()
	}
}

func (kern *App) GetPromptData() domain.PromptData {
	promptDataOnce.Do(func() {
		promptDataInstance = &promptData{
			app: kern,
		}

		promptDataInstance.RefreshData()

		go promptDataInstance.keepRefreshing()
	})

	return promptDataInstance
}
