package dashboards

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGauge(t *testing.T) {
	Convey("Init debouncer", t, func() {
		timeDebounceChan := make(chan interface{}, 1000)
		TimeChangeChan := make(chan interface{}, 1)
		go DebounceChan(timeDebounceChan, 13*time.Millisecond, TimeChangeChan)

		ctr := 0

		for i := 0; i < 10; i++ {
			ctr = i
			timeDebounceChan <- ctr
		}
		r := <-TimeChangeChan
		So(r, ShouldEqual, ctr)
		close(TimeChangeChan)
	})

}
