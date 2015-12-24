package main

import (
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	DBC "github.com/influxdb/influxdb/client/v2"
	DASH "github.com/vrecan/FluxDash/dashboards"
	DB "github.com/vrecan/FluxDash/influx"
	"os"
	FP "path/filepath"
)

var json = flag.String("f", "dashboards", "-f can be passed a single file or a folder of json dashboards")

func main() {
	flag.Parse()
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		log.Warn("Failed to load config", err)
	}

	log.ReplaceLogger(logger)
	c := DBC.HTTPConfig{Addr: "http://127.0.0.1:8086", Username: "admin", Password: "logrhythm!1"}
	db, err := DB.NewInflux(c)
	if nil != err {
		panic(err)
	}
	defer db.Close()
	dashboards := GetDashbordsFromFlag(db)

	if len(dashboards) <= 0 {
		fmt.Println("No valid dashboards found in path: ", *json)
		os.Exit(2)
	}
	d := DASH.NewMonitor(dashboards...)
	d.Start()

}

//GetDashboardsFromFlag parses the flag and gets all the json dashboards from the path supplied.
func GetDashbordsFromFlag(db DB.DBI) (dashboards []DASH.Stats) {
	path := FP.Clean(*json)
	info, err := os.Stat(path)
	if nil != err {
		fmt.Println("Failed to parse json file/directory path: ", path)
		os.Exit(1)
	}
	if info.IsDir() {
		globPattern := fmt.Sprintf("%s%s*.json", path, string(FP.Separator))
		files, err := FP.Glob(globPattern)
		fmt.Println(files)
		if nil != err {
			log.Error("Failed parsing json files in path: ", path, " with error: ", err)
			return dashboards
		}
		for _, f := range files {
			dashboards = append(dashboards, DASH.NewDashboardFromFile(db, f))
		}

	} else {
		//parse as a single dashboard file
		dashboards = append(dashboards, DASH.NewDashboardFromFile(db, path))

	}
	return dashboards
}
