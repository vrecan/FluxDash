package influx

import (
	DB "github.com/influxdata/influxdb/client/v2"
)

type DBI interface {
	Query(string) ([]DB.Result, error)
	Close() error
}

//Influx is an influxdb wrapper to alow simpler querying.
type Influx struct {
	CLI  DB.Client
	Conf DB.HTTPConfig
}

//Create influx db connection
func NewInflux(conf DB.HTTPConfig) (db *Influx, err error) {
	db = &Influx{Conf: conf}
	cli, err := DB.NewHTTPClient(conf)
	db.CLI = cli
	return db, err

}

// queryDB convenience function to query the database
func (i *Influx) Query(cmd string) (res []DB.Result, err error) {
	q := DB.Query{
		Command:  cmd,
		Database: "stats",
	}
	if response, err := i.CLI.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}

func (i *Influx) Close() error {
	if nil != i {
		return i.CLI.Close()
	}
	return nil
}
