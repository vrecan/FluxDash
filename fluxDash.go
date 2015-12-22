package main

import (
	log "github.com/cihub/seelog"
	DBC "github.com/influxdb/influxdb/client/v2"
	DASH "github.com/vrecan/FluxDash/dashboards"
	DB "github.com/vrecan/FluxDash/influx"
	// "os"
	// "fmt"
)

func main() {

	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		log.Warn("Failed to load config", err)
	}

	log.ReplaceLogger(logger)

	// DASH.CreateExampleDash()
	// fmt.Println(string(dash))
	// os.Exit(0)

	c := DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"}
	db, err := DB.NewInflux(c)
	if nil != err {
		panic(err)
	}

	// DASH.CreateExampleDash()

	// dash := DASH.ExampleDash(db)
	dash := DASH.NewDashboardFromFile(db, "dashboards/example.json")
	d := DASH.NewMonitor(dash)
	d.Start()

}
