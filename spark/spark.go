package spark

import (
	"errors"
	// log "github.com/cihub/seelog"
	ui "github.com/gizak/termui"
	INFLUX "github.com/influxdb/influxdb/client"
	DB "github.com/vrecan/FluxDash/influx"
)

type SparksConf struct {
	Query  string
	Title  string
	Height int
	Width  int
}

type Sparks struct {
	Data       []*INFLUX.Series
	SparkLines *ui.Sparklines
}

type SparkConf struct {
	Height     int
	Title      string
	LineColor  int
	TitleColor int
	Data       []int
}

//Spark is a set of sparklines derived from a query.
func NewSparks(conf SparksConf, db *DB.Influx) (sparks *Sparks, err error) {
	// log.Info("New Spark")
	sparks = &Sparks{}
	data, err := db.Query(conf.Query)
	if nil != err {
		return sparks, err
	}
	Data, ok := data.([]*INFLUX.Series)
	if !ok {
		return sparks, errors.New("Invalid data returned for sparks")
	}
	sparks.Data = Data
	color := int(2)
	temp := make([]ui.Sparkline, 0)
	for _, s := range sparks.Data {
		spark := newSpark(s, &SparkConf{
			Height:     5,
			Title:      s.GetName(),
			LineColor:  color,
			Data:       SeriesToInt(s),
			TitleColor: int(ui.ColorWhite)})
		temp = append(temp, spark)
		color++
	}
	sparks.SparkLines = ui.NewSparklines(temp...)
	sparks.SparkLines.SetWidth(conf.Width)
	sparks.SparkLines.BgColor = ui.Attribute(1)
	sparks.SparkLines.Border.Label = conf.Title
	sparks.SparkLines.Height = conf.Height
	sparks.SparkLines.IsDisplay = true

	return sparks, err
}

//convert to int for sparkline
func SeriesToInt(s *INFLUX.Series) (d []int) {
	d = make([]int, 0)
	for _, p := range s.GetPoints() {
		// data :=
		d = append(d, int(p[1].(float64))) //ignore time data
	}
	// log.Debug(d)
	return d
}

func (s *Sparks) Render() *ui.Sparklines {
	return s.SparkLines
}

//Create an individual sparkline that will go into sparks.
func newSpark(series *INFLUX.Series, conf *SparkConf) (spark ui.Sparkline) {
	spark = ui.NewSparkline()
	spark.Height = conf.Height
	spark.LineColor = ui.Attribute(conf.LineColor)
	spark.Title = conf.Title
	spark.TitleColor = ui.Attribute(conf.TitleColor)
	spark.Data = conf.Data
	// log.Debug(spark)
	return spark
}
