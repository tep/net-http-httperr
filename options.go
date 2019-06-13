// Copyright 2019 Timothy E. Peoples
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

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
