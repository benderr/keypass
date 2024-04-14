package logger

import "go.uber.org/zap"

type Logger interface {
	Infoln(args ...interface{})
	Errorln(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

func New() (Logger, func() error) {
	l, lerr := zap.NewDevelopment()
	if lerr != nil {
		panic(lerr)
	}

	sugar := *l.Sugar()
	return &sugar, l.Sync
}
