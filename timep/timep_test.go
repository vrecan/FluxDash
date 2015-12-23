package timep

import (
	// "fmt"
	. "github.com/smartystreets/goconvey/convey"
	TS "github.com/vrecan/FluxDash/timeselect"
	"testing"
)

func TestTimeP(t *testing.T) {
	Convey("Init with valid timeselect", t, func() {
		time := TS.TimeSelect{}
		tp := NewTimeP(&TimeP{Height: 1})
		tp.Update(time)
		So(tp.Par.Text, ShouldEqual, "Time: 15m Interval: 5s Refresh: 5s")
		So(tp.Par.Height, ShouldEqual, 1)
	})
}
