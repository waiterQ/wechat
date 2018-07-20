package errors

const (
	DEBUG = iota + 1
	INFO
	WARN
	ERROR
	FATAL
)

// level error
type Lerror interface {
	error
	Level() int
}

func New(text string, level ...int) Lerror {
	if len(level) > 0 {
		return &errorLevel{text, level[0]}
	}
	return &errorLevel{text, ERROR}
}

func Transform(err error, level ...int) Lerror {
	if len(level) > 0 {
		return &errorLevel{err.Error(), level[0]}
	}
	return &errorLevel{err.Error(), ERROR}
}

type errorLevel struct {
	s     string
	level int
}

func (e *errorLevel) Error() string {
	return e.s
}

func (e *errorLevel) Level() int {
	return e.level
}
