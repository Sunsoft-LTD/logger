package logger

func (*logger) AccessLog(log *Access) {
	Queue.Publish(log)
}
