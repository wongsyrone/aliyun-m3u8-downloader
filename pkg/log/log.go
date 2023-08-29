package log

import "github.com/TarsCloud/TarsGo/tars/util/rogger"

var log = rogger.GetLogger("aliyun-m3u8-downloader")

func init() {
	Init()
}

func Init() {
	rogger.SetCallerFlag(false)
	rogger.SetLevel(rogger.INFO)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
