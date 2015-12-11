package dashboards

import (
	"sync"

	ui "github.com/gizak/termui"
	TS "github.com/vrecan/FluxDash/timeselect"
	// "log"
)

const (
	KBD_Q     = 1
	KBD_T     = 2
	KBD_Y     = 3
	KBD_SPACE = 4
	KBD_N     = 5
	TIME      = 6
	RESIZE    = 7
)

type Event struct {
	Type int
}

func CommandQ(inputQ chan Event, done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
loop:
	for {
		select {
		case <-done:
			break loop
		case e := <-inputQ:
			if e.Type == KBD_Q {
				ui.StopLoop()
			}
		}
	}
}

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
	m.cDash.Create()
	m.cDash.UpdateAll(m.time)

}

func (m *Monitor) NextDash() {
	m.dashPos++
	if m.dashPos > len(m.Dashes)-1 {
		m.dashPos = 0
	}
	m.cDash = m.Dashes[m.dashPos]
	ui.Clear()
	m.StartDash()
	ui.Render(m.cDash.GetGrid())
}
func (m *Monitor) Close() error {
	ui.StopLoop()
	return nil
}
