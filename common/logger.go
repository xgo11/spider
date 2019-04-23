package common

import (
	"github.com/sirupsen/logrus"
)

type logAdapter struct {
	params logrus.Fields
	logger logrus.FieldLogger
}

func (lg *logAdapter) WithField(key string, value interface{}) *logrus.Entry {
	return lg.logger.WithFields(lg.params).WithField(key, value)
}

func (lg *logAdapter) WithFields(fields logrus.Fields) *logrus.Entry {
	return lg.logger.WithFields(lg.params).WithFields(fields)
}

func (lg *logAdapter) WithError(err error) *logrus.Entry {
	return lg.logger.WithFields(lg.params).WithError(err)
}

func (lg *logAdapter) Debugf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Debugf(format, args...)
}

func (lg *logAdapter) Infof(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Infof(format, args...)
}

func (lg *logAdapter) Printf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Printf(format, args...)
}

func (lg *logAdapter) Warnf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Warnf(format, args...)
}

func (lg *logAdapter) Warningf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Warningf(format, args...)
}

func (lg *logAdapter) Errorf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Errorf(format, args...)
}

func (lg *logAdapter) Fatalf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Fatalf(format, args...)
}

func (lg *logAdapter) Panicf(format string, args ...interface{}) {
	lg.logger.WithFields(lg.params).Panicf(format, args...)
}

func (lg *logAdapter) Debug(args ...interface{}) {
	lg.logger.WithFields(lg.params).Debug(args...)
}

func (lg *logAdapter) Info(args ...interface{}) {
	lg.logger.WithFields(lg.params).Info(args...)
}

func (lg *logAdapter) Print(args ...interface{}) {
	lg.logger.WithFields(lg.params).Print(args...)
}

func (lg *logAdapter) Warn(args ...interface{}) {
	lg.logger.WithFields(lg.params).Warn(args...)
}

func (lg *logAdapter) Warning(args ...interface{}) {
	lg.logger.WithFields(lg.params).Warning(args...)
}

func (lg *logAdapter) Error(args ...interface{}) {
	lg.logger.WithFields(lg.params).Error(args...)
}

func (lg *logAdapter) Fatal(args ...interface{}) {
	lg.logger.WithFields(lg.params).Fatal(args...)
}

func (lg *logAdapter) Panic(args ...interface{}) {
	lg.logger.WithFields(lg.params).Panic(args...)
}

func (lg *logAdapter) Debugln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Debugln(args...)
}

func (lg *logAdapter) Infoln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Infoln(args...)
}

func (lg *logAdapter) Println(args ...interface{}) {
	lg.logger.WithFields(lg.params).Println(args...)
}

func (lg *logAdapter) Warnln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Warnln(args...)
}

func (lg *logAdapter) Warningln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Warningln(args...)
}

func (lg *logAdapter) Errorln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Errorln(args...)
}

func (lg *logAdapter) Fatalln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Fatalln(args...)
}

func (lg *logAdapter) Panicln(args ...interface{}) {
	lg.logger.WithFields(lg.params).Panicln(args...)
}

var (
	loggers = make(map[string]logrus.FieldLogger)
)

func GetLoggerWithParent(name string, parent logrus.FieldLogger, params ...interface{}) logrus.FieldLogger {
	if _, ok := loggers[name]; ok {
		return loggers[name]
	}
	var lg = &logAdapter{}
	if parent != nil {
		lg.logger = parent
	} else {
		lg.logger = logrus.New()
	}

	lg.params = logrus.Fields{"name": name}
	for i, s := 0, len(params); i+1 < s; i += 2 {
		k, v := params[i], params[i+1]
		if k == nil || v == nil {
			continue
		}
		if ks, ok := k.(string); ok {
			lg.params[ks] = v
		}
	}

	loggers[name] = lg
	return lg
}

func GetLogger(name string, params ...interface{}) logrus.FieldLogger {
	if lg, ok := loggers[name]; ok {
		return lg
	}
	return GetLoggerWithParent(name, nil, params...)
}
