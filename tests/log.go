package tests

import (
	"io"
	"testing"
	"log"
)

// From https://github.com/sirupsen/logrus/issues/834, adapted for standard log

// LogCapturer reroutes testing.T log output
type LogCapturer interface {
	Release()
}

type logCapturer struct {
	*testing.T
	origOut io.Writer
}

func (tl logCapturer) Write(p []byte) (n int, err error) {
	tl.Logf((string)(p))
	return len(p), nil
}

func (tl logCapturer) Release() {
	log.SetOutput(tl.origOut)
}

// CaptureLog redirects log output to testing.Log
func CaptureLog(t *testing.T) LogCapturer {
	lc := logCapturer{T: t, origOut: log.Default().Writer()}
	if !testing.Verbose() {
		log.SetOutput(lc)
	}
	return &lc
}