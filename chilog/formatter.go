package chilog

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
)

type KitLogger struct {
	Logger log.Logger
}

func (l *KitLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &KitLogEntry{Logger: l.Logger}

	logFields := make([]interface{}, 0, 16)
	//logFields = append(logFields, "ts", time.Now().UTC().Format(time.RFC3339Nano))
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, "req_id", reqID)
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields = append(logFields, "http_scheme", scheme)
	logFields = append(logFields, "http_proto", r.Proto)
	logFields = append(logFields, "http_method", r.Method)
	logFields = append(logFields, "remote_addr", r.RemoteAddr)
	logFields = append(logFields, "user_agent", r.UserAgent())
	logFields = append(logFields, "uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))

	entry.Logger = log.With(entry.Logger, logFields...)
	return entry
}
