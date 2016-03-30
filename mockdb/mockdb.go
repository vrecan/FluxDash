package mockdb

import (
	DB "github.com/influxdata/influxdb/client/v2"
)

type MockDB struct {
}

func (i *MockDB) Query(cmd string) (res []DB.Result, err error) {
	return res, nil
}

func (i *MockDB) Close() error {
	return nil
}
