package timep

import (
	"fmt"
	// log "github.com/cihub/seelog"
	ui "github.com/gizak/termui"
	TS "github.com/vrecan/FluxDash/timeselect"
)

type TimeP struct {
	Par    *ui.Par `json:"-"`
	Height int     `json:"height"`
	Border bool    `json:"border"`
}

func NewTimeP(t *TimeP) *TimeP {
	par := ui.NewPar("")
	par.Border = t.Border
	par.Height = t.Height
	t.Par = par
	return t
}

func (t *TimeP) Update(time TS.TimeSelect) {
	dt, di, dr := time.DisplayTimes()
	displayTimes := fmt.Sprintf("Time: %s Interval: %s Refresh: %vs", dt, di, dr)
	t.Par.Text = displayTimes
}
