package log

type Logger struct {
	logging func(msg string)
}

func NewLogger(logging func(msg string)) *Logger {
	return &Logger{
		logging: logging,
	}
}

func (l Logger) Log(msg string) {
	l.logging(msg)
}

func (l Logger) LogSubmsg(msg string, submsgs ...string) {
	l.logging(msg)
	for _, submsg := range submsgs {
		l.logging(" > " + submsg)
	}
}
