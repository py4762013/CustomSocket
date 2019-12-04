package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	//配置文件路径
	configPath string
)

type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

func init()  {
	home, _ := homedir.Dir()
	// 默认配置文件名称
	configFileName := ".customsockets.json"
	//如果用户配置文件，就使用用户传入的配置文件
	if len(os.Args) == 2 {
		configFileName = os.Args[1]
	}
	configPath = path.Join(home, configFileName)
}

//保存配置到配置文件
func (config *Config) SaveConfig() {
	configJson, _ := json.MarshalIndent(config, "", "     ")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Errorf("保存到配置文件 %s 出错: %s", configPath, err)
	}
	log.Printf("保存到配置文件 %s 成功\n", configPath)
}

func (config *Config) ReadConfig() {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 config {
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("从文件 %s 中读取配置\n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("打开配置文件 %s 出错 %s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("格式不是合法的 JSON 配置文件:\n%s", file.Name())
		}
	}
}