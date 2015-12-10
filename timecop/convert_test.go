package timecop

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConvert(t *testing.T) {

	Convey("return case", t, func() {
		time, unit, err := GetRoundedTime(1, "seconds")
		So(err, ShouldBeNil)
		So(time, ShouldEqual, 1)
		So(unit, ShouldEqual, "seconds")
		time, unit, err = GetRoundedTime(59, "seconds")
		So(err, ShouldBeNil)
		So(time, ShouldEqual, 59)
		So(unit, ShouldEqual, "seconds")
	})
	Convey("next case", t, func() {
		time, unit, err := GetRoundedTime(60, "seconds")
		So(err, ShouldBeNil)
		So(time, ShouldEqual, 1)
		So(unit, ShouldEqual, "minutes")
	})
	Convey("prev case", t, func() {
		time, unit, err := GetRoundedTime(.1, "seconds")
		So(err, ShouldBeNil)
		So(time, ShouldEqual, 100)
		So(unit, ShouldEqual, "milliseconds")
	})
}
