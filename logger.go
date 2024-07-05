package logger

import (
	"runtime"

	"github.com/go-resty/resty/v2"
)

const (
	URL string = "http://localhost:7001/"
	Inf int    = iota + 1
	Warn
	Err
	Fat
)

type (
	Logger struct {
		App string
	}

	Log struct {
		Level   int    `json:"level"`
		Line    int    `json:"line"`
		File    string `json:"file"`
		Func    string `json:"func"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	User struct {
		Name string `json:"name"`
		Id   any    `json:"id"`
		Role any    `json:"role"`
	}

	Access struct {
		Ip        string `json:"ip"`
		Route     string `json:"route"`
		Method    string `json:"method"`
		UserAgent string `json:"user_agent"`
		User      *User  `json:"user"`
	}
)

func Register(app string) *Logger {
	return &Logger{App: app}
}

// Error log
func (lg *Logger) Error(err error, msg string, tree ...int) {
	num := 1
	if len(tree) > 0 {
		num = tree[0]
	}
	pc, file, line, ok := runtime.Caller(num)
	if ok {
		log := &Log{
			File:    file,
			Line:    line,
			Level:   Err,
			Message: msg,
			Error:   err.Error(),
			Func:    runtime.FuncForPC(pc).Name(),
		}
		client := resty.New()
		go client.SetHeaders(map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Logger-App-Name": lg.App,
		}).R().SetBody(log).Post(URL)
	}
}

func (lg *Logger) Warning(msg string, tree ...int) {
	num := 1
	if len(tree) > 0 {
		num = tree[0]
	}
	pc, file, line, ok := runtime.Caller(num)
	if ok {
		log := &Log{
			File:    file,
			Line:    line,
			Level:   Warn,
			Message: msg,
			Func:    runtime.FuncForPC(pc).Name(),
		}
		client := resty.New()
		go client.SetHeaders(map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Logger-App-Name": lg.App,
		}).R().SetBody(log).Post(URL)
	}
}

func (lg *Logger) Info(msg string, tree ...int) {
	num := 1
	if len(tree) > 0 {
		num = tree[0]
	}
	pc, file, line, ok := runtime.Caller(num)
	if ok {
		log := &Log{
			File:    file,
			Line:    line,
			Level:   Inf,
			Message: msg,
			Func:    runtime.FuncForPC(pc).Name(),
		}
		client := resty.New()
		go client.SetHeaders(map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Logger-App-Name": lg.App,
		}).R().SetBody(log).Post(URL)
	}
}

func (lg *Logger) Fatal(err error, msg string, tree ...int) {
	num := 1
	if len(tree) > 0 {
		num = tree[0]
	}
	pc, file, line, ok := runtime.Caller(num)
	if ok {
		log := &Log{
			File:    file,
			Line:    line,
			Level:   Fat,
			Message: msg,
			Error:   err.Error(),
			Func:    runtime.FuncForPC(pc).Name(),
		}
		client := resty.New()
		go client.SetHeaders(map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Logger-App-Name": lg.App,
		}).R().SetBody(log).Post(URL)
	}
}

func (lg *Logger) AccessLog(log *Access) {
	client := resty.New()
	go client.SetHeaders(map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"Logger-App-Name": lg.App,
	}).R().SetBody(log).Post(URL)
}
