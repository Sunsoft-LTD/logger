package logger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
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
	header := http.Header{
		"Logger-App-Name": []string{app},
	}
	serverUrl, err := url.Parse("ws://127.0.0.1:7000/logger")
	if err != nil {
		os.Exit(0)
	}
START:
	con, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), header)
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
			if isConnected {
				if err := con.WriteJSON(e); err != nil {
					errChan <- fmt.Errorf("error sending data: %v", err)
				}
			}
		case err := <-errChan:
			reconnect = true
			isConnected = false
			panic(err)
		case <-quit:
			close(chExit)
		case <-pingTicker.C:
			if isConnected {
				if err := con.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					errChan <- fmt.Errorf("ping sending error")
					return
				}
			}
		case <-chExit:
			con.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			os.Exit(0)
			return
		}
	}
}

func Register(app string) (*Logger, error) {
	Queue = goob.New(context.Background())
	connect(app)
	if !isConnected {
		return nil, errors.New("client did not connect")
	}
	return &Logger{}, nil
}
