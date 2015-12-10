package timecop

import (
	"errors"
	"fmt"
)

//UnitInfo that describes the bounding box, previous and subsequent unit
type UnitInfo struct {
	MinRounded float64
	MaxRounded float64
	Next       string
	NextRatio  float64
	Prev       string
	PrevRatio  float64
}

//The defined units
var Units map[string]*UnitInfo

func init() {
	Units = make(map[string]*UnitInfo)
	Units["nanoseconds"] = &UnitInfo{MinRounded: 0, MaxRounded: 999, NextRatio: 1000, PrevRatio: 0, Prev: "", Next: "milliseconds"}
	Units["milliseconds"] = &UnitInfo{MinRounded: 1, MaxRounded: 999, NextRatio: 1000, PrevRatio: 1000, Prev: "nanoseconds", Next: "seconds"}
	Units["seconds"] = &UnitInfo{MinRounded: 1, MaxRounded: 59, NextRatio: 60, PrevRatio: 1000, Prev: "milliseconds", Next: "minutes"}
	Units["minutes"] = &UnitInfo{MinRounded: 1, MaxRounded: 59, NextRatio: 60, PrevRatio: 60, Prev: "seconds", Next: "hours"}
	Units["hours"] = &UnitInfo{MinRounded: 1, MaxRounded: 23, NextRatio: 24, PrevRatio: 60, Prev: "minutes", Next: "days"}
	Units["days"] = &UnitInfo{MinRounded: 1, MaxRounded: 364, NextRatio: 365, PrevRatio: 24, Prev: "hours", Next: "years"}
	Units["years"] = &UnitInfo{MinRounded: 1, MaxRounded: 99, NextRatio: 100, PrevRatio: 365, Prev: "days", Next: "centuries"}
	Units["centuries"] = &UnitInfo{MinRounded: 1, MaxRounded: 9, NextRatio: 10, PrevRatio: 100, Prev: "years", Next: "millenia"}
	Units["millenia"] = &UnitInfo{MinRounded: 1, MaxRounded: 0, NextRatio: 0, PrevRatio: 10, Prev: "centuries", Next: ""}
}

//GetRoundedTime with best possible unit
func GetRoundedTime(time float64, unit string) (newtime float64, newunit string, err error) {

	currentUnits, ok := Units[unit]
	if !ok {
		return time, unit, errors.New("invalid time unit")
	}
	if time < 0 {
		return time, unit, errors.New("Invalid negative time")
	}
	if time < currentUnits.MinRounded {
		return GetRoundedTime(time*currentUnits.PrevRatio, currentUnits.Prev)
	}

	if time > currentUnits.MaxRounded {
		if currentUnits.MaxRounded == 0 {
			return time, unit, nil
		}
		return GetRoundedTime(time/currentUnits.NextRatio, currentUnits.Next)
	}
	return time, unit, nil
}

//GetCommaString that represents the value and unit
func GetCommaString(time float64, unit string) string {
	newTime, newUnit, err := GetRoundedTime(time, unit)
	if err != nil {
		return "NaN"
	}
	return fmt.Sprintf("%d %s", int64(newTime), newUnit)
}
