package sparkline

import (
	// "fmt"
	. "github.com/smartystreets/goconvey/convey"
	M "github.com/vrecan/FluxDash/mockdb"
	TS "github.com/vrecan/FluxDash/timeselect"
	"testing"
)

func TestSparkline(t *testing.T) {
	Convey("Init Sparkline", t, func() {
		time := TS.TimeSelect{}
		sl := &SparkLine{Height: 1}
		lines := make([]*SparkLine, 0)
		lines = append(lines, sl)
		db := &M.MockDB{}
		spark := NewSparkLines(db, &SparkLines{Height: 3, Lines: lines})
		spark.Update(time)
		So(spark.Lines[0].Height, ShouldEqual, 1)
		So(spark.SL.Height, ShouldEqual, 5)
	})
}
