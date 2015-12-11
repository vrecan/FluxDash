package dashboards

import (
	// DBC "github.com/influxdb/influxdb/client/v2"
	"fmt"
	ui "github.com/gizak/termui"
	// BC "github.com/vrecan/FluxDash/barchart"
	G "github.com/vrecan/FluxDash/gauge"
	DB "github.com/vrecan/FluxDash/influx"
	MS "github.com/vrecan/FluxDash/multispark"
	// SL "github.com/vrecan/FluxDash/sparkline"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type Carpenter struct {
	JMXMemory *MS.MultiSpark
	Tables    *MS.MultiSpark
	Disk      *G.Gauge
	TimePar   *ui.Par

	Time *TS.TimeSelect
	db   *DB.Influx
	Grid *ui.Grid
}

func NewCarpenter(db *DB.Influx) *Carpenter {
	return &Carpenter{Time: &TS.TimeSelect{}, db: db}
}

func (s *Carpenter) Create() {
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

	JMXMemoryi := MS.MultiSparkInfo{From: `/jmx\.memory.*/`,
		Title:     "JMX Memory",
		Where:     `"service"= 'carpenter'`,
		DataType:  MS.Bytes,
		LineColor: ui.ColorGreen,
	}
	s.JMXMemory = MS.NewMultiSpark(s.db, JMXMemoryi)

	Tablesi := MS.MultiSparkInfo{From: `/(table|cleartables)/`,
		Title:     "Table Inserts",
		Where:     `"service"= 'carpenter'`,
		DataType:  MS.Short,
		LineColor: ui.ColorMagenta,
	}
	s.Tables = MS.NewMultiSpark(s.db, Tablesi)

	// build layout
	grid := ui.NewGrid()
	grid.BgColor = ui.ThemeAttr("bg")
	grid.Width = ui.TermWidth()

	grid.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, s.TimePar),
			ui.NewCol(6, 0, s.Disk.Gauges())),
		ui.NewRow(s.Tables.GetColumns()...),
		ui.NewRow(s.JMXMemory.GetColumns()...))

	// calculate layout
	grid.Align()
	s.Grid = grid
}

func (s *Carpenter) GetGrid() *ui.Grid {
	return s.Grid
}

func (s *Carpenter) UpdateAll(time *TS.TimeSelect) {
	s.Time = time
	dt, di, dr := s.Time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	s.TimePar.Text = displayTimes
	ctime, interval, _ := s.Time.CurTime()
	ui.Render(s.Grid)
	s.JMXMemory.Update(ctime, interval)
	s.Tables.Update(ctime, interval)
	s.Disk.Update(ctime)
	ui.Render(s.Grid)

}
