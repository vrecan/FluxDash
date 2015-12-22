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

type MultiSpark struct {
	From        string       `json:"from"`
	BorderLabel string       `json:"borderlabel"`
	Border      bool         `json:"border"`
	Where       string       `json:"where"`
	DataType    int          `json:"type"`
	LineColor   ui.Attribute `json:"linecolor"`
	TitleColor  ui.Attribute `json:"titlecolor"`

	SL.SparkLines
	db *DB.Influx `json:"-"`
}

func NewMultiSpark(db *DB.Influx, ms *MultiSpark) *MultiSpark {
	ms.db = db
	ms.SL = ui.NewSparklines()
	log.Info(ms)
	return ms
}

func (s *MultiSpark) Update(time TS.TimeSelect) {
	t, groupByInterval, _ := time.CurTime()
	s.SetDataAndTitle(t, groupByInterval)
}

func (s *MultiSpark) SetDataAndTitle(time string, groupBy string) {
	data, labels := query.GetIntDataFromTags(s.db, query.Build("mean(value)", s.From, s.Where, time, groupBy))
	meanTotal, _ := query.GetIntDataFromTags(s.db, query.Build("mean(value)", s.From, s.Where, time, ""))
	maxTotal, _ := query.GetIntDataFromTags(s.db, query.Build("max(value)", s.From, s.Where, time, ""))
	var uiSparks []ui.Sparkline
	for i, _ := range data {
		line := ui.NewSparkline()
		line.Data = data[i]
		line.LineColor = s.LineColor
		line.TitleColor = s.TitleColor
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
	if s.SL == nil {
		s.SL = ui.NewSparklines(uiSparks...)
	} else {
		s.SL.Lines = uiSparks
	}
	s.SL.BorderLabel = s.BorderLabel
	s.SL.Border = s.Border

	s.SL.Height = 3 + len(data)*2

}
