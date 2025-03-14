package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)                  // Capture the response body
	return rec.ResponseWriter.Write(b) // Write to the original response
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Incoming Request: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		log.Printf("Request Body: %s", string(requestBody))

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(recorder, r) // Process request

		duration := time.Since(start)
		log.Printf("Response: Status=%d, Duration=%s", recorder.statusCode, duration)
		log.Printf("Response Body: %s", recorder.body.String())
	})
}
