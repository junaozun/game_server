package config

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"path"

	"github.com/a8m/envsubst"
	"gopkg.in/yaml.v2"
)

type IConfigSource interface {
	Load(v interface{}) error
}

// FileConfigSource 数据源是本地文件
type FileConfigSource struct {
	FilePath string // 文件路径
}

func (f *FileConfigSource) Load(v interface{}) error {
	contents, err := ioutil.ReadFile(f.FilePath)
	if err != nil {
		return err
	}

	ext := path.Ext(f.FilePath)
	buf, err := envsubst.Bytes(contents)
	return unmarshal(ext, buf, v)
}

// LoadConfigFromFile 从本地文件加载配置
func LoadConfigFromFile(filePath string, configValue interface{}) error {
	return LoadConfigFromSource(&FileConfigSource{
		FilePath: filePath,
	}, configValue)
}

// LoadConfigFromSource 加载配置
func LoadConfigFromSource(ics IConfigSource, configValue interface{}) error {
	return ics.Load(configValue)
}

func unmarshal(ext string, buf []byte, v interface{}) error {
	switch ext {
	case ".xml":
		if err := xml.Unmarshal(buf, v); err != nil {
			return err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(buf, v); err != nil {
			return err
		}
	case ".json":
		if err := json.Unmarshal(buf, v); err != nil {
			return err
		}
	default:
		if err := yaml.Unmarshal(buf, v); err != nil {
			return err
		}
	}
	return nil
}
