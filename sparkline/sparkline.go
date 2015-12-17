package sparkline

import (
	"fmt"
	log "github.com/cihub/seelog"
	H "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/query"
	"github.com/vrecan/FluxDash/timecop"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type SparkLines struct {
	SL          *ui.Sparklines `json:"-"`
	Lines       []*SparkLine   `json:"lines"`
	BorderLabel string         `json:"borderlabel"`
	Border      bool           `json:"border"`
}

type SparkLine struct {
	SL       *ui.Sparkline `json:"-"`
	From     string        `json:"from"`
	db       *DB.Influx    `json:"-"`
	Title    string        `json:"title"`
	Where    string        `json:"where"`
	DataType int           `json:"type"`
	Height   int           `json:"height"`
}

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

func NewSparkLines(db *DB.Influx, s *SparkLines) *SparkLines {
	s.SL = ui.NewSparklines()
	h := defaultHeight

	for _, line := range s.Lines {
		line.db = db
		h += line.Height + 1
		spark := ui.NewSparkline()
		spark.Height = line.Height
		spark.Title = line.Title
		line.SL = &spark
	}
	s.SL.Height = h
	s.SL.BorderLabelFg = ui.ColorGreen | ui.AttrBold
	s.SL.Border = true
	return s
}

func (s *SparkLines) Update(time TS.TimeSelect) {
	var uiSparks []ui.Sparkline
	t, groupByInterval, _ := time.CurTime()
	for _, sl := range s.Lines {
		sl.SetData(t, groupByInterval)
		sl.SetTitle(t)
		uiSparks = append(uiSparks, *sl.SL)
	}
	s.SL.Lines = uiSparks
}

func (s *SparkLine) SetData(time string, groupBy string) {
	s.SL.Data = query.GetIntData(s.db, query.Build("mean(value)", s.From, s.Where, time, groupBy))
}

func (s *SparkLine) SetTitle(time string) {
	meanTotal := query.GetIntData(s.db, query.Build("mean(value)", s.From, s.Where, time, ""))
	maxTotal := query.GetIntData(s.db, query.Build("max(value)", s.From, s.Where, time, ""))
	if len(meanTotal) < 1 || len(maxTotal) < 1 {
		log.Error("No data for mean/max totals")
		s.SL.Title = fmt.Sprintf("%s No data", s.Title)
		return
	}
	switch s.DataType {
	case Percent:
		s.SL.Title = fmt.Sprintf("%s mean:%v%% max:%v%%", s.Title, meanTotal[0], maxTotal[0])
	case Bytes:
		s.SL.Title = fmt.Sprintf("%s mean:%v max:%v", s.Title, H.Bytes(uint64(meanTotal[0])), H.Bytes(uint64(maxTotal[0])))
	case Short:
		s.SL.Title = fmt.Sprintf("%s mean:%v max:%v", s.Title, H.Comma(int64(meanTotal[0])), H.Comma(int64(maxTotal[0])))
	case Time:
		s.SL.Title = fmt.Sprintf("%s mean:%v max:%v", s.Title, timecop.GetCommaString(float64(meanTotal[0]), "nanoseconds"), timecop.GetCommaString(float64(maxTotal[0]), "nanoseconds"))
	default:
		log.Critical("Data type is invalid: ", s.DataType)
	}

}
