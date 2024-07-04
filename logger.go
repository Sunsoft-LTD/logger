package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/gobwas/ws"
	"github.com/ysmood/goob"
)

type Logger struct{}

var (
	isConnected bool
	chExit      = make(chan struct{})
	errChan     = make(chan error, 1)
	Queue       *goob.Observable
)

func consoleMessage() {
	c := color.New(color.FgHiGreen, color.Bold)
	c.Println("Logger client successfully connected to the Server")
}

func connect(app string) {
	reconnect := false
	defer func() {
		isConnected = false
		if r := recover(); r != nil {
			err, _ := r.(error)
			fmt.Println("websocket panic error", err)
			if reconnect {
				time.Sleep(5 * time.Second)
				connect(app)
			}
		}
	}()
	header := ws.HandshakeHeaderHTTP{
		"Logger-App-Name": []string{app},
	}
	serverUrl, err := url.Parse("ws://localhost:7000/logger")
	if err != nil {
		os.Exit(0)
	}
START:
	dialer := ws.Dialer{
		Header:          header,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
	con, _, _, err := dialer.Dial(context.Background(), serverUrl.String())
	if err != nil {
		time.Sleep(5 * time.Second)
		goto START
	}
	isConnected = true
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)
	consoleMessage()
	pingTicker := time.NewTicker(3 * time.Minute)
	defer func() {
		pingTicker.Stop()
		if isConnected {
			con.Close()
			isConnected = false
		}
	}()
	for {
		select {
		case e := <-Queue.Subscribe(context.TODO()):
			js, err := json.Marshal(e)
			if err == nil {
				con.Write(js)
			}
		case err := <-errChan:
			reconnect = true
			isConnected = false
			panic(err)
		case <-quit:
			close(chExit)
		case <-pingTicker.C:
			if isConnected {
				if _, err := con.Write(ws.CompiledPing); err != nil {
					errChan <- fmt.Errorf("ping sending error")
					return
				}
			}
		case <-chExit:
			con.Write(ws.CompiledClose)
			con.Write(ws.CompiledCloseNormalClosure)
			os.Exit(0)
			return
		}
	}
}

func Register(app string) {
	Queue = goob.New(context.Background())
	connect(app)
}

func Error(err error, msg string, tree ...int) {
	num := 0
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
		Queue.Publish(log)
	}
}

func Warning(msg string, tree ...int) {
	num := 0
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
		Queue.Publish(log)
	}
}

func Info(msg string, tree ...int) {
	num := 0
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
		Queue.Publish(log)
	}
}

func Fatal(err error, msg string, tree ...int) {
	num := 0
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
		Queue.Publish(log)
	}
}
