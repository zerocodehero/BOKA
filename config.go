/*
 @author: lynn
 @date: 2023/5/3
 @time: 21:27
*/

package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

const ConfigFile = "./config/config.yaml"

var (
	CONFIG *Config
)

func setDefault(v *viper.Viper) {
	v.SetDefault("custId", "")
	v.SetDefault("compId", "")
	v.SetDefault("userName", "")
	v.SetDefault("passWord", "")
	v.SetDefault("source", "")
	v.SetDefault("sec", 0)
}

type Config struct {
	BoKa KaConfig `yaml:"boka" mapstructure:"boka" json:"boka"`
}

type KaConfig struct {
	CustId   string `json:"custId" yaml:"custId" mapstructure:"custId" `
	CompId   string `json:"compId" yaml:"compId" mapstructure:"compId" `
	UserName string `json:"userName" yaml:"userName" mapstructure:"userName" `
	PassWord string `json:"passWord" yaml:"passWord" mapstructure:"passWord" `
	Source   string `json:"source" yaml:"source" mapstructure:"source" `
	Sec      int64  `mapstructure:"sec" json:"sec" yaml:"sec"`
}

func InitConf(path ...string) {

	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "boka", "", "choose conf file.")
		flag.Parse()
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv("BOKACONFIG"); configEnv == "" {
				config = ConfigFile
				log.Printf("您正在使用config的默认值,config的路径为%v\n", config)
			} else {
				config = configEnv
				log.Printf("您正在使用configEnv环境变量,config的路径为%v\n", config)
			}
		} else {
			log.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
		}
	} else {
		config = path[0]
		log.Printf("您正在使用InitConf传递的值,config的路径为%v\n", config)
	}
	v := viper.New()
	// 设置默认值
	setDefault(v)

	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error conf file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("conf file changed:", e.Name)
		if err = v.Unmarshal(&CONFIG); err != nil {
			log.Println(err)
		}
	})
	if err = v.Unmarshal(&CONFIG); err != nil {
		log.Println(err)
	}
}
