package chilog

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type KitLogEntry struct {
	Logger log.Logger // field logger interface, created by RequestLogger
}

func (l *KitLogEntry) Write(status, bytes int, elapsed time.Duration) {
	logger := log.With(l.Logger,
		"resp_status", status,
		"resp_bytes_length", bytes,
		"resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0,
	)

	level.Info(logger).Log("msg", "request complete")
}

func (l *KitLogEntry) Panic(rec interface{}, stack []byte) {
	logger := log.With(l.Logger,
		"stack", string(stack),
		"panic", fmt.Sprintf("%+v", rec),
	)
	level.Error(logger).Log("msg", "panic recovered")
}
