package guage

import (
	"fmt"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/merge"
	"github.com/vrecan/FluxDash/query"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type GaugeInfo struct {
}

type Gauge struct {
	From                    string       `json:"from"`
	Where                   string       `json:"where"`
	BorderLabel             string       `json:"borderlabel"`
	Border                  bool         `json:"border"`
	BorderFg                ui.Attribute `json:"borderfg"`
	BorderBg                ui.Attribute `json:"borderbg"`
	BorderLeft              bool         `json:borderleft"`
	BorderRight             bool         `json:"borderright"`
	BorderTop               bool         `json:"bordertop"`
	BorderBottom            bool         `json:"borderbottom"`
	BorderLabelFg           ui.Attribute `json:"borderlabelfg"`
	BorderLabelBg           ui.Attribute `json:"borderlabelbg"`
	Display                 bool         `json:"display"`
	Bg                      ui.Attribute `json:"bg"`
	Width                   int          `json:"width"`
	Height                  int          `json:"height"`
	PaddingTop              int          `json:"paddingtop"`
	PaddingBottom           int          `json:"paddingbottom"`
	PaddingLeft             int          `json:"paddingleft"`
	PaddingRight            int          `json:"paddingright"`
	BarColor                ui.Attribute `json:"barcolor"`
	PercentColor            ui.Attribute `json:"percentcolor"`
	PercentColorHighlighted ui.Attribute `json:"percentcolorhighlighted"`
	Label                   string       `json:"label"`
	LabelAlign              ui.Align     `json:"labelalign"`
	G                       *ui.Gauge
	db                      *DB.Influx
}

func NewGauge(db *DB.Influx, g *Gauge) *Gauge {
	g.db = db
	g.G = ui.NewGauge()
	merge.Merge(g, g.G, "G", "db")
	return g
}

func (s *Gauge) Gauges() *ui.Gauge {
	return s.G
}

func (s *Gauge) Update(time TS.TimeSelect) {
	t, _, _ := time.CurTime()
	s.SetData(t)
	s.SetTitle(t)
}

func (s *Gauge) SetData(time string) {
	meanTotal := query.GetIntData(s.db, query.Build("mean(value)", s.From, s.Where, time, ""))
	s.G.Percent = meanTotal[0]
}

func (s *Gauge) SetTitle(time string) {
	meanTotal := query.GetIntData(s.db, query.Build("mean(value)", s.From, s.Where, time, ""))
	if len(meanTotal) > 0 {
		s.G.Percent = meanTotal[0]
	} else {
		s.G.Percent = 0
	}
	maxTotal := query.GetIntData(s.db, query.Build("max(value)", s.From, s.Where, time, ""))
	s.G.Label = fmt.Sprintf("%s mean:%v%% max:%v%%", s.Label, s.G.Percent, maxTotal[0])
}
