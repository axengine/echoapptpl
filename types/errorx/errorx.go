package errorx

import (
	"fmt"
	"golang.org/x/text/language"
	"runtime"
)

// Message is type of map,key is code,value is message string
type Message map[int]string

// langToMessage key is lang,value is Message
var langToMessage = map[language.Tag]Message{}

// defile error code
const (
	// 0-999 reserved
	CodeSuccess = 0
	CodeFailed  = 1
	CodeSystem  = 2

	CodeInvalidToken = 401 // token无效
	CodeForbidden    = 403
	CodeNotFound     = 404 // 未找到相关资源

)

// define Err
var (
	ErrSystem       = NewError(CodeSystem, GetMessage(CodeSystem))
	ErrFailed       = NewError(CodeFailed, GetMessage(CodeFailed))
	ErrNotFound     = NewError(CodeNotFound, GetMessage(CodeNotFound))
	ErrInvalidToken = NewError(CodeInvalidToken, GetMessage(CodeInvalidToken))
	ErrForbidden    = NewError(CodeForbidden, GetMessage(CodeForbidden))
)

type Error struct {
	Code    int      `json:"code" xml:"code"`
	Message string   `json:"message" xml:"message"`
	Stack   []string `json:"-" xml:"-"`
}

func NewError(code int, msg string) Error {
	return Error{Code: code, Message: msg, Stack: make([]string, 0, 1)}
}

func (e Error) Error() string {
	return fmt.Sprintf("code:%d message:%s statck:%v", e.Code, e.Message, e.Stack)
}

func (e Error) String() string {
	return e.Message
}

func (e Error) MultiErr(err error) Error {
	if err != nil {
		if e1, ok := err.(Error); ok {
			e.Message += ", " + e1.Message
			return e.addStack(e1.Stack...)
		} else {
			e.Message += ", " + err.Error()
		}
	}
	return e
}

func (e Error) MultiMsg(v ...interface{}) Error {
	msg := fmt.Sprint(v...)
	if msg != "" {
		e.Message += ", " + msg
	}
	return e
}

func (e Error) CodeMsg() (int, string) {
	return e.Code, e.Message
}

func (e Error) addStack(v ...string) Error {
	e.Stack = append(e.Stack, v...)
	return e
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	frames := runtime.CallersFrames(callers())
	stack := ""
	frame, more := frames.Next()
	if more {
		stack = fmt.Sprintf("\n%s:%d %s", frame.File, frame.Line, frame.Function)
	}

	if e, ok := err.(Error); ok {
		return e.addStack(stack)
	}
	return NewError(CodeFailed, err.Error()).addStack(stack)
}

func callers() []uintptr {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

// ToSystemError if not Error, to system error
func ToSystemError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return e
	}
	// to internal
	return NewError(CodeSystem, err.Error())
}

func GetMessage(code int, langs ...language.Tag) string {
	for _, lang := range langs {
		if v, ok := langToMessage[lang]; ok {
			return v[code]
		}
	}
	return langToMessage[language.English][code]
}
