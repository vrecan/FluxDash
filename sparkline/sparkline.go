package sparkline

import (
	"encoding/json"
	"fmt"
	DASH "github.com/vrecan/FluxDash/dashboard"
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
	SL *ui.Sparkline
	D  DASH.SparkLineData
	db *DB.Influx
}

// type SparkLine struct {
// 	SL       *ui.Sparkline
// 	From     string
// 	Time     string
// 	db       *DB.Influx
// 	Title    string
// 	Where    string
// 	DataType int
// }
const (
	defaultHeight   = 3
	defaultInterval = "5s"
)

const (
	Short   = "short"
	Percent = "percent"
	Bytes   = "bytes"
	NS      = "time/ns" //nano second
)

// func NewSparkLinesFromData(s *SparkLinesData) *SparkLines {

// 	s.SL

// 	spark := ui.NewSparklines()
// 	sparkLines := SparkLines{SL: spark, lines: s.SL}
// 	h := defaultHeight
// 	for _, sl := range s.SL {
// 		h += sl.SL.Height + 1
// 	}
// 	spark.Height = h
// 	spark.BorderLabelFg = ui.ColorGreen | ui.AttrBold
// 	spark.Border = true
// 	return &sparkLines
// }

func NewSparkLineFromData(db *DB.Influx, d ...DASH.SparkLineData) (lines []*SparkLine) {
	for _, sl := range d {
		s := ui.NewSparkline()
		s.Height = sl.Height
		spark := &SparkLine{SL: &s, D: sl, db: db}
		lines = append(lines, spark)
	}
	return lines
}

func NewSparkLinesFromData(db *DB.Influx, d DASH.SparkLinesData) *SparkLines {
	uisparks := ui.NewSparklines()
	sparks := SparkLines{SL: uisparks, lines: NewSparkLineFromData(db, d.SL...)}
	return &sparks
}

// func NewSparkLines(s ...*SparkLine) *SparkLines {
// 	spark := ui.NewSparklines()
// 	sparkLines := SparkLines{SL: spark, lines: s}
// 	h := defaultHeight
// 	for _, sl := range s {
// 		h += sl.SL.Height + 1
// 	}
// 	spark.Height = h
// 	spark.BorderLabelFg = ui.ColorGreen | ui.AttrBold
// 	spark.Border = true
// 	return &sparkLines
// }

func (s *SparkLines) Sparks() *ui.Sparklines {
	return s.SL
}

func (s *SparkLines) Update() {
	var uiSparks []ui.Sparkline
	s.SetHeight()
	for _, sl := range s.lines {

		sl.SetData()
		sl.SetTitle()
		uiSparks = append(uiSparks, *sl.SL)
	}
	s.SL.Lines = uiSparks
}

func buildQuery(sel string, from string, where string, time string, groupBy string) string {
	if len(sel) == 0 || len(from) == 0 || len(time) == 0 {
		log.Fatal("invalid query string :", fmt.Sprintf("SELECT %s FROM %s WHERE %s AND time > %s %s", sel, from, where, groupBy))
	}
	if len(where) > 0 {
		return fmt.Sprintf("SELECT %s FROM %s WHERE %s AND time > %s %s", sel, from, where, time, groupBy)
	} else {
		return fmt.Sprintf("SELECT %s FROM %s WHERE time > %s %s", sel, from, time, groupBy)
	}
}

func (s *SparkLines) SetHeight() {
	h := defaultHeight
	for _, sl := range s.lines {
		h += sl.D.Height
	}
	s.Sparks().Height = h
}

func (s *SparkLine) SetData() {
	// s.SL.Data = getData(s.db, fmt.Sprintf("Select mean(value) FROM %s WHERE time > %s GROUP BY time(%s)", s.From, s.Time, defaultInterval))
	s.SL.Data = getData(s.db, buildQuery("mean(value)", s.D.From, s.D.Where, s.D.Time, fmt.Sprintf("GROUP BY time(%s)", defaultInterval)))
}

func (s *SparkLine) SetTitle() {
	// meanTotal := getData(s.db, fmt.Sprintf("Select mean(value) FROM %s WHERE time > %s", s.From, s.Time))
	meanTotal := getData(s.db, buildQuery("mean(value)", s.D.From, s.D.Where, s.D.Time, ""))
	maxTotal := getData(s.db, buildQuery("max(value)", s.D.From, s.D.Where, s.D.Time, ""))
	switch s.D.DataType {
	case Percent:
		s.SL.Title = fmt.Sprintf("%s mean:%v%% max:%v%%", s.D.Title, meanTotal[0], maxTotal[0])
	case Bytes:
		s.SL.Title = fmt.Sprintf("%s mean:%v max:%v", s.D.Title, H.Bytes(uint64(meanTotal[0])), H.Bytes(uint64(maxTotal[0])))
	case Short:
		s.SL.Title = fmt.Sprintf("%s mean:%v max:%v", s.D.Title, H.Comma(int64(meanTotal[0])), H.Comma(int64(maxTotal[0])))
	case NS:
		s.SL.Title = fmt.Sprintf("%s mean:%v max:%v", s.D.Title, timecop.GetCommaString(float64(meanTotal[0]), "nanoseconds"), timecop.GetCommaString(float64(maxTotal[0]), "nanoseconds"))
	default:
		log.Fatal("Data type is invalid: ", s.D.DataType)
	}

}

func getData(db *DB.Influx, q string) (data []int) {
	r, err := db.Query(q)
	if nil != err {
		log.Fatal(err)
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
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
