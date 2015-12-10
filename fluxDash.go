package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	DBC "github.com/influxdb/influxdb/client/v2"
	// tm "github.com/nsf/termbox-go"
	DB "github.com/vrecan/FluxDash/influx"
	SL "github.com/vrecan/FluxDash/sparkline"
	TS "github.com/vrecan/FluxDash/timeselect"
)

func main() {

	Run()
}

func Run() {
	time := TS.TimeSelect{}
	c := DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"}
	db, err := DB.NewInflux(c)
	if nil != err {
		panic(err)
	}

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	cpu := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorRed | ui.AttrBold},
		"/system.cpu/", db, "CPU", "")
	cpu.DataType = SL.Percent
	memFree := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.free/", db, "MEM Free", "")
	memFree.DataType = SL.Bytes
	memCached := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.cached/", db, "MEM Cached", "")
	memCached.DataType = SL.Bytes
	memBuffers := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.buffers/", db, "MEM Buffers", "")
	memBuffers.DataType = SL.Bytes
	gcPause := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/gc.pause.ns/", db, "GC Pause Time", "")
	gcPause.DataType = SL.Time
	sp1 := SL.NewSparkLines(cpu, memFree, memCached, memBuffers, gcPause)
	sp1.SL.Block.BorderLabel = "System"

	relayIncoming := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/Relay.IncomingMessages/", db, "Relay Incomming", `"service"= 'anubis'`)
	anubis := SL.NewSparkLines(relayIncoming)
	anubis.SL.Block.BorderLabel = "Anubis"
	dt, di := time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s", dt, di)
	_times := ui.NewPar(displayTimes)

	_times.Height = 1
	_times.Border = false
	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, _times)),
		ui.NewRow(
			ui.NewCol(12, 0, sp1.Sparks())),
		ui.NewRow(
			ui.NewCol(12, 0, anubis.Sparks())))

	// calculate layout
	ui.Body.Align()

	qTime, interval := time.CurTime()

	sp1.Update(qTime, interval)
	anubis.Update(qTime, interval)
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	//adjust time range
	ui.Handle("/sys/kbd/t", func(ui.Event) {

		qTime, interval = time.NextTime()
		dt, di = time.DisplayTimes()
		displayTimes = fmt.Sprintf("Time: %s Interval: %s", dt, di)
		_times.Text = displayTimes
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/y", func(ui.Event) {

		qTime, interval = time.PrevTime()
		dt, di = time.DisplayTimes()
		displayTimes = fmt.Sprintf("Time: %s Interval: %s", dt, di)
		_times.Text = displayTimes
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()

	})
	ui.Handle("/timer/1s", func(e ui.Event) {

		sp1.Update(qTime, interval)
		anubis.Update(qTime, interval)
		ui.Render(ui.Body)

	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
