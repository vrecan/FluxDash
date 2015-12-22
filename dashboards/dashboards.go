package dashboards

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	ui "github.com/gizak/termui"
	BC "github.com/vrecan/FluxDash/barchart"
	G "github.com/vrecan/FluxDash/gauge"
	DB "github.com/vrecan/FluxDash/influx"
	MS "github.com/vrecan/FluxDash/multispark"
	SL "github.com/vrecan/FluxDash/sparkline"
	TP "github.com/vrecan/FluxDash/timep"
	TS "github.com/vrecan/FluxDash/timeselect"
	"io/ioutil"
	// "os"
)

type Dashboard struct {
	Rows []*Row         `json:"rows"`
	Time *TS.TimeSelect `json:"-"` //TODO: Remove this?? It's in monitor shouldn't be in 2 places
	db   *DB.Influx     `json:"-"`
	Grid *ui.Grid       `json:"-"` //json:"-" omits a field from being encoded
}

type P struct {
	Height int     `json:"height"`
	Text   string  `json:"text"`
	Par    *ui.Par `json:"-"`
	Border bool    `json:"border"`
}

type Row struct {
	Height  int       `json:"height"`
	Span    int       `json:"span"`
	Offset  int       `json:"offset"`
	row     *ui.Row   `json:"-"`
	Columns []*Column `json:"columns"`
}

type Column struct {
	Height     int            `json:"height"`
	Span       int            `json:"span"`
	Offset     int            `json:"offset"`
	row        *ui.Row        `json:"-"`
	P          *P             `json:"p,omitempty"`
	TimeP      *TP.TimeP      `json:"timep,omitempty"`
	SparkLines *SL.SparkLines `json:"sparklines,omitempty"`
	BarChart   *BC.BarChart   `json:"barchart,omitempty"`
	Gauge      *G.Gauge       `json:"gauge,omitempty"`
	MultiSpark *MS.MultiSpark `json:"multispark,omitempty"`
}

func CreateExampleDash() {
	dash := ExampleDash(nil)

	raw, err := json.Marshal(dash)
	fmt.Println("err: ", err)
	fmt.Println("json: ", string(raw))
}

//ExampleDash returns an example dashboard with all basic stuff filled out.
func ExampleDash(db *DB.Influx) *Dashboard {
	dash := NewDashboard(db)
	sl1 := &SL.SparkLine{Title: "CPU", From: "/system.cpu/", Where: "", Height: 1, DataType: 1}
	sl2 := &SL.SparkLine{Title: "Dispatch GC", Height: 1, From: "/gc.pause.ns/", Where: `"service"= 'godispatch'`, DataType: 1}
	sparks := make([]*SL.SparkLine, 0)
	sparks = append(sparks, sl1, sl2)
	sparkLines := &Column{Height: 1, Span: 12, SparkLines: &SL.SparkLines{BorderLabel: "System", Border: true, Lines: sparks}}
	columns := make([]*Column, 0)
	columns = append(columns, sparkLines)
	dash.Rows = append(dash.Rows, &Row{Height: 1, Span: 12, Offset: 0, Columns: columns})
	p1 := &Column{Height: 1, Span: 6, Offset: 0, P: &P{Text: "Static text is all the rage!!!", Height: 3, Border: true}}
	columns2 := make([]*Column, 0)
	columns2 = append(columns2, p1)
	dash.Rows = append(dash.Rows, &Row{Height: 1, Span: 12, Offset: 0, Columns: columns2})
	ptRow := &Column{Height: 1, Span: 6, Offset: 0, TimeP: &TP.TimeP{Height: 3, Border: true}}
	columns3 := make([]*Column, 0)
	columns3 = append(columns3, ptRow)

	dash.Rows = append(dash.Rows, &Row{Height: 1, Span: 12, Offset: 0, Columns: columns3})

	return dash
}

func NewDashboard(db *DB.Influx) *Dashboard {
	return &Dashboard{db: db}

}

