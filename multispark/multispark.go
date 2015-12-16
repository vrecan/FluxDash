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

type MultiSparkInfo struct {
	From       string
	Time       string
	Title      string
	Where      string
	DataType   int
	SpanSize   int
	LineColor  ui.Attribute
	TitleColor ui.Attribute
}

type MultiSpark struct {
	SL.SparkLines
	db *DB.Influx
	I  MultiSparkInfo
}

func NewMultiSpark(db *DB.Influx, i MultiSparkInfo) *MultiSpark {
	ms := &MultiSpark{db: db, I: i}
	ms.SetDataAndTitle(fmt.Sprintf("now() - %s", "15m"), fmt.Sprintf("GROUP BY time(%s)", "5s"))
	return ms
}

func (s *MultiSpark) Update(time string, groupBy string) {
	s.SetDataAndTitle(time, groupBy)
}

func (s *MultiSpark) SetDataAndTitle(time string, groupBy string) {
	data, labels := query.GetIntDataFromTags(s.db, query.Build("mean(value)", s.I.From, s.I.Where, time, groupBy))
	meanTotal, _ := query.GetIntDataFromTags(s.db, query.Build("mean(value)", s.I.From, s.I.Where, time, ""))
	maxTotal, _ := query.GetIntDataFromTags(s.db, query.Build("max(value)", s.I.From, s.I.Where, time, ""))
	var uiSparks []ui.Sparkline
	for i, _ := range data {
		line := ui.NewSparkline()
		line.Data = data[i]
		line.LineColor = s.I.LineColor
		line.TitleColor = s.I.TitleColor
		switch s.I.DataType {
		case Percent:
			line.Title = fmt.Sprintf("%s mean:%v%% max:%v%% cur: %v", labels[i], meanTotal[i][0], maxTotal[i][0], data[i][len(data[i])-1])
		case Bytes:
			line.Title = fmt.Sprintf("%s mean:%v max:%v cur: %v", labels[i], H.Bytes(uint64(meanTotal[i][0])), H.Bytes(uint64(maxTotal[i][0])), H.Bytes(uint64(data[i][len(data[i])-1])))
		case Short:
			line.Title = fmt.Sprintf("%s mean:%v max:%v cur: %v", labels[i], H.Comma(int64(meanTotal[i][0])), H.Comma(int64(maxTotal[i][0])), H.Comma(int64(data[i][len(data[i])-1])))
		case Time:
			line.Title = fmt.Sprintf("%s mean:%v max:%v cur: %v", labels[i], timecop.GetCommaString(float64(meanTotal[i][0]), "nanoseconds"), timecop.GetCommaString(float64(maxTotal[i][0]), "nanoseconds"), timecop.GetCommaString(float64(data[i][len(data[i])-1]), "nanoseconds"))
		default:
			log.Critical("Data type is invalid: ", s.I.DataType)
		}
		uiSparks = append(uiSparks, line)
	}
	if s.SL == nil {
		s.SL = ui.NewSparklines(uiSparks...)
	} else {
		s.SL.Lines = uiSparks
	}
	s.SL.BorderLabel = s.I.Title

	s.SL.Height = 3 + len(data)*2

}
func (s *MultiSpark) GetColumns() []*ui.Row {
	return []*ui.Row{ui.NewCol(12, 0, s.SL)}
}
