package logger

type LogType int

const (
	TC_INFO  LogType = 0
	TC_WARN  LogType = 1
	TC_ERROR LogType = 2
	TC_FATAL LogType = 3
	TC_PANIC LogType = 4
)
