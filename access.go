package logger

func AccessLog(log *Access) {
	Queue.Publish(log)
}
