// Package httperr is in need of a good package description
package httperr

import (
	"errors"
	"fmt"
	"net/http"
)

// Error... needs to be better documented.
type Error struct {
	code   int
	level  LogLevel
	logged bool
	err    error
	mesg   *Message
}

// NewError returns a new Error that formats as the given text and applies
// the default options Status(http.StatusInternalServerError) and LogError.
// For different option values, use New instead of NewError -- or, you may
// call WithOptions on the return value to modify its options.
func NewError(text string) *Error {
	return &Error{
		code:  http.StatusInternalServerError,
		level: LogError,
		err:   errors.New(text),
	}
}

// New is a convenience wrapper around NewError(text).WithOptions(options...).
func New(text string, options ...Option) *Error {
	return NewError(text).WithOptions(options...)
}

// Errorf is a convenience wrapper around NewError(fmt.Sprintf(text, args...)).
func Errorf(text string, args ...interface{}) *Error {
	return NewError(fmt.Sprintf(text, args...))
}

// LogErrorf is a convenience wrapper around Errorf(text, args...).Log()
func LogErrorf(text string, args ...interface{}) *Error {
	return NewError(fmt.Sprintf(text, args...)).log(2)
}

// WithError returns err as an *Error or nil if err is nil. If err is
// already and Error, its type assertion is taken, otherwise, err is
// used to create a new Error. If any options are provided they override
// those from an existing Error or are applied to a newly created Error.
// If a new Error is created and no options are given, default options
// (as from NewError) are applied.
func WithError(err error, options ...Option) *Error {
	if err == nil {
		return nil
	}

	var e *Error
	switch err.(type) {
	case *Error:
		e = err.(*Error)
	default:
		e = NewError(err.Error())
	}

	return e.WithOptions(options...)
}

func LogWithError(err error, options ...Option) *Error {
	return WithError(err, options...).log(2)
}

// Abort is a wrapper around WithError(err, options).Abort(w). See the Abort
// method for details.
func Abort(w http.ResponseWriter, err error, options ...Option) *Error {
	return WithError(err, options...).abort(w)
}

// Abort is a wrapper around both Send and Log. If e is non-nil, Send will be
// called to deliver an http error to w. It then makes a final attempt to log
// the error by calling Log.  If this error has previously been logged, no log
// message will be emitted.  If e is nil, nothing will be done. Either way,
// e will be returned.
func (e *Error) Abort(w http.ResponseWriter) *Error {
	return e.abort(w)
}

// This method exists for both the Abort method and function to call, and we
// can keep the same call stack depth when calling log.
func (e *Error) abort(w http.ResponseWriter) *Error {
	if e != nil {
		e.Send(w).log(3)
	}
	return e
}

// WithOptions applies all options to e and returns the results. Note that the
// return value is a mere convenience; the reciever will also be modified.
func (e *Error) WithOptions(options ...Option) *Error {
	for _, o := range options {
		o.update(e)
	}

	return e
}

// Error implements the error interface
func (e *Error) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return e.err.Error()
}

func Annotate(err error, text string) *Error {
	return WithError(err).Annotate(text)
}

func Annotatef(err error, text string, args ...interface{}) *Error {
	return Annotate(err, fmt.Sprintf(text, args...))
}

func (e *Error) Annotate(text string) *Error {
	if e == nil {
		return NewError(text)
	}

	e.err = fmt.Errorf("%s: %v", text, e.err)
	e.logged = false

	return e
}

// Send calls http.Error using w, a derived error message and the http status
// code attached to e -- if and only if both e and w are non-nil and e's status
// code is a valid http status (as determined by calling http.StatusText). If
// any of these tests fail, nothing will be done.  Otherwise, the error message
// will be derived based on one of the following Message option values:
//
// * HTTPMessage: The results from http.StatusText will be used. (default)
// * ErrorMessage: The result of calling e.Error will be used.
// * AltMessage: The alternative message text will be used.
//
func (e *Error) Send(w http.ResponseWriter) *Error {
	if e == nil || w == nil {
		return e
	}

	st := http.StatusText(e.code)
	if st == "" {
		return e
	}

	var text string

	switch {
	case e.mesg == nil:
		text = st

	case e.mesg.text == "":
		text = e.Error()

	default:
		text = e.mesg.text
	}

	http.Error(w, text, e.code)

	return e
}

// Log emits a log message containing the text for the attached error at the
// level specified by the LogLevel option then returns e. Each error may only
// be logged once; if e has already been logged, no message will be written.
//
// If the LogLevel option is LogNone, no log message will be written -- but e
// will still be marked as logged nonetheless. If the LogLevel option is later
// updated, the log message will continue to be suppressed on repeated calls
// to Log.
//
// A call to one of the Annotate* functions will clear this flag and allow
// another log message to be emitted.
func (e *Error) Log() *Error {
	return e.log(2)
}

func (e *Error) log(depth int) *Error {
	if e != nil && e.err != nil && !e.logged {
		e.level.logMesg(depth+1, e.err)
		e.logged = true
	}

	return e
}