//Dashboard get dash from path.
func NewDashboardFromFile(db *DB.Influx, f string) *Dashboard {
	mem, e := ioutil.ReadFile(f)
	if e != nil {
		log.Critical("File error: ", e)
	}

	// var jsontype jsonobject
	dash := &Dashboard{}
	err := json.Unmarshal(mem, dash)
	if nil != err {
		panic(err)
	}
	dash.db = db
	dash.Time = &TS.TimeSelect{}
	return dash
}

func (d *Dashboard) Create() {

	rows := make([]*ui.Row, 0)

	for _, r := range d.Rows {
		r.row = ui.NewRow()
		r.row.Height = r.Height
		r.row.Span = r.Span
		r.row.Offset = r.Offset
		columns := make([]*ui.Row, 0)
		for _, c := range r.Columns {
			if nil != c.P {
				par := ui.NewPar(c.P.Text)
				par.Border = c.P.Border
				par.Height = c.P.Height
				c.P.Par = par
				col := ui.NewCol(c.Span, c.Offset, par)
				columns = append(columns, col)
			} else if nil != c.TimeP {
				c.TimeP = TP.NewTimeP(c.TimeP)
				col := ui.NewCol(c.Span, c.Offset, c.TimeP.Par)
				columns = append(columns, col)
			} else if nil != c.SparkLines {
				c.SparkLines = SL.NewSparkLines(d.db, c.SparkLines)
				col := ui.NewCol(c.Span, c.Offset, c.SparkLines.SL)
				columns = append(columns, col)
			} else if nil != c.BarChart {
				c.BarChart = BC.NewBarChart(d.db, c.BarChart)
				colBar := ui.NewCol(c.Span, c.Offset, c.BarChart.BarCharts())
				colLabel := ui.NewCol(c.Span, c.Offset, c.BarChart.Labels())
				columns = append(columns, colBar, colLabel)
			} else if nil != c.MultiSpark {
				c.MultiSpark = MS.NewMultiSpark(d.db, c.MultiSpark)
				col := ui.NewCol(c.Span, c.Offset, c.MultiSpark.SL)
				columns = append(columns, col)
			} else if nil != c.Gauge {
				c.Gauge = G.NewGauge(d.db, c.Gauge)
				col := ui.NewCol(c.Span, c.Offset, c.Gauge.G)
				columns = append(columns, col)
			}
		}
		r.row.Cols = columns
		rows = append(rows, r.row)
	}
	d.Grid = ui.NewGrid(rows...)
	d.Grid.BgColor = ui.ThemeAttr("bg")
	d.Grid.Width = ui.TermWidth()
	d.Grid.Align()
}

func asyncUpdate(f func(TS.TimeSelect), t TS.TimeSelect, done chan bool) {
	go func() {
		f(t)
		done <- true
	}()
}

func (d *Dashboard) UpdateAll(time *TS.TimeSelect) {
	finChan := make(chan bool, 0)
	exp := 0
	for _, r := range d.Rows {
		for _, c := range r.Columns {
			if nil != c.TimeP {
				exp++
				asyncUpdate(c.TimeP.Update, *time, finChan)
			} else if nil != c.SparkLines {
				exp++
				asyncUpdate(c.SparkLines.Update, *time, finChan)
			} else if nil != c.BarChart {
				exp++
				asyncUpdate(c.BarChart.Update, *time, finChan)
			} else if nil != c.MultiSpark {
				exp++
				asyncUpdate(c.MultiSpark.Update, *time, finChan)
			} else if nil != c.Gauge {
				exp++
				asyncUpdate(c.Gauge.Update, *time, finChan)
			}
		}
	}
	//select with timer might be more robust here
	for _ = range finChan {
		exp--
		d.Grid.Align()
		ui.Render(d.Grid)
		if exp == 0 {
			break
		}
	}
}

func (d *Dashboard) GetGrid() *ui.Grid {
	return d.Grid
}
