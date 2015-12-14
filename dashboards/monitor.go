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
	Type    int
	Time    *TS.TimeSelect
	Dash    *Stats
	Monitor *Monitor
}

func CommandQ(inputQ chan Event, done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	counter := uint64(0)
loop:
	for {
		select {
		case <-done:
			break loop
		case e := <-inputQ:
			if e.Type == KBD_Q {
				ui.StopLoop()
			} else if e.Type == KBD_T {
				e.Time.NextTime()
				(*e.Dash).UpdateAll(e.Time)
			} else if e.Type == KBD_Y {
				e.Time.PrevTime()
				(*e.Dash).UpdateAll(e.Time)
			} else if e.Type == KBD_SPACE {
				(*e.Dash).UpdateAll(e.Time)
			} else if e.Type == TIME {
				counter++
				_, _, refresh := e.Time.CurTime()
				if counter%uint64(refresh) == 0 {
					(*e.Dash).UpdateAll(e.Time)
				}
			} else if e.Type == KBD_N {
				e.Monitor.NextDash()
			} else if e.Type == RESIZE {
				(*e.Dash).GetGrid().Width = ui.TermWidth()
				(*e.Dash).GetGrid().Align()
				ui.Render((*e.Dash).GetGrid())
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

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	m.time.CurTime()

	inputQ := make(chan Event, 100)
	done := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go CommandQ(inputQ, done, wg)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		inputQ <- Event{Type: KBD_Q}
	})
	//adjust time range
	ui.Handle("/sys/kbd/t", func(ui.Event) {
		inputQ <- Event{Type: KBD_T, Time: m.time, Dash: &m.cDash}

	})

	ui.Handle("/sys/kbd/y", func(ui.Event) {
		inputQ <- Event{Type: KBD_Y, Time: m.time, Dash: &m.cDash}
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		inputQ <- Event{Type: KBD_Q}
	})
	ui.Handle("/sys/kbd/<space>", func(e ui.Event) {
		inputQ <- Event{Type: KBD_SPACE, Time: m.time, Dash: &m.cDash}
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		inputQ <- Event{Type: TIME, Time: m.time, Dash: &m.cDash}
	})

	ui.Handle("/sys/kbd/n", func(e ui.Event) {
		inputQ <- Event{Type: KBD_N, Monitor: m}
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		inputQ <- Event{Type: RESIZE, Dash: &m.cDash}
	})

	m.StartDash()
	ui.Loop()
	done <- true
	wg.Wait()
}

func (m *Monitor) StartDash() {
	m.cDash = m.Dashes[m.dashPos]
	if m.cDash.GetGrid() == nil {
		m.cDash.Create()
	}
	m.cDash.UpdateAll(m.time)
	ui.Clear()
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
