package logger

import (
	"encoding/xml"
	"os"
)

type LogConfig struct {
	TimestampFormat string `xml:"timestampFormat"`
	Pattern         string `xml:"pattern"`
	Level           string `xml:"level"`
}

func LoadLogConfig(path string) (*LogConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg LogConfig
	err = xml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
