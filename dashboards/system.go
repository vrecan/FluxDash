package dashboards

import (
	// DBC "github.com/influxdb/influxdb/client/v2"
	"fmt"
	ui "github.com/gizak/termui"
	BC "github.com/vrecan/FluxDash/barchart"
	G "github.com/vrecan/FluxDash/gauge"
	DB "github.com/vrecan/FluxDash/influx"
	MS "github.com/vrecan/FluxDash/multispark"
	SL "github.com/vrecan/FluxDash/sparkline"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type System struct {
	Monitor  *SL.SparkLines
	Anubis   *SL.SparkLines
	Disk     *G.Gauge
	Indices  *BC.BarChart
	Dispatch *MS.MultiSpark
	Time     *TS.TimeSelect
	TimePar  *ui.Par
	db       *DB.Influx
	Grid     *ui.Grid
}

func NewSystem(db *DB.Influx) *System {
	return &System{Time: &TS.TimeSelect{}, db: db}
}

func (s *System) Create() {
	cpu := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorRed},
		"/system.cpu/", s.db, "CPU", "")
	cpu.DataType = SL.Percent
	memFree := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		"/system.mem.free/", s.db, "MEM Free", "")
	memFree.DataType = SL.Bytes
	memCached := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		"/system.mem.cached/", s.db, "MEM Cached", "")
	memCached.DataType = SL.Bytes
	memBuffers := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		"/system.mem.buffers/", s.db, "MEM Buffers", "")
	memBuffers.DataType = SL.Bytes
	s.Monitor = SL.NewSparkLines(cpu, memFree, memCached, memBuffers)
	s.Monitor.SL.Block.BorderLabel = "System"
	relayIncoming := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		"/Relay.IncomingMessages/", s.db, "Relay Incomming", `"service"= 'anubis'`)
	s.Anubis = SL.NewSparkLines(relayIncoming)
	s.Anubis.SL.Block.BorderLabel = "Anubis"
	dt, di, dr := s.Time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	s.TimePar = ui.NewPar(displayTimes)
	s.TimePar.Height = 3
	s.TimePar.Border = true
	idisk := G.GaugeInfo{From: `/es\..*.FS.Used.Percent/`,
		Title:  "Disk Percent Used",
		Height: 3,
		Border: true,
		Where:  `"service"= 'gomaintain'`}
	s.Disk = G.NewGauge(ui.ColorCyan, s.db, idisk)
	iind := BC.BarChartInfo{From: `/es.*\.shards/`,
		Title:  "ES Shards",
		Height: 10,
		Where:  `"service"= 'gomaintain'`}
	s.Indices = BC.NewBarChart(s.db, iind)

	dispatchi := MS.MultiSparkInfo{From: `/Dispatch.*/`,
		Title:    "Dispatch",
		Where:    `"service"= 'godispatch'`,
		DataType: MS.Short,
	}
	s.Dispatch = MS.NewMultiSpark(s.db, dispatchi)
	// build layout
	grid := ui.NewGrid()
	grid.BgColor = ui.ThemeAttr("bg")
	grid.Width = ui.TermWidth()

	grid.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, s.TimePar),
			ui.NewCol(6, 0, s.Disk.Gauges())),
		ui.NewRow(
			ui.NewCol(12, 0, s.Monitor.Sparks())),
		ui.NewRow(
			ui.NewCol(12, 0, s.Anubis.Sparks())),
		ui.NewRow(s.Indices.GetColumns()...),
		ui.NewRow(s.Dispatch.GetColumns()...),
	)
	// calculate layout
	grid.Align()
	s.Grid = grid
}

func (s *System) GetGrid() *ui.Grid {
	return s.Grid
}

func (s *System) UpdateAll(time *TS.TimeSelect) {
	s.Time = time
	dt, di, dr := s.Time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	s.TimePar.Text = displayTimes
	ctime, interval, _ := s.Time.CurTime()
	ui.Render(s.Grid)
	s.Monitor.Update(ctime, interval)
	s.Anubis.Update(ctime, interval)
	s.Disk.Update(ctime)
	s.Indices.Update(ctime)
	s.Dispatch.Update(ctime, interval)
	ui.Render(s.Grid)
}
