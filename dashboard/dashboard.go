package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Dashboard struct {
	Lines SparkLinesData `json:"sparklines"`
}

type SparkLineData struct {
	Title    string `json:"title"`
	From     string `json:"from"`
	Time     string `json:"time"`
	Where    string `json:"where"`
	Height   int    `json:"height"`
	DataType string `json:"dataType"`
}

type SparkLinesData struct {
	SL []SparkLineData `json:"sparkline"`
}

// Time     string `json:"time"`
// Title    string `json:"title"`
// Where    string `json:"where"`
// dataType string `json:"dataType"`
func CreateExampleDash() {
	dash := ExampleDash()

	raw, err := json.Marshal(dash)
	fmt.Println("err: ", err)
	fmt.Println(string(raw))
}

//ExampleDash returns an example dashboard with all basic stuff filled out.
func ExampleDash() Dashboard {
	dash := Dashboard{}
	dash.Lines.SL = append(dash.Lines.SL, SparkLineData{From: "/system.cpu/", Time: "now - 15m", Title: "CPU", Where: "", DataType: "percent"})
	return dash
}

//Dashboard get dash from path.
func NewDashboard(f string) Dashboard {
	mem, e := ioutil.ReadFile(f)
	if e != nil {
		log.Fatal("File error: ", e)
	}
	fmt.Printf("%s\n", string(mem))

	// var jsontype jsonobject
	dash := Dashboard{}
	err := json.Unmarshal(mem, &dash)
	if nil != err {
		panic(err)
	}
	return dash

}
