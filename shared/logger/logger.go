package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
	mu     sync.Mutex
}

type Config struct {
	Level        LogLevel
	OutputType   string // "console", "file", "custom", or "multi"
	FilePath     string
	CustomWriter io.Writer
}

func New(cfg Config) (*Logger, error) {
	var output io.Writer

	switch cfg.OutputType {
	case "console":
		output = os.Stdout
	case "file":
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		output = file
	case "custom":
		if cfg.CustomWriter == nil {
			return nil, fmt.Errorf("custom writer is nil")
		}
		output = cfg.CustomWriter
	case "multi":
		writers := []io.Writer{os.Stdout}
		if cfg.FilePath != "" {
			file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file: %v", err)
			}
			writers = append(writers, file)
		}
		if cfg.CustomWriter != nil {
			writers = append(writers, cfg.CustomWriter)
		}
		output = io.MultiWriter(writers...)
	default:
		return nil, fmt.Errorf("invalid log output type: %s", cfg.OutputType)
	}

	return &Logger{
		level:  cfg.Level,
		logger: log.New(output, "", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var levelStr string
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	fieldStr := ""
	for k, v := range fields {
		fieldStr += fmt.Sprintf(" %s=%v", k, v)
	}

	l.logger.Printf("[%s] %s%s", levelStr, msg, fieldStr)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.log(DEBUG, msg, fields)
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log(INFO, msg, fields)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log(WARN, msg, fields)
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.log(ERROR, msg, fields)
}
