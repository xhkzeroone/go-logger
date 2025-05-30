package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"regexp"
	"strings"
)

type FunctionNameFormatter interface {
	Format(fullName string) string
}

type DefaultFunctionNameFormatter struct{}

func (f *DefaultFunctionNameFormatter) Format(fullName string) string {
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullName
}

type MessageFormater interface {
	Format(message string) string
}

type DefaultMessageFormater struct {
}

func (d *DefaultMessageFormater) Format(message string) string {
	return message
}

type DynamicFormatter struct {
	Pattern               string
	TimestampFormat       string
	MsgFormatter          MessageFormater
	FunctionNameFormatter FunctionNameFormatter
}

func (f *DynamicFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())

	message := f.MsgFormatter.Format(entry.Message)

	file := "???"
	line := 0
	function := "???"
	if entry.Caller != nil {
		file = path.Base(entry.Caller.File)
		line = entry.Caller.Line
		function = f.FunctionNameFormatter.Format(entry.Caller.Function)
	}

	out := f.Pattern
	out = strings.ReplaceAll(out, "%timestamp%", timestamp)
	out = strings.ReplaceAll(out, "%level%", level)
	out = strings.ReplaceAll(out, "%file%", file)
	out = strings.ReplaceAll(out, "%line%", fmt.Sprintf("%d", line))
	out = strings.ReplaceAll(out, "%function%", function)
	out = strings.ReplaceAll(out, "%message%", message)
	for _, k := range extractPlaceholders(f.Pattern) {
		placeholder := "%" + k + "%"
		value, ok := entry.Data[k]
		if !ok || value == nil {
			out = strings.ReplaceAll(out, placeholder, "null")
		} else {
			out = strings.ReplaceAll(out, placeholder, fmt.Sprint(value))
		}
	}

	return []byte(out + "\n"), nil
}

func extractPlaceholders(pattern string) []string {
	re := regexp.MustCompile(`%([a-zA-Z0-9_]+)%`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	var keys []string
	for _, match := range matches {
		if len(match) > 1 {
			keys = append(keys, match[1])
		}
	}
	return keys
}
