package multispark

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	H "github.com/dustin/go-humanize"

	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
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
	From     string
	Time     string
	Title    string
	Where    string
	DataType int
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

func (s *MultiSpark) SetDataAndTitle(time string, groupBy string) {
	data, labels := getData(s.db, buildQuery("mean(value)", s.I.From, s.I.Where, time, groupBy))
	meanTotal, _ := getData(s.db, buildQuery("mean(value)", s.I.From, s.I.Where, time, ""))
	maxTotal, _ := getData(s.db, buildQuery("max(value)", s.I.From, s.I.Where, time, ""))
	var uiSparks []ui.Sparkline
	for i, _ := range data {
		line := ui.NewSparkline()
		line.Data = data[i]
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
			log.Fatal("Data type is invalid: ", s.I.DataType)
		}
		uiSparks = append(uiSparks, line)
	}
	if s.SL == nil {
		s.SL = ui.NewSparklines(uiSparks...)
		s.SL.BorderLabel = s.I.Title
	} else {
		s.SL.Lines = uiSparks
	}

	s.SL.Height = 3 + len(data)*2

}
func (s *MultiSpark) GetColumns() []*ui.Row {
	return []*ui.Row{ui.NewCol(12, 0, s.Sparks())}
}
func getData(db *DB.Influx, q string) (data [][]int, labels []string) {
	r, err := db.Query(q)
	if nil != err {
		log.Fatal(err)
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
		log.Fatal(q)
	}
	labels = make([]string, len(r[0].Series))
	data = make([][]int, len(r[0].Series))
	for i, result := range r[0].Series {
		labels[i] = result.Name
		for _, row := range result.Values {
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
					data[i] = append(data[i], int(val))
				}
			}
		}
	}
	return data, labels
}
