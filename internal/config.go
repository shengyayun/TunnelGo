package internal

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Config struct {
	Tunnels map[string]Tunnel `yaml:"tunnels"`
}

// NewConfig 根据配置文件路径，生成配置对象
func NewConfig(path string) (*Config, error) {
	filename, _ := filepath.Abs(path)
	log.Printf("Load %s\n", filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
