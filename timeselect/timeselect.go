package timeselect

import (
	"fmt"
)

var times = [...]string{"15m", "30m", "1h", "8h", "24h", "2d", "5d", "7d"}
var intervals = [...]string{"5s", "10s", "20s", "160s", "8m", "24m", "32m", "45m"}
var refresh = [...]int{5, 10, 20, 40, 80, 160, 320, 640}

type TimeSelect struct {
	pos int
}

func (t *TimeSelect) CurTime() (string, string, int) {
	if t.pos > len(times)-1 {
		t.pos = 0
	}
	return fmt.Sprintf("now() - %s", times[t.pos]), fmt.Sprintf("GROUP BY time(%s)", intervals[t.pos]), refresh[t.pos]

}
func (t *TimeSelect) NextTime() (string, string, int) {
	t.pos++
	return t.CurTime()
}

func (t *TimeSelect) PrevTime() (string, string, int) {
	t.pos--
	if t.pos < 0 {
		t.pos = int(len(times) - 1)
	}
	return t.CurTime()
}

func (t *TimeSelect) DisplayTimes() (string, string, int) {
	if t.pos > len(times)-1 {
		t.pos = 0
	}
	return times[t.pos], intervals[t.pos], refresh[t.pos]
}
