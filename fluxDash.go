package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	DBC "github.com/influxdb/influxdb/client/v2"
	BC "github.com/vrecan/FluxDash/barchart"
	G "github.com/vrecan/FluxDash/gauge"
	DB "github.com/vrecan/FluxDash/influx"
	MS "github.com/vrecan/FluxDash/multispark"
	SL "github.com/vrecan/FluxDash/sparkline"
	TS "github.com/vrecan/FluxDash/timeselect"
)

func main() {

	Run()
}

func Run() {
	counter := uint64(0)
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
	dt, di, dr := time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	_times := ui.NewPar(displayTimes)
	_times.Height = 1
	_times.Border = false
	idisk := G.GaugeInfo{From: `/es\..*.FS.Used.Percent/`,
		Time:  "now() - 15m",
		Title: "Disk Percent Used",
		Where: `"service"= 'gomaintain'`}
	diskUsed := G.NewGauge(ui.ColorCyan, db, idisk)
	iind := BC.BarChartInfo{From: `/es.*\.shards/`,
		Time:  "now() - 1m",
		Title: "ES Shards",
		Where: `"service"= 'gomaintain'`}
	indices := BC.NewBarChart(db, iind)

	dispatchi := MS.MultiSparkInfo{From: `/Dispatch.*/`,
		Time:     "now() - 15m",
		Title:    "Dispatch Info",
		Where:    `"service"= 'godispatch'`,
		DataType: 1,
	}
	dispatch := MS.NewMultiSpark(db, dispatchi)
	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, _times)),
		ui.NewRow(
			ui.NewCol(12, 0, sp1.Sparks())),
		ui.NewRow(
			ui.NewCol(12, 0, anubis.Sparks())),
		ui.NewRow(
			ui.NewCol(12, 0, diskUsed.Gauges())),
		ui.NewRow(diskUsed.GetColumns()...),
		ui.NewRow(indices.GetColumns()...),
		ui.NewRow(dispatch.GetColumns()...),
	)

	// calculate layout
	ui.Body.Align()
	qTime, interval, refresh := time.CurTime()

	updateAll := func() {
		sp1.Update(qTime, interval)
		anubis.Update(qTime, interval)
		diskUsed.Update(qTime)
		indices.Update(qTime)
		dispatch.Update(qTime, interval)
		ui.Render(ui.Body)
	}
	updateAll()

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	//adjust time range
	ui.Handle("/sys/kbd/t", func(ui.Event) {

		qTime, interval, refresh = time.NextTime()
		dt, di, dr = time.DisplayTimes()
		displayTimes = fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
		_times.Text = displayTimes
		ui.Render(ui.Body)
		updateAll()
	})

	ui.Handle("/sys/kbd/y", func(ui.Event) {

		qTime, interval, refresh = time.PrevTime()
		dt, di, dr = time.DisplayTimes()
		displayTimes = fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
		_times.Text = displayTimes
		ui.Render(ui.Body)
		updateAll()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()

	})
	ui.Handle("/sys/kbd/<space>", func(e ui.Event) {
		updateAll()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		counter++
		if counter%uint64(refresh) == 0 {
			updateAll()
		}

	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
