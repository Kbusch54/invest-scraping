package logg

import (
	"log"
	"os"
)

type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Panic(v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	Panicf(format string, v ...any)
}

type DefaultLogger struct {
	flags int
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func NewDefaultLog() Logger {
	f := log.Ldate | log.Ltime
	i := log.New(os.Stdout, "INFO: ", f)
	w := log.New(os.Stdout, "WARN: ", f)
	e := log.New(os.Stdout, "ERROR: ", f)

	return &DefaultLogger{
		flags: f,
		info:  i,
		warn:  w,
		error: e,
	}
}

func (d *DefaultLogger) Info(v ...any) {
	d.info.Println(v...)
}

func (d *DefaultLogger) Warn(v ...any) {
	d.warn.Println(v...)
}

func (d *DefaultLogger) Error(v ...any) {
	d.error.Println(v...)
}

func (d *DefaultLogger) Infof(format string, v ...any) {
	d.info.Printf(format, v...)
}

func (d *DefaultLogger) Warnf(format string, v ...any) {
	d.info.Printf(format, v...)
}

func (d *DefaultLogger) Errorf(format string, v ...any) {
	d.info.Printf(format, v...)
}

func (d *DefaultLogger) Panicf(format string, v ...any) {
	log.Panicf(format, v...)
}

func (d *DefaultLogger) Panic(v ...any) {
	log.Panic(v...)
}
