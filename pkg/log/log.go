package log

import (
	"log"
	"time"
)

type Config struct {
	Level      string `mapstructure:"level"`      // e.g., info, debug, error
	Format     string `mapstructure:"format"`     // e.g., json, pretty
	OutputPath string `mapstructure:"outputPath"` // e.g., stdout, stderr, or file path
}

type LogCustom struct{}

type LogData struct {
	Err         error
	Description string
	StartTime   interface{} // pakai time.Time kalau mau waktu mulai
	Response    interface{}
}

func NewLogger() *LogCustom {
	return &LogCustom{}
}

func (l *LogCustom) Info(data LogData) {
	duration := ""
	if t, ok := data.StartTime.(time.Time); ok {
		duration = time.Since(t).String()
	}
	log.Printf("[INFO] %s | err: %v | duration: %s | response: %v\n", data.Description, data.Err, duration, data.Response)
}

func (l *LogCustom) Error(data LogData) {
	log.Printf("[ERROR] %s | err: %v\n", data.Description, data.Err)
}
