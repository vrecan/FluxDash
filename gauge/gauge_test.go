package gauge

import (
	. "github.com/smartystreets/goconvey/convey"
	M "github.com/vrecan/FluxDash/mockdb"
	TS "github.com/vrecan/FluxDash/timeselect"
	"testing"
)

func TestGauge(t *testing.T) {
	Convey("Init Gauge", t, func() {
		time := TS.TimeSelect{}
		db := &M.MockDB{}
		gauge := &Gauge{Height: 1}
		gauge = NewGauge(db, gauge)

		gauge.Update(time)
		So(gauge.G.Percent, ShouldEqual, 0)
		So(gauge.G.Height, ShouldEqual, 1)
	})
}
