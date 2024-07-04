package logger

import (
	"context"
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
	chExit  = make(chan struct{})
	errChan = make(chan error, 1)
	Queue   *goob.Observable
)

func consoleMessage() {
	c := color.New(color.FgHiGreen, color.Bold)
	c.Println("Logger client successfully connected to the Server")
}

func connect(app string) {
	reconnect := false
	defer func() {
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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)
	consoleMessage()
	pingTicker := time.NewTicker(3 * time.Minute)
	defer func() {
		pingTicker.Stop()
		con.Close()
	}()
	for {
		select {
		case e := <-Queue.Subscribe(context.TODO()):
			if err := con.WriteJSON(e); err != nil {
				errChan <- fmt.Errorf("error sending data: %v", err)
			}
		case err := <-errChan:
			reconnect = true
			panic(err)
		case <-quit:
			close(chExit)
		case <-pingTicker.C:
			if err := con.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				errChan <- fmt.Errorf("ping sending error")
				return
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
	return &Logger{}, nil
}
