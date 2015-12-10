package timeselect

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTimeSelect(t *testing.T) {

	Convey("get Current Time", t, func() {
		t := TimeSelect{}
		r, i := t.CurTime()
		So(r, ShouldEqual, fmt.Sprintf("now() - %s", times[0]))
		So(i, ShouldEqual, fmt.Sprintf("GROUP BY time(%s)", intervals[0]))
	})

	Convey("get Prev Time starting at 0", t, func() {
		t := TimeSelect{}
		r, i := t.PrevTime()
		So(r, ShouldEqual, fmt.Sprintf("now() - %s", times[len(times)-1]))
		So(i, ShouldEqual, fmt.Sprintf("GROUP BY time(%s)", intervals[len(intervals)-1]))
	})

	Convey("get Next Time starting at max", t, func() {
		t := TimeSelect{}
		t.pos = len(times) - 1
		r, i := t.NextTime()
		So(r, ShouldEqual, fmt.Sprintf("now() - %s", times[0]))
		So(i, ShouldEqual, fmt.Sprintf("GROUP BY time(%s)", intervals[0]))
	})
}
