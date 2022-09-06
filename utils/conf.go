package utils

import (
	"fmt"
	yaml "github.com/asim/go-micro/plugins/config/encoder/yaml/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
	"sync"
)

var conf config.Config
var inst sync.Once

//单例获取配置
func Instance() config.Config {
	inst.Do(func() {
		load()
	})
	return conf
}

//动态获取yaml的配置文件
func load() {
	//读取动态配置
	enc := yaml.NewEncoder()
	conf, _ = config.NewConfig(
		config.WithReader(
			json.NewReader( // json reader for internal config merge
				reader.WithEncoder(enc),
			),
		),
	)

	if err := conf.Load(file.NewSource(file.WithPath("../config.yaml"))); err != nil {
		fmt.Println(err.Error())
		return
	}

}
