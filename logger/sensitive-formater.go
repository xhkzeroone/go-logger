package logger

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Patterns struct {
	XMLName xml.Name  `xml:"patterns"`
	Rules   []Pattern `xml:"pattern"`
}

type Pattern struct {
	Type        string `xml:"type"`
	Regex       string `xml:"regex"`
	Replacement string `xml:"replacement"`
}

func loadPatterns(path string) ([]Pattern, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var patterns Patterns
	err = xml.Unmarshal(data, &patterns)
	if err != nil {
		return nil, err
	}

	return patterns.Rules, nil
}

func sensitiveMessage(message string, patterns []Pattern) string {
	sanitized := message
	for _, p := range patterns {
		re := regexp.MustCompile(p.Regex)
		sanitized = re.ReplaceAllString(sanitized, p.Replacement)
	}
	return sanitized
}

type SensitiveMessageFormater struct {
	patterns []Pattern
}

func NewSensitiveMessageFormater() MessageFormater {
	dir, _ := os.Getwd()
	patterns, err := loadPatterns(filepath.Join(dir, "/sensitive-patterns.xml"))
	if err != nil {
		fmt.Println("Failed to load patterns:", err)
		return &DefaultMessageFormater{}
	}
	return &SensitiveMessageFormater{
		patterns: patterns,
	}
}

func RegisterSensitiveMessageFormater() {
	RegisterMessageFormater(NewSensitiveMessageFormater())
}

func (d *SensitiveMessageFormater) Format(message string) string {
	return sensitiveMessage(message, d.patterns)
}
