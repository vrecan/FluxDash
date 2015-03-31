package main

import (
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
	"math"
	"time"
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	sinps := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()
	lc := ui.NewLineChart()
	lc.Border.Label = "System CPU"
	lc.Data = sinps
	lc.Width = 200
	lc.Height = 15
	lc.X = 0
	lc.Y = 14
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorBlue | ui.AttrBold
	lc.DataLabels = make([]string, 0)
	lc.DataLabels = append(lc.DataLabels, "woo")
	lc.DataLabels = append(lc.DataLabels, "Somethign else")
	lc.Mode = "dot"
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, lc)))
	ui.Body.Align()

	draw := func(t int) {
		lc.Data = sinps[t/2:]
		ui.Render(lc)
	}

	evt := make(chan tm.Event)
	go func() {
		for {
			evt <- tm.PollEvent()
		}
	}()
	i := 0
	draw(i)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case e := <-evt:
			if e.Type == tm.EventKey && e.Ch == 'q' {
				return
			}
		case <-ticker.C:
			i++
			draw(i)
		}
	}
}
