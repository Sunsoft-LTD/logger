package logger

import (
	"fmt"
	"runtime"
)

// Error log
func Error(err error, msg string, tree ...int) {
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
		if err := conn.WriteJSON(log); err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}
}

func Warning(msg string, tree ...int) {
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
		if err := conn.WriteJSON(log); err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}
}

func Info(msg string, tree ...int) {
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
		if err := conn.WriteJSON(log); err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}
}

func Fatal(err error, msg string, tree ...int) {
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
		if err := conn.WriteJSON(log); err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}
}
