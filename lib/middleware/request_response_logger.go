package middleware

import (
	"log"
	"net/http"
	"strings"
)

// Wrapper around http.ResponseWriter
// used to store statusCode and body to be logged
type responseWriterRecorder struct {
	responseWriter http.ResponseWriter

	statusCode int
	body       strings.Builder
}

func (wr *responseWriterRecorder) Write(bytes []byte) (int, error) {
	wr.body.Write(bytes)
	return wr.responseWriter.Write(bytes)
}

func (wr *responseWriterRecorder) WriteHeader(statusCode int) {
	wr.statusCode = statusCode
	wr.responseWriter.WriteHeader(statusCode)
}

func (wr *responseWriterRecorder) Header() http.Header {
	return wr.responseWriter.Header()
}

func RequestResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Context().Value(reqIDKey)

		log.Println("Received request", rid, r.Method, r.URL, r.Body)
		recorder := &responseWriterRecorder{
			responseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(recorder, r)

		log.Println("Responded", rid, recorder.statusCode, recorder.body.String())
	})
}
