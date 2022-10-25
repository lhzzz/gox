package recovery

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	reallyCrash = true
)

var panicHandlers = []func(interface{}){logPanic}

func HandleCrash(additionalHandlers ...func(interface{})) {
	if r := recover(); r != nil {
		for _, fn := range panicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		if reallyCrash {
			// Actually proceed to panic.
			panic(r)
		}
	}
}

func RecoverFromPanic(err *error) {
	if r := recover(); r != nil {
		// Same as stdlib http server code. Manually allocate stack trace buffer size
		// to prevent excessively large logs
		const size = 64 << 10
		stacktrace := make([]byte, size)
		stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]

		*err = fmt.Errorf(
			"recovered from panic %q. (err=%v) Call stack:\n%s",
			r,
			*err,
			stacktrace)
	}
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func logPanic(r interface{}) {
	// Same as stdlib http server code. Manually allocate stack trace buffer size
	// to prevent excessively large logs
	const size = 64 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	if _, ok := r.(string); ok {
		logrus.Errorf("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		logrus.Errorf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
