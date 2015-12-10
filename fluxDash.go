package main

import (
	ui "github.com/gizak/termui"
	DBC "github.com/influxdb/influxdb/client/v2"
	// tm "github.com/nsf/termbox-go"
	DB "github.com/vrecan/FluxDash/influx"
	SL "github.com/vrecan/FluxDash/sparkline"
)

func main() {
	c := DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"}
	db, err := DB.NewInflux(c)
	if nil != err {
		panic(err)
	}
	// fmt.Println(db)

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	cpu := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorRed | ui.AttrBold},
		"/system.cpu/", "now() - 15m", db, "CPU")
	cpu.DataType = SL.Percent
	memFree := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.free/", "now() - 15m", db, "MEM Free")
	memFree.DataType = SL.Bytes
	gcPause := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/gc.pause.ns/", "now() - 15m", db, "GC Pause Time")
	sp1 := SL.NewSparkLines(cpu, memFree, gcPause)

	relayIncoming := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/Relay.IncomingMessages/", "now() - 15m", db, "Relay Incomming")
	anubis := SL.NewSparkLines(relayIncoming)

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, sp1.Sparks())),
		ui.NewRow(
			ui.NewCol(12, 0, anubis.Sparks())))

	// calculate layout
	ui.Body.Align()
	sp1.Update()
	anubis.Update()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {

		sp1.Update()
		anubis.Update()
		ui.Render(ui.Body)

	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()

	// fmt.Println("Exiting...")

}
