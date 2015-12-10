package guage

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
)

type Gauge struct {
	G     *ui.Gauge
	From  string
	Time  string
	db    *DB.Influx
	Title string
	Where string
}

func NewGauge(s ui.Gauge, from string, time string, db *DB.Influx, title string, where string) *Gauge {
	g := &Gauge{G: &s, From: from, Time: time, db: db, Title: title, Where: where}
	return g
}

func (s *Gauge) Gauges() *ui.Gauge {
	return s.G
}

func (s *Gauge) Update() {
	s.SetData()
	s.SetTitle()
}

func (s *Gauge) SetData() {
	meanTotal := getData(s.db, buildQuery("mean(value)", s.From, s.Where, s.Time, ""))
	s.G.Percent = meanTotal[0]
}

func (s *Gauge) SetTitle() {
	meanTotal := getData(s.db, buildQuery("mean(value)", s.From, s.Where, s.Time, ""))
	maxTotal := getData(s.db, buildQuery("max(value)", s.From, s.Where, s.Time, ""))
	s.G.Label = fmt.Sprintf("%s mean:%v%% max:%v%%", s.Title, meanTotal[0], maxTotal[0])
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
