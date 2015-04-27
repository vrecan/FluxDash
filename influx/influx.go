package influx

import (
	DB "github.com/influxdb/influxdb/client"
)

type Influx struct {
	CLI  *DB.Client
	Conf *DB.ClientConfig
}

//Create influx db connection
func NewInflux(conf *DB.ClientConfig) (db *Influx, err error) {
	cli, err := DB.NewClient(conf)
	if nil != err {
		return db, err
	}
	db = &Influx{
		CLI:  cli,
		Conf: conf,
	}
	return db, err
}

func (i *Influx) Query(q string) (interface{}, error) {
	return i.CLI.Query(q)

}
