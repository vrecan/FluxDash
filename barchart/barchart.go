package barchart

import (
	"fmt"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/query"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type BarChart struct {
	From        string       `json:"from"`
	BorderLabel string       `json:"borderlabel"`
	Where       string       `json:"where"`
	Height      int          `json:"height"`
	GroupBy     string       `json:"groupby,omitempty"`
	C           *ui.BarChart `json:"-"`
	L           *ui.List     `json:"-"`
	db          *DB.Influx   `json:"-"`
}

func NewBarChart(db *DB.Influx, bc *BarChart) *BarChart {
	bc.C = ui.NewBarChart()
	bc.L = ui.NewList()
	bc.db = db
	bc.C.DataLabels = make([]string, 0)
	bc.C.Height = bc.Height
	bc.L.Height = bc.Height
	bc.L.BorderLabel = bc.BorderLabel
	return bc
}

func (s *BarChart) BarCharts() *ui.BarChart {
	return s.C
}
func (s *BarChart) Labels() *ui.List {
	return s.L
}
func (s *BarChart) Update(ts *TS.TimeSelect) {
	time, _, _ := ts.CurTime()
	s.SetData(time)
	s.SetTitle()
}

func (s *BarChart) SetData(time string) {
	meanTotal, labels := query.GetDataForBar(s.db, query.Build("mean(value)", s.From, s.Where, time, ""))
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
	s.C.BorderLabel = fmt.Sprintf("%s", s.BorderLabel)
}
