package main

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request-id", id)
		r = r.WithContext(ctx)
		w.Header().Set("request-id", id)
		next.ServeHTTP(w, r)
		return
	})
}

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	return w.ResponseWriter.Write(body)
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := r.Context()
		requestID := getRequestID(ctx)
		logRespWriter := LogResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&logRespWriter, r)
		duration := time.Now().Sub(startTime)
		trace := appinsights.NewRequestTelemetry(r.Method, r.URL.Path, duration, strconv.Itoa(logRespWriter.statusCode))
		trace.Id = requestID
		trace.Timestamp = time.Now()
		dtLog.Track(trace)
	})
}

func getRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request-id").(string); ok {
		return requestID
	}
	return ""
}
