package timeselect

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTimeSelect(t *testing.T) {

	Convey("get Current Time", t, func() {
		t := TimeSelect{}
		res, i, ref := t.CurTime()
		So(res, ShouldEqual, fmt.Sprintf("now() - %s", times[0]))
		So(i, ShouldEqual, fmt.Sprintf("GROUP BY time(%s)", intervals[0]))
		So(ref, ShouldEqual, refresh[0])
	})

	Convey("get Prev Time starting at 0", t, func() {
		t := TimeSelect{}
		res, i, ref := t.PrevTime()
		So(res, ShouldEqual, fmt.Sprintf("now() - %s", times[len(times)-1]))
		So(i, ShouldEqual, fmt.Sprintf("GROUP BY time(%s)", intervals[len(intervals)-1]))
		So(ref, ShouldEqual, refresh[len(refresh)-1])
	})

	Convey("get Next Time starting at max", t, func() {
		t := TimeSelect{}
		t.pos = len(times) - 1
		res, i, ref := t.NextTime()
		So(res, ShouldEqual, fmt.Sprintf("now() - %s", times[0]))
		So(i, ShouldEqual, fmt.Sprintf("GROUP BY time(%s)", intervals[0]))
		So(ref, ShouldEqual, refresh[0])
	})

	Convey("get Next Time starting at max", t, func() {
		t := TimeSelect{}
		time, interval, r := t.DisplayTimes()
		So(time, ShouldEqual, times[0])
		So(interval, ShouldEqual, intervals[0])
		So(r, ShouldEqual, refresh[0])
	})
}
