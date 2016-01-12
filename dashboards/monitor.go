package dashboards

import (
	// "fmt"
	ui "github.com/gizak/termui"
	TS "github.com/vrecan/FluxDash/timeselect"
	"sync"
	"time"
)

//Keybaord consts
const (
	KBD_Q     = 1
	KBD_T     = 2
	KBD_Y     = 3
	KBD_SPACE = 4
	KBD_N     = 5
	TIME      = 6
	RESIZE    = 7
)

//Event are the events sent to the commandQ.
type Event struct {
	Type    int
	Time    *TS.TimeSelect
	Dash    Stats
	Monitor *Monitor
}

//CommandQ is our main loop for processing input
func CommandQ(inputQ <-chan interface{}, done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	counter := uint64(0)
	timeDebounceChan := make(chan interface{}, 1000)
	TimeChangeChan := make(chan interface{}, 1)
	go DebounceChan(timeDebounceChan, 13*time.Millisecond, TimeChangeChan)
loop:
	for {
		select {
		case <-done:
			break loop

		case event := <-TimeChangeChan:
			e := event.(Event)
			if e.Type == KBD_T {
				e.Time.NextTime()
				(e.Dash).UpdateAll(e.Time)
			} else if e.Type == KBD_Y {
				e.Time.PrevTime()
				(e.Dash).UpdateAll(e.Time)
			} else if e.Type == KBD_SPACE {
				(e.Dash).UpdateAll(e.Time)
			} else if e.Type == TIME {
				counter++
				_, _, refresh := e.Time.CurTime()
				if counter%uint64(refresh) == 0 {
					(e.Dash).UpdateAll(e.Time)
				}
			} else if e.Type == KBD_N {
				e.Monitor.NextDash()
			}

		case event := <-inputQ:
			e := event.(Event)
			if e.Type == KBD_T {
				timeDebounceChan <- e
			} else if e.Type == KBD_Y {
				timeDebounceChan <- e
			} else if e.Type == KBD_SPACE {
				timeDebounceChan <- e
			} else if e.Type == TIME {
				timeDebounceChan <- e
			} else if e.Type == KBD_N {
				timeDebounceChan <- e
			} else if e.Type == RESIZE {
				(e.Dash).GetGrid().Width = ui.TermWidth()
				(e.Dash).GetGrid().Align()
				ui.Clear()
				ui.Render((e.Dash).GetGrid())
			}
		}
	}
}

//DebounceChan is a simple debouncer using channels that will grab the last message after the timeout.
func DebounceChan(input chan interface{}, wait time.Duration, res chan interface{}) {
	var dRes interface{}
	timer := time.NewTimer(wait)
	timer.Stop()
	var started bool
	for {
		select {
		case d, ok := <-input:
			if !ok {
				break
			}
			if !started {
				timer.Reset(wait)
				started = true
			}
			dRes = d
		case <-timer.C:
			timer.Stop()
			started = false
			res <- dRes
		}
	}
	timer.Stop()
}

//Monitor is the main struct to monitor the dashboards.
type Monitor struct {
	time    *TS.TimeSelect
	Dashes  []Stats
	dashPos int
	cDash   Stats
}

//NewMonitor creates a new monitor struct that can display dashboards.
func NewMonitor(s ...Stats) *Monitor {
	return &Monitor{time: &TS.TimeSelect{}, Dashes: s, dashPos: 0}
}

//Start the monitor.
func (m *Monitor) Start() {
	m.run()
}

//Stats is an interface for all the widgets
type Stats interface {
	UpdateAll(*TS.TimeSelect)
	Create()
	GetGrid() *ui.Grid
}

//run is the main loop for a monitor.
func (m *Monitor) run() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	m.time.CurTime()

	inputQ := make(chan interface{}, 1000)

	done := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go CommandQ(inputQ, done, wg)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	//adjust time range
	ui.Handle("/sys/kbd/t", func(ui.Event) {
		inputQ <- Event{Type: KBD_T, Time: m.time, Dash: m.cDash}

	})

	ui.Handle("/sys/kbd/y", func(ui.Event) {
		inputQ <- Event{Type: KBD_Y, Time: m.time, Dash: m.cDash}
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		inputQ <- Event{Type: KBD_Q}
	})
	ui.Handle("/sys/kbd/<space>", func(e ui.Event) {
		inputQ <- Event{Type: KBD_SPACE, Time: m.time, Dash: m.cDash}
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		inputQ <- Event{Type: TIME, Time: m.time, Dash: m.cDash}
	})

	ui.Handle("/sys/kbd/n", func(e ui.Event) {
		inputQ <- Event{Type: KBD_N, Monitor: m}
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		inputQ <- Event{Type: RESIZE, Dash: m.cDash}
	})

	m.StartDash()
	ui.Loop()
	done <- true
	wg.Wait()
}

//StartDash enables a dashboard, if it' hasn't been created it creates a new one.
func (m *Monitor) StartDash() {
	m.cDash = m.Dashes[m.dashPos]
	if m.cDash.GetGrid() == nil {
		m.cDash.Create()
	}
	m.cDash.UpdateAll(m.time)
	ui.Clear()
	ui.Render(m.cDash.GetGrid())
}

//NextDash moves to the next dashboard.
func (m *Monitor) NextDash() {
	m.dashPos++
	if m.dashPos > len(m.Dashes)-1 {
		m.dashPos = 0
	}
	m.cDash = m.Dashes[m.dashPos]
	m.StartDash()
}

//Close stop's the main loop and exits the monitor.
func (m *Monitor) Close() error {
	ui.StopLoop()
	return nil
}
