package logger

import (
	"github.com/andyollylarkin/smudge-custom-transport"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	l *logrus.Entry
}

func NewLogrusLogger(l *logrus.Logger, logLvl logrus.Level) *LogrusLogger {
	l.Level = logLvl

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"

	e := l.WithFields(logrus.Fields{
		"context": "gossip",
	})
	e.Logger.Formatter = customFormatter

	return &LogrusLogger{
		l: e,
	}
}

func (l *LogrusLogger) Log(level smudge.LogLevel, a ...interface{}) (int, error) {
	l.l.Log(toLogrusLogLevel(level), a...)

	return 0, nil
}

func (l *LogrusLogger) Logf(level smudge.LogLevel, format string, a ...interface{}) (int, error) {
	l.l.Logf(toLogrusLogLevel(level), format, a...)

	return 0, nil
}

func toLogrusLogLevel(level smudge.LogLevel) logrus.Level {
	switch level {
	case smudge.LogAll:
		return logrus.FatalLevel
	case smudge.LogTrace:
		return logrus.TraceLevel
	case smudge.LogDebug:
		return logrus.DebugLevel
	case smudge.LogInfo:
		return logrus.InfoLevel
	case smudge.LogWarn:
		return logrus.WarnLevel
	case smudge.LogError:
		return logrus.ErrorLevel
	case smudge.LogFatal:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
