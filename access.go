package logger

func (*Logger) AccessLog(log *Access) {
	Queue.Publish(log)
}
