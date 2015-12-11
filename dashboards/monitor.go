package dashboards

import (
	ui "github.com/gizak/termui"
	TS "github.com/vrecan/FluxDash/timeselect"
	// "log"
)

type Monitor struct {
	time    *TS.TimeSelect
	Dashes  []Stats
	dashPos int
	cDash   Stats
}

func NewMonitor(s ...Stats) *Monitor {
	return &Monitor{time: &TS.TimeSelect{}, Dashes: s, dashPos: 0}
}

func (m *Monitor) Start() {
	m.run()
}

type Stats interface {
	UpdateAll(*TS.TimeSelect)
	Create()
	GetGrid() *ui.Grid
}

func (m *Monitor) run() {
	counter := uint64(0)

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	time, interval, refresh := m.time.CurTime()

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	//adjust time range
	ui.Handle("/sys/kbd/t", func(ui.Event) {
		time, interval, refresh = m.time.NextTime()
		m.cDash.UpdateAll(m.time)
	})

	ui.Handle("/sys/kbd/y", func(ui.Event) {

		time, interval, refresh = m.time.PrevTime()
		m.cDash.UpdateAll(m.time)
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()

	})
	ui.Handle("/sys/kbd/<space>", func(e ui.Event) {
		m.cDash.UpdateAll(m.time)
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		counter++
		if counter%uint64(refresh) == 0 {
			m.cDash.UpdateAll(m.time)
		}

	})

	ui.Handle("/sys/kbd/n", func(e ui.Event) {
		m.NextDash()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		m.cDash.GetGrid().Width = ui.TermWidth()
		m.cDash.GetGrid().Align()
		ui.Render(m.cDash.GetGrid())
	})

	m.StartDash()
	ui.Loop()
}

func (m *Monitor) StartDash() {
	m.cDash = m.Dashes[m.dashPos]
	if m.cDash.GetGrid() == nil {
		m.cDash.Create()
	}
	m.cDash.UpdateAll(m.time)
	ui.Render(m.cDash.GetGrid())
}

func (m *Monitor) NextDash() {
	m.dashPos++
	if m.dashPos > len(m.Dashes)-1 {
		m.dashPos = 0
	}
	m.cDash = m.Dashes[m.dashPos]
	m.StartDash()
}
func (m *Monitor) Close() error {
	ui.StopLoop()
	return nil
}