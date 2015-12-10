package main

import (
	ui "github.com/gizak/termui"
	DBC "github.com/influxdb/influxdb/client/v2"
	// tm "github.com/nsf/termbox-go"
	"fmt"
	DASH "github.com/vrecan/FluxDash/dashboard"
	DB "github.com/vrecan/FluxDash/influx"
	SL "github.com/vrecan/FluxDash/sparkline"
)

func main() {
	d := DASH.NewDashboard("example.json")
	DASH.CreateExampleDash()
	Run(d)
}

func Run(d DASH.Dashboard) {
	c := DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"}
	db, err := DB.NewInflux(c)
	if nil != err {
		panic(err)
	}
	fmt.Println(db)
	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	sparks := SL.NewSparkLinesFromData(db, d.Lines)
	// cpu := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorRed | ui.AttrBold},
	// 	"/system.cpu/", "now() - 15m", db, "CPU", "")
	// cpu.DataType = SL.Percent
	// memFree := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
	// 	"/system.mem.free/", "now() - 15m", db, "MEM Free", "")
	// memFree.DataType = SL.Bytes
	// memCached := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
	// 	"/system.mem.cached/", "now() - 15m", db, "MEM Cached", "")
	// memCached.DataType = SL.Bytes
	// memBuffers := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
	// 	"/system.mem.buffers/", "now() - 15m", db, "MEM Buffers", "")
	// memBuffers.DataType = SL.Bytes
	// gcPause := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
	// 	"/gc.pause.ns/", "now() - 15m", db, "GC Pause Time", "")
	// gcPause.DataType = SL.Time
	// sp1 := SL.NewSparkLines(cpu, memFree, memCached, memBuffers, gcPause)

	// relayIncoming := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
	// 	"/Relay.IncomingMessages/", "now() - 15m", db, "Relay Incomming", `"service"= 'anubis'`)
	// anubis := SL.NewSparkLines(relayIncoming)

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, sparks.Sparks())))
	// calculate layout
	ui.Body.Align()
	sparks.Update()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()

	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		sparks.Update()
		// sp1.Update()
		// anubis.Update()
		// ui.Render(ui.Body)

	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
