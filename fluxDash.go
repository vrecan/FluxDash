package main

import (
	log "github.com/cihub/seelog"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
	C "github.com/vrecan/FluxDash/c"
	DB "github.com/vrecan/FluxDash/influx"
	SPARK "github.com/vrecan/FluxDash/spark"
	"time"
)

func main() {
	c, err := C.GetConf(".")
	if nil != err {
		log.Critical(err)
	}

	db, err := DB.NewInflux(&c.DB)

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	var sparks *SPARK.Sparks

	draw := func() {
		sparks, err = SPARK.NewSparks(SPARK.SparksConf{
			Query:  "select mean(value) from /default\\.localhost\\.vitals\\.system\\.cpu\\..*/ where time > now() - 30m group by time(5s) fill(0) order asc",
			Title:  "WOOO",
			Width:  50,
			Height: 50,
		}, db)
		if nil != err {
			log.Error(err)
		}
		ui.Render(sparks.Render())
	}

	draw()
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, sparks.Render())))
	ui.Body.Align()

	evt := make(chan tm.Event)
	go func() {
		for {
			evt <- tm.PollEvent()
		}
	}()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case e := <-evt:
			if e.Type == tm.EventKey && e.Ch == 'q' {
				return
			}
		case <-ticker.C:
			draw()
		}
	}
}
