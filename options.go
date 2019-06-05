package httperr

import (
	"fmt"

	"toolman.org/base/log/v2"
)

type Option interface {
	update(*Error)
}

type Status int

const StatusNone Status = 0

func (s Status) update(e *Error) {
	e.code = int(s)
}

type LogLevel int

const (
	LogNone LogLevel = iota
	LogInfo
	LogWarning
	LogError
	LogFatal
)

func LogVerbose(v log.Level) LogLevel {
	return LogLevel(int(v) * -1)
}

func (l LogLevel) update(e *Error) {
	e.level = l
}

func (l LogLevel) logMesg(depth int, args ...interface{}) {
	if v := log.Level(-1 * int(l)); v > 0 && log.V(v) {
		l = LogInfo
	}

	var logfunc func(int, ...interface{})

	switch l {
	case LogInfo:
		logfunc = log.InfoDepth

	case LogWarning:
		logfunc = log.WarningDepth

	case LogError:
		logfunc = log.ErrorDepth

	case LogFatal:
		logfunc = log.FatalDepth

	default:
		return
	}

	logfunc(depth, args...)
}

type Message struct {
	text string
}

var (
	HTTPMessage  *Message
	ErrorMessage = &Message{""}
)

func AltMessage(text string) *Message {
	return &Message{text}
}

func AltMessagef(text string, args ...interface{}) *Message {
	return AltMessage(fmt.Sprintf(text, args...))
}

func (m *Message) update(e *Error) {
	e.mesg = m
}
