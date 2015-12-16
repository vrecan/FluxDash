package guage

import (
	"fmt"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/query"
)

type GaugeInfo struct {
	From   string
	Time   string
	Title  string
	Where  string
	Height int
	Border bool
}

type Gauge struct {
	I  GaugeInfo
	G  *ui.Gauge
	db *DB.Influx
}

func NewGauge(barColor ui.Attribute, db *DB.Influx, info GaugeInfo) *Gauge {

	g := &Gauge{G: ui.NewGauge(), db: db, I: info}
	g.G.BarColor = barColor
	// g.G.PercentColor = ui.ColorRed
	// g.G.PercentColorHighlighted = ui.ColorMagenta
	g.G.Border = info.Border
	g.G.Height = info.Height
	g.G.Label = info.Title
	return g
}

func (s *Gauge) Gauges() *ui.Gauge {
	return s.G
}

func (s *Gauge) Update(time string) {
	s.SetDataAndTitle(time)
}

func (s *Gauge) SetDataAndTitle(time string) {

}

func (s *Gauge) SetTitle(time string) {
	meanTotal := query.GetIntData(s.db, query.Build("mean(value)", s.I.From, s.I.Where, time, ""))
	if len(meanTotal) > 0 {
		s.G.Percent = meanTotal[0]
	} else {
		s.G.Percent = 0
	}
	maxTotal := query.GetIntData(s.db, query.Build("max(value)", s.I.From, s.I.Where, time, ""))
	s.G.Label = fmt.Sprintf("%s mean:%v%% max:%v%%", s.I.Title, s.G.Percent, maxTotal[0])
}
