package dashboards

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
	SL "github.com/vrecan/FluxDash/sparkline"
	TS "github.com/vrecan/FluxDash/timeselect"
	"io/ioutil"
	// "os"
)

type Dashboard struct {
	Rows []Row          `json:"rows"`
	Time *TS.TimeSelect `json:"-"`
	db   *DB.Influx     `json:"-"`
	Grid *ui.Grid       `json:"-` //json:"-" omits a field from being encoded
}

// type SparkLine struct {
// 	Title    string        `json:"title"`
// 	Height   int           `json:"height"`
// 	From     string        `json:"from"`
// 	Where    string        `json:"where"`
// 	DataType string        `json:"dataType"`
// 	SL       SL.SparkLines `json:-"`
// }

// type SparkLines struct {
// 	SL []SparkLine `json:"sparkline"`
// }

type P struct {
	Height int     `json:"height"`
	Text   string  `json:"text"`
	Par    *ui.Par `json:"-"`
	Border bool    `json:"border"`
}

type Row struct {
	Height  int      `json:"height"`
	Span    int      `json:"span"`
	Offset  int      `json:"offset"`
	row     *ui.Row  `json:"-"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Height int         `json:"height"`
	Span   int         `json:"span"`
	Offset int         `json:"offset"`
	row    *ui.Row     `json:"-"`
	Widget interface{} `json:"widget"`
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

	r1 := Column{Height: 1, Span: 6, Offset: 0, Widget: SL.SparkLine{Title: "CPU", From: "/system.cpu/", Where: "", Height: 1, DataType: 1}}
	r2 := Column{Height: 1, Span: 6, Offset: 6, Widget: SL.SparkLine{Title: "Dispatch GC", Height: 1, From: "/gc.pause.ns/", Where: `"service"= 'godispatch'`, DataType: 1}}
	columns := make([]Column, 0)
	columns = append(columns, r1)
	columns = append(columns, r2)
	dash.Rows = append(dash.Rows, Row{Height: 1, Span: 12, Offset: 0, Columns: columns})
	p1 := Column{Height: 1, Span: 6, Offset: 0, Widget: P{Text: "!!!WOOOO!!!", Height: 3, Border: true}}
	columns2 := make([]Column, 0)
	columns2 = append(columns2, p1)
	dash.Rows = append(dash.Rows, Row{Height: 1, Span: 12, Offset: 0, Columns: columns2})
	// dash.Lines.SL = append(dash.Lines.SL, SparkLineData{From: "/system.cpu/", Time: "now - 15m", Title: "CPU", Where: "", DataType: "percent"})

	return dash
}

func NewDashboard(db *DB.Influx) *Dashboard {
	return &Dashboard{db: db, Time: &TS.TimeSelect{}}

}

//Dashboard get dash from path.
func NewDashboardFromFile(f string) *Dashboard {
	mem, e := ioutil.ReadFile(f)
	if e != nil {
		log.Critical("File error: ", e)
	}
	fmt.Printf("%s\n", string(mem))

	// var jsontype jsonobject
	dash := &Dashboard{}
	err := json.Unmarshal(mem, dash)
	if nil != err {
		panic(err)
	}
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

			switch t := c.Widget.(type) {
			case SL.SparkLine:
				SL.NewSparkLinex(&t, d.db)
				col := ui.NewCol(c.Span, c.Offset, ui.NewSparklines(*t.SL))
				col.Height = c.Height
				columns = append(columns, col)
			case P:
				par := ui.NewPar(t.Text)
				par.Border = t.Border
				par.Height = t.Height
				col := ui.NewCol(c.Span, c.Offset, par)
				columns = append(columns, col)
			default:
				log.Error("Invalid type in dashboard: ", t)

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

func (d *Dashboard) UpdateAll(time *TS.TimeSelect) {
}

func (d *Dashboard) GetGrid() *ui.Grid {
	return d.Grid
}
