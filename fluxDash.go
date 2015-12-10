package main

import (
	"io"
	"os"
	SYS "syscall"

	ui "github.com/gizak/termui"
	DBC "github.com/influxdb/influxdb/client/v2"
	DEATH "github.com/vrecan/death"
	// tm "github.com/nsf/termbox-go"
	DB "github.com/vrecan/FluxDash/influx"
	SL "github.com/vrecan/FluxDash/sparkline"
)

func main() {
	var goRoutines []io.Closer
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	theUi := closeUI{}
	go theUi.Start()

	goRoutines = append(goRoutines, closeUI{})
	death.WaitForDeath(goRoutines...)

	// fmt.Println("Exiting...")

}

type closeUI struct{}

func (theUI closeUI) Start() {
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
		"/system.cpu/", "now() - 15m", db, "CPU", "")
	cpu.DataType = SL.Percent
	memFree := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.free/", "now() - 15m", db, "MEM Free", "")
	memFree.DataType = SL.Bytes
	memCached := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.cached/", "now() - 15m", db, "MEM Cached", "")
	memCached.DataType = SL.Bytes
	memBuffers := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/system.mem.buffers/", "now() - 15m", db, "MEM Buffers", "")
	memBuffers.DataType = SL.Bytes
	gcPause := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/gc.pause.ns/", "now() - 15m", db, "GC Pause Time", "")
	gcPause.DataType = SL.Time
	sp1 := SL.NewSparkLines(cpu, memFree, memCached, memBuffers, gcPause)

	relayIncoming := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue | ui.AttrBold},
		"/Relay.IncomingMessages/", "now() - 15m", db, "Relay Incomming", `"service"= 'anubis'`)
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
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
		ui.Close()
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(os.Interrupt)
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
}
func (c closeUI) Close() error {
	ui.StopLoop()
	return nil
}
