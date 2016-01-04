package multispark

import (
	"fmt"
	log "github.com/cihub/seelog"
	H "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/merge"
	"github.com/vrecan/FluxDash/query"
	SL "github.com/vrecan/FluxDash/sparkline"
	"github.com/vrecan/FluxDash/timecop"
	TS "github.com/vrecan/FluxDash/timeselect"
)

const (
	defaultHeight   = 3
	defaultInterval = "5s"
)
const (
	Short   = 1
	Percent = 2
	Bytes   = 3
	Time    = 4
)

//MultiSpark is a collection of sparklines generated based on tags from an influxdb query.
type MultiSpark struct {
	From  string `json:"from"`
	Where string `json:"where"`

	SL.SparkLines
	db DB.DBI `json:"-"`

	BorderLabel   string       `json:"borderlabel"`
	Border        bool         `json:"border"`
	BorderFg      ui.Attribute `json:"borderfg"`
	BorderBg      ui.Attribute `json:"borderbg"`
	BorderLeft    bool         `json:borderleft"`
	BorderRight   bool         `json:"borderright"`
	BorderTop     bool         `json:"bordertop"`
	BorderBottom  bool         `json:"borderbottom"`
	BorderLabelFg ui.Attribute `json:"borderlabelfg"`
	BorderLabelBg ui.Attribute `json:"borderlabelbg"`
	Display       bool         `json:"display"`
	Bg            ui.Attribute `json:"bg"`
	Width         int          `json:"width"`
	Height        int          `json:"height"`
	PaddingTop    int          `json:"paddingtop"`
	PaddingBottom int          `json:"paddingbottom"`
	PaddingLeft   int          `json:"paddingleft"`
	PaddingRight  int          `json:"paddingright"`
	AutoColor     bool         `json:"autocolor"`
	DataType      int          `json:"type"`
	LineColor     ui.Attribute `json:"linecolor"`
	TitleColor    ui.Attribute `json:"titlecolor"`
}

//NewMultiSpark builds a multispark from a partial multispark that has been generated from a json dashboard.
func NewMultiSpark(db DB.DBI, ms *MultiSpark) *MultiSpark {
	ms.db = db
	ms.SL = ui.NewSparklines()
	merge.Merge(ms, ms.SL, "db", "from", "where")
	return ms
}

var colors = []ui.Attribute{ui.ColorWhite, ui.ColorGreen, ui.ColorMagenta, ui.ColorRed, ui.ColorYellow, ui.ColorBlack, ui.ColorBlue, ui.ColorCyan}

//Update a multispark will requery influxdb to update the sparklines.
func (s *MultiSpark) Update(time TS.TimeSelect) {
	t, groupByInterval, _ := time.CurTime()
	s.SetDataAndTitle(t, groupByInterval)
}

//SetDatAndTitle will update all the data for all sparklines in the multispark.
func (s *MultiSpark) SetDataAndTitle(time string, groupBy string) {
	data, labels := query.GetIntDataFromTags(s.db, query.Build("mean(value)", s.From, s.Where, time, groupBy))
	meanTotal, _ := query.GetIntDataFromTags(s.db, query.Build("mean(value)", s.From, s.Where, time, ""))
	maxTotal, _ := query.GetIntDataFromTags(s.db, query.Build("max(value)", s.From, s.Where, time, ""))
	autoColorc := 0
	uiSparks := make([]ui.Sparkline, 0)
	for i, _ := range data {
		DataExists := true
		line := ui.NewSparkline()
		line.Height = 1
		if s.AutoColor {
			if autoColorc > len(colors)-1 {
				autoColorc = 0
			}
			//skip bg color so that we don't paint black on black...
			if colors[autoColorc] == s.Bg {
				autoColorc++
			}
			line.LineColor = colors[autoColorc]
			line.TitleColor = colors[autoColorc]
			autoColorc++
		} else {
			line.LineColor = s.LineColor
			line.TitleColor = s.TitleColor
		}
		if len(data[i]) < 1 || len(meanTotal[i]) < 1 || len(maxTotal[i]) < 1 {
			DataExists = false
		}

		if DataExists {
			for in, d := range data[i] {
				if d < 0 {
					data[i][in] = 0
				}
			}
			line.Data = data[i]
			switch s.DataType {

			case Percent:
				line.Title = fmt.Sprintf("%s mean:%v%% max:%v%% cur: %v", labels[i], meanTotal[i][0], maxTotal[i][0], data[i][len(data[i])-1])
			case Bytes:
				line.Title = fmt.Sprintf("%s mean:%v max:%v cur: %v", labels[i], H.Bytes(uint64(meanTotal[i][0])), H.Bytes(uint64(maxTotal[i][0])), H.Bytes(uint64(data[i][len(data[i])-1])))
			case Short:
				line.Title = fmt.Sprintf("%s mean:%v max:%v cur: %v", labels[i], H.Comma(int64(meanTotal[i][0])), H.Comma(int64(maxTotal[i][0])), H.Comma(int64(data[i][len(data[i])-1])))
			case Time:
				line.Title = fmt.Sprintf("%s mean:%v max:%v cur: %v", labels[i], timecop.GetCommaString(float64(meanTotal[i][0]), "nanoseconds"), timecop.GetCommaString(float64(maxTotal[i][0]), "nanoseconds"), timecop.GetCommaString(float64(data[i][len(data[i])-1]), "nanoseconds"))
			default:
				log.Critical("Data type is invalid: ", s.DataType)
			}
		} else {
			if len(labels[i]) < 1 {
				line.Title = "No Data"
			}
			line.Title = fmt.Sprintf("%s No Data")
		}
		uiSparks = append(uiSparks, line)
	}
	if len(uiSparks) > 0 {
		s.SL.Lines = uiSparks
		s.SL.Height = 3 + len(data)*2
	}

}
