package dashboards

import (
	// DBC "github.com/influxdb/influxdb/client/v2"
	"fmt"
	ui "github.com/gizak/termui"
	// BC "github.com/vrecan/FluxDash/barchart"
	G "github.com/vrecan/FluxDash/gauge"
	DB "github.com/vrecan/FluxDash/influx"
	MS "github.com/vrecan/FluxDash/multispark"
	SL "github.com/vrecan/FluxDash/sparkline"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type GoDispatch struct {
	Dispatch *MS.MultiSpark
	App      *SL.SparkLines
	Disk     *G.Gauge
	TimePar  *ui.Par

	Time *TS.TimeSelect
	db   *DB.Influx
	Grid *ui.Grid
}

func NewGoDispatch(db *DB.Influx) *GoDispatch {
	return &GoDispatch{Time: &TS.TimeSelect{}, db: db}
}

func (s *GoDispatch) Create() {
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

	Dispatchi := MS.MultiSparkInfo{From: `/Dispatch.*/`,
		Title:     "Dispatch",
		Where:     `"service"= 'godispatch'`,
		DataType:  MS.Short,
		LineColor: ui.ColorMagenta,
	}

	s.Dispatch = MS.NewMultiSpark(s.db, Dispatchi)

	cpu := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorRed},
		"/cpu.percent/", s.db, "CPU", `"service"='godispatch'`)
	cpu.DataType = SL.Percent

	gcPause := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		"/gc.pause.ns/", s.db, "GC Pause", `"service"='godispatch'`)
	gcPause.DataType = SL.Time
	s.App = SL.NewSparkLines(cpu, gcPause)
	s.App.SL.Block.BorderLabel = "App"

	// build layout
	grid := ui.NewGrid()
	grid.BgColor = ui.ThemeAttr("bg")
	grid.Width = ui.TermWidth()

	grid.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, s.TimePar),
			ui.NewCol(6, 0, s.Disk.Gauges())),
		ui.NewRow(s.Dispatch.GetColumns()...),
		ui.NewRow(ui.NewCol(12, 0, s.App.Sparks())))

	// calculate layout
	grid.Align()
	s.Grid = grid
}

func (s *GoDispatch) GetGrid() *ui.Grid {
	return s.Grid
}

func (s *GoDispatch) UpdateAll(time *TS.TimeSelect) {
	s.Time = time
	dt, di, dr := s.Time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	s.TimePar.Text = displayTimes
	ctime, interval, _ := s.Time.CurTime()
	ui.Render(s.Grid)
	s.App.Update(ctime, interval)
	s.Dispatch.Update(ctime, interval)
	s.Disk.Update(ctime)
	ui.Render(s.Grid)

}
