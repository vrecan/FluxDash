package barchart

import (
	"fmt"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/query"
)

type BarChartInfo struct {
	From   string
	Time   string
	Title  string
	Where  string
	Height int
}

type BarChart struct {
	I  BarChartInfo
	C  *ui.BarChart
	L  *ui.List
	db *DB.Influx
}

func NewBarChart(db *DB.Influx, info BarChartInfo) *BarChart {
	barchart := ui.NewBarChart()
	list := ui.NewList()
	g := &BarChart{C: barchart, L: list, db: db, I: info}
	g.C.DataLabels = make([]string, 0)
	g.C.Height = info.Height
	g.L.Height = info.Height
	g.L.BorderLabel = info.Title
	// g.L.ItemFgColor = ui.ColorBlack
	// g.L.ItemBgColor = ui.ColorWhite
	return g
}

func (s *BarChart) BarCharts() *ui.BarChart {
	return s.C
}
func (s *BarChart) Labels() *ui.List {
	return s.L
}
func (s *BarChart) Update(time string) {
	s.SetData(time)
	s.SetTitle()
}

func (s *BarChart) SetData(time string) {
	meanTotal, labels := query.GetDataForBar(s.db, query.Build("mean(value)", s.I.From, s.I.Where, time, ""))
	s.C.Data = meanTotal
	series := make([]string, len(labels))
	items := make([]string, len(labels))
	i := 0
	for _, v := range labels {
		series[i] = v[0]
		items[i] = fmt.Sprintf("%s: %s", v[0], v[1])
		i++
	}
	s.C.DataLabels = series
	s.L.Items = items
}
func (s *BarChart) GetColumns() []*ui.Row {
	return []*ui.Row{ui.NewCol(6, 0, s.BarCharts()), ui.NewCol(6, 0, s.Labels())}
}

func (s *BarChart) SetTitle() {
	s.C.BorderLabel = fmt.Sprintf("%s", s.I.Title)
}
