package game_config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}

type mapBuildConf struct {
	Title  string `json:"title"`
	Cfg    []cfg  `json:"cfg"`
	cfgMap map[int8][]cfg
}

var MapBuildConf = &mapBuildConf{
	cfgMap: make(map[int8][]cfg),
}

const mapBuildConfFile = "/internal/logic/game_config/map_build.json"

func (m *mapBuildConf) Load() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + mapBuildConfFile

	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + mapBuildConfFile
		}
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		log.Println("json格式不正确，解析出错")
		panic(err)
	}
	for _, v := range m.Cfg {
		_, ok := m.cfgMap[v.Type]
		if !ok {
			m.cfgMap[v.Type] = make([]cfg, 0)
		}
		m.cfgMap[v.Type] = append(m.cfgMap[v.Type], v)
	}
}

func (m *mapBuildConf) BuildConfig(buildType int8, level int8) *cfg {
	cfgs := m.cfgMap[buildType]
	for _, v := range cfgs {
		if v.Level == level {
			return &v
		}
	}
	return nil
}
