package multispark

import (
	"fmt"
	log "github.com/cihub/seelog"
	H "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
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
	From        string       `json:"from"`
	BorderLabel string       `json:"borderlabel"`
	Border      bool         `json:"border"`
	Where       string       `json:"where"`
	DataType    int          `json:"type"`
	Bg          ui.Attribute `json:"bg"`
	LineColor   ui.Attribute `json:"linecolor"`
	TitleColor  ui.Attribute `json:"titlecolor"`
	AutoColor   bool         `json:"autocolor"`
	SL.SparkLines
	db *DB.Influx `json:"-"`
}

//NewMultiSpark builds a multispark from a partial multispark that has been generated from a json dashboard.
func NewMultiSpark(db *DB.Influx, ms *MultiSpark) *MultiSpark {
	ms.db = db
	ms.SL = ui.NewSparklines()
	ms.SL.Bg = ms.Bg
	ms.SL.BorderLabel = ms.BorderLabel
	ms.SL.Border = ms.Border
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
	var uiSparks []ui.Sparkline
	for i, _ := range data {
		line := ui.NewSparkline()
		line.Data = data[i]
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
		uiSparks = append(uiSparks, line)
	}
	s.SL.Lines = uiSparks
	s.SL.Height = 3 + len(data)*2

}
