package query

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	DB "github.com/vrecan/FluxDash/influx"
	"time"
)

func Build(sel string, from string, where string, time string, groupBy string) string {
	if len(sel) == 0 || len(from) == 0 || len(time) == 0 {
		log.Critical("invalid query string :", fmt.Sprintf("SELECT %s FROM %s WHERE %s AND time > %s %s fill(0)", sel, from, where, groupBy))
		return ""
	}
	if len(where) > 0 {
		return fmt.Sprintf("SELECT %s FROM %s WHERE %s AND time > %s %s fill(0)", sel, from, where, time, groupBy)
	} else {
		return fmt.Sprintf("SELECT %s FROM %s WHERE time > %s %s fill(0)", sel, from, time, groupBy)
	}
}

func GetIntData(db DB.DBI, q string) (data []int) {
	r, err := db.Query(q)
	if nil != err {
		log.Error("No data from query:", q)
		return data
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
		log.Error("No data from query:", q)
		return data
	}

	for _, row := range r[0].Series[0].Values {
		_, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Critical(err)
			return data
		}
		if len(row) > 1 {
			if nil != row[1] {
				val, err := row[1].(json.Number).Float64()
				if nil != err {
					log.Error("Failed to parse data: ", err)
				} else {
					data = append(data, int(val))
				}
			}
		}

	}
	return data
}

func GetIntDataFromTags(db DB.DBI, q string) (data [][]int, labels []string) {
	r, err := db.Query(q)
	if nil != err {
		log.Error("No data from query:", q)
		return data, labels
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
		log.Error("No data from query:", q)
		return data, labels
	}
	labels = make([]string, len(r[0].Series))
	data = make([][]int, len(r[0].Series))
	for i, result := range r[0].Series {
		labels[i] = result.Name
		for _, row := range result.Values {
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Critical(err)
				return data, labels
			}
			if len(row) > 1 {
				if nil != row[1] {
					val, err := row[1].(json.Number).Float64()
					if nil != err {
						log.Error("Failed to parse data: ", err)
					} else {
						data[i] = append(data[i], int(val))
					}
				}
			}
		}
	}
	return data, labels
}

func GetDataForBar(db DB.DBI, q string) (data []int, labels [][]string) {
	r, err := db.Query(q)
	if nil != err {
		log.Error("No data from query:", q)
		return data, labels
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
		log.Error("No data from query:", q)
		return data, labels
	}
	labels = make([][]string, len(r[0].Series))
	for i, result := range r[0].Series {
		series := fmt.Sprintf("S%d", i)
		labels[i] = []string{series, result.Name}
		for _, row := range result.Values {
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Critical(err)
				return data, labels
			}
			if len(row) > 1 {
				if nil != row[1] {
					val, err := row[1].(json.Number).Float64()
					if nil != err {
						log.Error("Failed to parse data: ", err)
					} else {
						data = append(data, int(val))
					}
				}
			}

		}
	}

	return data, labels
}
