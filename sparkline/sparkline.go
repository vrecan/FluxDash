package sparkline

import (
	"fmt"
	log "github.com/cihub/seelog"
	H "github.com/dustin/go-humanize"
	ST "github.com/fatih/structs"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	"github.com/vrecan/FluxDash/query"
	"github.com/vrecan/FluxDash/timecop"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type SparkLines struct {
	SL    *ui.Sparklines `json:"-"`
	Lines []*SparkLine   `json:"lines"`

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
}

type SparkLine struct {
	SL         *ui.Sparkline `json:"-"`
	From       string        `json:"from"`
	db         *DB.Influx    `json:"-"`
	DataType   int           `json:"type"`
	Title      string        `json:"title"`
	TitleColor ui.Attribute  `json:"titlecolor"`
	Where      string        `json:"where"`
	Height     int           `json:"height"`
	LineColor  ui.Attribute  `json:"linecolor"`
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
	sStruct := ST.New(s)

	for _, field := range sStruct.Fields() {
		if field.Name() == "SL" || field.Name() == "Lines" {
			continue
		}
		slStruct := ST.New(s.SL)
		err := slStruct.Field(field.Name()).Set(field.Value())
		if nil != err {
			panic(err)
		}
	}

	for _, line := range s.Lines {
		line.db = db
		h += line.Height + 1
		spark := ui.NewSparkline()
		spark.LineColor = line.LineColor
		spark.TitleColor = line.TitleColor
		spark.Height = line.Height
		spark.Title = line.Title
		line.SL = &spark
	}

	s.SL.Height = h
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
