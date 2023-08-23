package gogrpcgeneric

import "log"

var Debug bool

type Logger interface {
	Info(msg ...any)
	Infof(format string, msg ...any)
	Debug(msg ...any)
	Debugf(format string, msg ...any)
}

type DefaultLogger struct{}

func (l *DefaultLogger) Info(msg ...any) {
	log.Default().Println(msg...)
}

func (l *DefaultLogger) Infof(format string, msg ...any) {
	log.Default().Printf(format, msg...)
}

func (l *DefaultLogger) Debug(msg ...any) {
	if Debug {

		log.Default().Print("\033[31m[DEBUG]")
		log.Default().Print(msg...)
		log.Default().Print("\033[0m\n")
	}
}

func (l *DefaultLogger) Debugf(format string, msg ...any) {
	if Debug {
		log.Default().Print("\033[31m[DEBUG]")
		log.Default().Printf(format, msg...)
		log.Default().Print("\033[0m\n")
	}

}
