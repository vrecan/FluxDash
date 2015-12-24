package merge

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type A struct {
	Title string
	Users []string
}

type B struct {
	Title string
	Users []string
}

type D struct {
	Title string
	Other string
}

type E struct {
	Title int
}

func TestMerge(t *testing.T) {
	Convey("Test merging fields with same names", t, func() {
		a := A{Title: "title", Users: []string{"user"}}
		b := B{}
		r := Merge(a, &b).(*B)
		So(r.Title, ShouldEqual, "title")
		So(r.Users[0], ShouldEqual, "user")
	})

	Convey("Test merging fields with same names ignore title", t, func() {
		a := A{Title: "title", Users: []string{"user"}}
		b := B{}
		r := Merge(a, &b, "Title").(*B)
		So(r.Title, ShouldBeEmpty)
		So(r.Users[0], ShouldEqual, "user")
	})

	Convey("Test merging fields with some fields overlapping but others not", t, func() {
		a := A{Title: "title", Users: []string{"user"}}
		d := D{}
		r := Merge(a, &d).(*D)
		So(r.Title, ShouldEqual, "title")
		So(r.Other, ShouldBeEmpty)
	})

	Convey("Test merging fields with samefield but different type", t, func() {
		a := A{Title: "title", Users: []string{"user"}}
		e := E{}
		r := Merge(a, &e).(*E)
		So(r.Title, ShouldEqual, 0)
	})
}
