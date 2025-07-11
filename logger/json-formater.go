package logger

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type JSONFormatter struct {
	TimestampFormat       string
	MsgFormatter          MessageFormater
	FunctionNameFormatter FunctionNameFormatter
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(map[string]interface{})
	if f.TimestampFormat != "" {
		data["timestamp"] = entry.Time.Format(f.TimestampFormat)
	} else {
		data["timestamp"] = entry.Time.Format("2006-01-02 15:04:05")
	}

	data["level"] = entry.Level.String()
	data["message"] = entry.Message

	if entry.Caller != nil {
		data["file"] = entry.Caller.File
		data["line"] = entry.Caller.Line
		if f.FunctionNameFormatter != nil {
			data["function"] = f.FunctionNameFormatter.Format(entry.Caller.Function)
		} else {
			data["function"] = entry.Caller.Function
		}
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var masked string
	if f.MsgFormatter != nil {
		masked = f.MsgFormatter.Format(string(serialized))
	} else {
		masked = string(serialized)
	}

	return []byte(masked + "\n"), nil
}
