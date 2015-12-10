package sparkline

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	H "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/timecop"
)

type SparkLines struct {
	SL      *ui.Sparklines
	Refresh time.Duration
	lines   []*SparkLine
}

type SparkLine struct {
	SL       *ui.Sparkline
	From     string
	Time     string
	db       *DB.Influx
	Title    string
	Where    string
	DataType int
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

func NewSparkLine(s ui.Sparkline, from string, time string, db *DB.Influx, title string) *SparkLine {
	sl := &SparkLine{SL: &s, From: from, Time: time, db: db, Title: title, DataType: Short}
	return sl
}

func (s *SparkLines) Update() {
	var uiSparks []ui.Sparkline
	for _, sl := range s.lines {
		sl.SetData()
		sl.SetTitle()
		uiSparks = append(uiSparks, *sl.SL)
	}
	s.SL.Lines = uiSparks
}

func (s *SparkLine) SetData() {
	s.SL.Data = getData(s.db, fmt.Sprintf("Select mean(value) FROM %s WHERE time > %s GROUP BY time(%s)", s.From, s.Time, defaultInterval))
}

func (s *SparkLine) SetTitle() {
	meanTotal := getData(s.db, fmt.Sprintf("Select mean(value) FROM %s WHERE time > %s", s.From, s.Time))
	var maxTotal []int
	if len(s.Where) > 0 {
		maxTotal = getData(s.db, fmt.Sprintf("Select max(value) FROM %s WHERE %s AND time > %s", s.From, s.Time))
	} else {
		maxTotal = getData(s.db, fmt.Sprintf("Select max(value) FROM %s WHERE time > %s", s.From, s.Time))
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
		log.Fatal("Data type is invalid: ", s.DataType)
	}

}

func getData(db *DB.Influx, q string) (data []int) {
	r, err := db.Query(q)
	if nil != err {
		log.Fatal(err)
	}
	if len(r) < 1 {
		log.Fatal(q)
	}
	for _, row := range r[0].Series[0].Values {
		_, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		if len(row) > 1 {
			if nil != row[1] {
				val, err := row[1].(json.Number).Float64()
				if nil != err {
					fmt.Println("ERR: ", err)
				}
				data = append(data, int(val))
			}
		}

	}
	return data
}
