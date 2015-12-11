package main

import (
	DBC "github.com/influxdb/influxdb/client/v2"
	DASH "github.com/vrecan/FluxDash/dashboards"
	DB "github.com/vrecan/FluxDash/influx"
)

func main() {
	c := DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"}
	db, err := DB.NewInflux(c)
	if nil != err {
		panic(err)
	}
	sys := DASH.NewSystem(db)
	ins := DASH.NewInsert(db)
	d := DASH.NewMonitor(sys, ins)
	d.Start()

}
