package sparkline

import (
	"fmt"
	log "github.com/cihub/seelog"
	H "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/query"
	"github.com/vrecan/FluxDash/timecop"
)

type SparkLines struct {
	SL    *ui.Sparklines `json:"-"`
	lines []*SparkLine   `json:"lines"`
}

type SparkLine struct {
	SL       *ui.Sparkline `json:"-"`
	From     string        `json:"from"`
	Time     string        `json:"time"`
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

func NewSparkLines(s ...*SparkLine) *SparkLines {
	spark := ui.NewSparklines()
	sparkLines := SparkLines{SL: spark, lines: s}
	h := defaultHeight
	for _, sl := range s {
		h += sl.SL.Height + 1
	}
	spark.Height = h
	spark.BorderLabelFg = ui.ColorGreen | ui.AttrBold
	spark.Border = true
	return &sparkLines
}

func (s *SparkLines) Sparks() *ui.Sparklines {
	return s.SL
}

func NewSparkLinex(s *SparkLine, db *DB.Influx) {
	s.db = db
	spark := ui.NewSparkline()
	s.SL = &spark
	s.SL.Height = s.Height
	s.SL.Title = s.Title
}

func NewSparkLine(s ui.Sparkline, from string, db *DB.Influx, title string, where string) *SparkLine {
	sl := &SparkLine{SL: &s, From: from, db: db, Title: title, DataType: Short, Where: where}
	return sl
}

func (s *SparkLines) Update(time string, groupBy string) {
	var uiSparks []ui.Sparkline
	for _, sl := range s.lines {
		sl.SetData(time, groupBy)
		sl.SetTitle(time)
		uiSparks = append(uiSparks, *sl.SL)
	}
	s.SL.Lines = uiSparks
}

func (s *SparkLine) SetData(time string, groupBy string) {
	// s.SL.Data = getData(s.db, fmt.Sprintf("Select mean(value) FROM %s WHERE time > %s GROUP BY time(%s)", s.From, s.Time, defaultInterval))
	s.SL.Data = query.GetIntData(s.db, query.Build("mean(value)", s.From, s.Where, time, groupBy))
}

func (s *SparkLine) SetTitle(time string) {
	// meanTotal := getData(s.db, fmt.Sprintf("Select mean(value) FROM %s WHERE time > %s", s.From, s.Time))
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
