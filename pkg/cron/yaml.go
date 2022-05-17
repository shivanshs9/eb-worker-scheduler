package cron

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type CronFile struct {
	Version int         `yaml:"version"`
	Crons   []CronEvent `yaml:"cron"`
}

func ParseYaml(filePath string) (*CronFile, error) {
	cronFile := &CronFile{}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cronFile, err
	}
	if err = yaml.Unmarshal(data, cronFile); err != nil {
		return cronFile, err
	}
	return cronFile, nil
}
