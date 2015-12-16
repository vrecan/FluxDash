package dashboards

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
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
	Grid *ui.Grid       `json:"-` //json:"-" omits a field from being encoded
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
	sl1 := &SL.SparkLine{Title: "CPU", From: "/system.cpu/", Where: "", Height: 1, DataType: 1}
	sl2 := &SL.SparkLine{Title: "Dispatch GC", Height: 1, From: "/gc.pause.ns/", Where: `"service"= 'godispatch'`, DataType: 1}
	sparks := make([]*SL.SparkLine, 0)
	sparks = append(sparks, sl1, sl2)
	sparkLines := &Column{Height: 1, Span: 12, Widget: SL.SparkLines{BorderLabel: "System", Border: true, Lines: sparks}}
	columns := make([]*Column, 0)
	columns = append(columns, sparkLines)
	dash.Rows = append(dash.Rows, &Row{Height: 1, Span: 12, Offset: 0, Columns: columns})
	p1 := &Column{Height: 1, Span: 6, Offset: 0, Widget: P{Text: "Static text is all the rage!!!", Height: 3, Border: true}}
	columns2 := make([]*Column, 0)
	columns2 = append(columns2, p1)
	dash.Rows = append(dash.Rows, &Row{Height: 1, Span: 12, Offset: 0, Columns: columns2})
	ptRow := &Column{Height: 1, Span: 6, Offset: 0, Widget: TP.TimeP{Height: 3, Border: true}}
	columns3 := make([]*Column, 0)
	columns3 = append(columns3, ptRow)

	dash.Rows = append(dash.Rows, &Row{Height: 1, Span: 12, Offset: 0, Columns: columns3})

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
			// case SL.SparkLine:
			// 	SL.NewSparkLinex(&t, d.db)
			// 	col := ui.NewCol(c.Span, c.Offset, ui.NewSparklines(*t.SL))
			// 	col.Height = c.Height
			// 	columns = append(columns, col)
			// 	c.Widget = t
			case SL.SparkLines:
				c.Widget = SL.NewSparkLines(d.db, &t)
				col := ui.NewCol(c.Span, c.Offset, t.SL)
				columns = append(columns, col)
			case P:
				par := ui.NewPar(t.Text)
				par.Border = t.Border
				par.Height = t.Height
				t.Par = par
				col := ui.NewCol(c.Span, c.Offset, par)
				columns = append(columns, col)
				c.Widget = t
			case TP.TimeP:
				c.Widget = TP.NewTimeP(&t)
				col := ui.NewCol(c.Span, c.Offset, t.Par)
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
	for _, r := range d.Rows {
		for _, c := range r.Columns {
			switch t := c.Widget.(type) {
			case P:
				continue //ignore static p tags
			case *TP.TimeP:
				t.Update(*time)
			case *SL.SparkLines:
				log.Info("SL UPDATE:", t)
				t.Update(*time)
			}
		}
	}
	d.Grid.Align()
	ui.Render(d.Grid)
}

func (d *Dashboard) GetGrid() *ui.Grid {
	return d.Grid
}
