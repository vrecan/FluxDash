package barchart

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
)

type BarChartInfo struct {
	From  string
	Time  string
	Title string
	Where string
}

type BarChart struct {
	I  BarChartInfo
	C  *ui.BarChart
	L  *ui.List
	db *DB.Influx
}

func NewBarChart(db *DB.Influx, info BarChartInfo) *BarChart {
	barchart := ui.NewBarChart()
	list := ui.NewList()
	g := &BarChart{C: barchart, L: list, db: db, I: info}
	g.C.DataLabels = make([]string, 0)
	g.C.Height = 10
	g.C.Width = 10
	g.L.Width = 10
	g.L.Height = 10
	g.L.BorderLabel = "Shard States"
	g.L.Y = 0
	// g.L.ItemFgColor = ui.ColorBlack
	// g.L.ItemBgColor = ui.ColorWhite
	return g
}

func (s *BarChart) BarCharts() *ui.BarChart {
	return s.C
}
func (s *BarChart) Labels() *ui.List {
	return s.L
}
func (s *BarChart) Update() {
	s.SetData()
	s.SetTitle()
}

func (s *BarChart) SetData() {
	meanTotal, labels := getData(s.db, buildQuery("mean(value)", s.I.From, s.I.Where, s.I.Time, ""))
	s.C.Data = meanTotal
	series := make([]string, len(labels))
	items := make([]string, len(labels))
	i := 0
	for _, v := range labels {
		series[i] = v[0]
		items[i] = fmt.Sprintf("%s: %s", v[0], v[1])
		i++
	}
	s.C.DataLabels = series
	s.L.Items = items
}
func (s *BarChart) GetColumns() []*ui.Row {
	return []*ui.Row{ui.NewCol(6, 0, s.BarCharts()), ui.NewCol(6, 0, s.Labels())}
}

func (s *BarChart) SetTitle() {
	s.C.BorderLabel = fmt.Sprintf("%s", s.I.Title)
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
func getData(db *DB.Influx, q string) (data []int, labels [][]string) {
	r, err := db.Query(q)
	if nil != err {
		log.Fatal(err)
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
		log.Fatal(q)
	}
	labels = make([][]string, len(r[0].Series))
	for i, result := range r[0].Series {
		series := fmt.Sprintf("S%d", i)
		labels[i] = []string{series, result.Name}
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
					data = append(data, int(val))
				}
			}

		}
	}

	return data, labels
}
