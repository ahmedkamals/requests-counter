package logger

import (
	"io"
	baseLogger "log"
	"os"
)

type (
	Logger struct {
		baseLogger *baseLogger.Logger
		formatter  FormatterAware
		output     io.Writer
		colorable  *Colorable
	}

	LogMessage struct {
		Level   Level
		Prefix  string
		Message string
	}

	Context interface{}

	Level uint

	LoggerFunc func(Level, string)
)

const (
	EMERGENCY Level = 9 - iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
	TRACE
	SUCCESS
)

func NewLogger() *Logger {
	output := os.Stdout
	return &Logger{
		baseLogger: baseLogger.New(output, "", baseLogger.Ldate|baseLogger.Ltime|baseLogger.Lshortfile),
		formatter:  NewFormatter(),
		output:     output,
		colorable:  NewColorable(output),
	}
}

func NewLogMessage(level Level, prefix, message string) *LogMessage {
	return &LogMessage{
		Level:   level,
		Prefix:  prefix,
		Message: message,
	}
}

func (self *Logger) Emergency(context ...Context) {
	self.LogWithLevel(EMERGENCY, context...)
}

func (self *Logger) Alert(context ...Context) {
	self.LogWithLevel(ALERT, context...)
}

func (self *Logger) Critical(context ...Context) {
	self.LogWithLevel(CRITICAL, context...)
}

func (self *Logger) Error(context ...Context) {
	self.LogWithLevel(ERROR, context...)
}

func (self *Logger) Warn(context ...Context) {
	self.LogWithLevel(WARNING, context...)
}

func (self *Logger) Notice(context ...Context) {
	self.LogWithLevel(NOTICE, context...)
}

func (self *Logger) Info(context ...Context) {
	self.LogWithLevel(INFO, context...)
}

func (self *Logger) Debug(context ...Context) {
	self.LogWithLevel(DEBUG, context...)
}

func (self *Logger) Trace(context ...Context) {
	self.LogWithLevel(TRACE, context...)
}

func (self *Logger) Success(context ...Context) {
	self.LogWithLevel(SUCCESS, context...)
}

func (self *Logger) Log(context ...Context) {
	self.verbose(self.formatter.format(context))
}

func (self *Logger) SetPrefix(prefix string) {
	self.baseLogger.SetPrefix(prefix)
}

func (self *Logger) LogWithLevel(level Level, context ...Context) {
	self.SetPrefix(self.getLevelAsString(level) + " ")
	self.verbose(self.FormatWithLevel(level, context...))
	self.SetPrefix("")
}

func (self *Logger) FormatWithLevel(level Level, context ...Context) string {
	var colorValue []ColorValue

	switch level {
	case EMERGENCY, ALERT, CRITICAL, ERROR:
		colorValue = []ColorValue{RED, BOLD}
		break
	case WARNING, NOTICE:
		colorValue = []ColorValue{YELLOW, ITALIC}
		break
	case INFO:
		colorValue = []ColorValue{BLUE, FAINT}
		break
	case DEBUG:
		colorValue = []ColorValue{MAGENTA, BLINK_RAPID}
		break
	case TRACE:
		colorValue = []ColorValue{CYAN}
		break
	case SUCCESS:
		colorValue = []ColorValue{GREEN}
		break
	default:
		colorValue = []ColorValue{WHITE}
	}

	message := self.formatter.format(context)
	return self.ColorFormat(message, colorValue...)
}

func (self *Logger) ColorFormat(message string, colorValue ...ColorValue) string {
	return self.colorable.Wrap(message, colorValue...)
}

func (self *Logger) verbose(messages string) {
	self.baseLogger.Output(3, messages)
}

func (self *Logger) getLevelAsString(l Level) string {
	levels := map[Level]string{
		EMERGENCY: "EMERGENCY",
		ALERT:     "ALERT",
		CRITICAL:  "CRITICAL",
		ERROR:     "ERROR",
		WARNING:   "WARNING",
		NOTICE:    "NOTICE",
		INFO:      "INFO",
		DEBUG:     "DEBUG",
		TRACE:     "TRACE",
		SUCCESS:   "SUCCESS",
	}

	return levels[l]
}
