package loading

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/vrecan/FluxDash/merge"
)

//Loading is a simple percentage gauge of the progress of the current dashboard loading.
type Loading struct {
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
	G                       *ui.Gauge    `json:"-"`
}

//NewLoading will create a gauge from a partial Loading generated from a json dashboard.
func NewLoading(g *Loading) *Loading {
	g.G = ui.NewGauge()
	merge.Merge(g, g.G, "G")
	return g
}

//Update the gauge data from influxdb queries.
func (s *Loading) Update(data int) {
	s.SetData(data)
	s.SetTitle()
}

//SetData will set the data for the bar.
func (s *Loading) SetData(data int) {
	s.G.Percent = data
}

//SetTitle will set the label of the gauge.
func (s *Loading) SetTitle() {
	if len(s.Label) <= 0 {
		s.Label = "Loading"
	}
	s.G.Label = fmt.Sprintf("%s %v%%", s.Label, s.G.Percent)

}
