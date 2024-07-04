package logger

import "fmt"

func AccessLog(log *Access) {
	if err := conn.WriteJSON(log); err != nil {
		fmt.Println(err)
		errChan <- err
	}
}
