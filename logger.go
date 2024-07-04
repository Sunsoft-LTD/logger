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

var (
	chExit  = make(chan struct{})
	errChan = make(chan error, 1)
	Queue   *goob.Observable
)

func init() {
	Queue = goob.New(context.Background())
}

func consoleMessage() {
	c := color.New(color.FgHiGreen, color.Bold)
	c.Println("Logger client successfully connected to the Server")
}

func writeMessage(con *websocket.Conn) {
	for e := range Queue.Subscribe(context.TODO()) {
		if err := con.WriteJSON(e); err != nil {
			errChan <- fmt.Errorf("error sending data: %v", err)
		}
	}
}

func Register(app string) {
	reconnect := false
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			fmt.Println("websocket panic error", err)
			if reconnect {
				time.Sleep(5 * time.Second)
				Register(app)
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
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		goto START
	}
	go writeMessage(con)
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
		case err := <-errChan:
			reconnect = true
			fmt.Println(err)
			panic(err)
		case <-quit:
			close(chExit)
		case <-pingTicker.C:
			if err := con.WriteJSON(map[string]bool{"ping": true}); err != nil {
				errChan <- fmt.Errorf("ping sending error")
				fmt.Println(err)
				return
			}
		case <-chExit:
			con.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			os.Exit(0)
			return
		}
	}
}
