package timeselect

import (
	"fmt"
)

var times = [...]string{"15m", "30m", "1h", "8h", "24h", "2d", "5d", "7d"}
var intervals = [...]string{"5s", "10s", "20s", "160s", "8m", "24m", "45m", "1h"}

type TimeSelect struct {
	pos int
}

func (t *TimeSelect) CurTime() (string, string) {
	if t.pos > len(times)-1 {
		t.pos = 0
	}
	return fmt.Sprintf("now() - %s", times[t.pos]), fmt.Sprintf("GROUP BY time(%s)", intervals[t.pos])

}
func (t *TimeSelect) NextTime() (string, string) {
	t.pos++
	return t.CurTime()
}

func (t *TimeSelect) PrevTime() (string, string) {
	t.pos--
	if t.pos < 0 {
		t.pos = int(len(times) - 1)
	}
	return t.CurTime()
}

func (t *TimeSelect) DisplayTimes() (string, string) {
	if t.pos > len(times)-1 {
		t.pos = 0
	}
	return times[t.pos], intervals[t.pos]
}
