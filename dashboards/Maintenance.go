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

type Maintenance struct {
	Ind     *MS.MultiSpark
	App     *SL.SparkLines
	Disk    *G.Gauge
	TimePar *ui.Par
	Indices *BC.BarChart
	Shards  *BC.BarChart

	Time *TS.TimeSelect
	db   *DB.Influx
	Grid *ui.Grid
}

func NewMaintenance(db *DB.Influx) *Maintenance {
	return &Maintenance{Time: &TS.TimeSelect{}, db: db}
}

func (s *Maintenance) Create() {
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

	iind := BC.BarChartInfo{From: `/es.indices/`,
		Title:  "Indices",
		Height: 10,
		Where:  `"service"= 'gomaintain'`}
	s.Indices = BC.NewBarChart(s.db, iind)

	ishard := BC.BarChartInfo{From: `/es.*\.shards/`,
		Title:  "ES Shards",
		Height: 10,
		Where:  `"service"= 'gomaintain'`}
	s.Shards = BC.NewBarChart(s.db, ishard)

	// /es.indices/
	gcPause := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		`/es\..*.JVM.gc.time.ms/`, s.db, "GC Pause", `"service"='gomaintain'`)
	gcPause.DataType = SL.Time
	gcCount := SL.NewSparkLine(ui.Sparkline{Height: 1, LineColor: ui.ColorBlue},
		`/es\..*.JVM.gc.count/`, s.db, "GC Count", `"service"='gomaintain'`)
	gcCount.DataType = SL.Short
	s.App = SL.NewSparkLines(gcPause, gcCount)
	s.App.SL.Block.BorderLabel = "App"

	// build layout
	grid := ui.NewGrid()
	grid.BgColor = ui.ThemeAttr("bg")
	grid.Width = ui.TermWidth()

	grid.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, s.TimePar),
			ui.NewCol(6, 0, s.Disk.Gauges())),
		ui.NewRow(s.Indices.GetColumns()...),
		ui.NewRow(s.Shards.GetColumns()...),
		ui.NewRow(ui.NewCol(12, 0, s.App.Sparks())))

	// calculate layout
	grid.Align()
	s.Grid = grid
}

func (s *Maintenance) GetGrid() *ui.Grid {
	return s.Grid
}

func (s *Maintenance) UpdateAll(time *TS.TimeSelect) {
	s.Time = time
	dt, di, dr := s.Time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	s.TimePar.Text = displayTimes
	ctime, interval, _ := s.Time.CurTime()
	ui.Render(s.Grid)
	s.Indices.Update(ctime)
	s.Shards.Update(ctime)
	s.App.Update(ctime, interval)
	s.Disk.Update(ctime)
	ui.Render(s.Grid)

}
