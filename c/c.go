package c

import (
	"fmt"
	DB "github.com/influxdb/influxdb/client"
	VIPER "github.com/spf13/viper"
)

//Simple configuration struct
type FluxConf struct {
	DB DB.ClientConfig
}

func GetConf(path string) (conf FluxConf, err error) {
	v := VIPER.New()
	v.SetConfigName("flux")
	v.AddConfigPath(path)
	err = v.ReadInConfig()
	if nil != err {
		return conf, err
	}
	conf = FluxConf{}
	err = v.Marshal(&conf)
	fmt.Println(conf)

	if nil != err {
		return conf, err
	}
	return conf, err
}
