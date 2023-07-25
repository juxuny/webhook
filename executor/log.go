package executor

type Logger interface {
	Printf(format string, values ...interface{})
	Println(message ...interface{})
}
