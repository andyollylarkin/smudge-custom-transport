package logger

import (
	"github.com/andyollylarkin/smudge-custom-transport"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	l *logrus.Logger
}

func NewLogrusLogger(l *logrus.Logger, logLvl logrus.Level) *LogrusLogger {
	l.Level = logLvl

	return &LogrusLogger{
		l: l,
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
