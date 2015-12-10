package sparkline

import (
	"fmt"
	DBC "github.com/influxdb/influxdb/client/v2"
	. "github.com/smartystreets/goconvey/convey"
	DASH "github.com/vrecan/FluxDash/dashboard"
	DB "github.com/vrecan/FluxDash/influx"
	"testing"
)

func TestSparklines(t *testing.T) {
	Convey("Build query with every option", t, func() {
		q := buildQuery("mean(value)", "/thing/", `"service='woo'`, "5s", " GROUP BY 'service'")
		So(q, ShouldEqual, `SELECT mean(value) FROM /thing/ WHERE "service='woo' AND time > 5s  GROUP BY 'service'`)
	})

	Convey("Build query with out where", t, func() {
		q := buildQuery("mean(value)", "/thing/", "", "5s", " GROUP BY 'service'")
		So(q, ShouldEqual, `SELECT mean(value) FROM /thing/ WHERE time > 5s  GROUP BY 'service'`)
	})

	Convey("Build Sparklines from data", t, func() {
		d := DASH.ExampleDash()
		db, err := DB.NewInflux(DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"})
		So(err, ShouldBeNil)
		sparks := NewSparkLinesFromData(db, d.Lines)
		So(len(sparks.lines), ShouldBeGreaterThanOrEqualTo, 1)
	})

	Convey("Build individual Sparkline from data", t, func() {
		d := DASH.ExampleDash()
		db, err := DB.NewInflux(DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"})
		So(err, ShouldBeNil)
		sparks := NewSparkLineFromData(db, d.Lines.SL...)
		So(len(sparks), ShouldBeGreaterThanOrEqualTo, 1)
		fmt.Println("woo")
	})

}
