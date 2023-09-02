package tracer

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/qiniu/x/humanize"
)

// -------------------------------------------------------------------------------

type teeResponseWriter struct {
	a, b http.ResponseWriter
}

func (p *teeResponseWriter) Header() http.Header {
	return p.a.Header()
}

func (p *teeResponseWriter) Write(buf []byte) (n int, err error) {
	n, err = p.a.Write(buf)
	p.b.Write(buf)
	return
}

func (p *teeResponseWriter) WriteHeader(statusCode int) {
	p.a.WriteHeader(statusCode)
	p.b.WriteHeader(statusCode)
}

func (p *teeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return p.a.(http.Hijacker).Hijack()
}

func Tee(a, b http.ResponseWriter) http.ResponseWriter {
	return &teeResponseWriter{a, b}
}

// -------------------------------------------------------------------------------

type ResponseRecorder struct {
	Code      int
	Bytes     int64
	HeaderMap http.Header
}

func (p *ResponseRecorder) Header() http.Header {
	return p.HeaderMap
}

func (p *ResponseRecorder) Write(buf []byte) (n int, err error) {
	p.Bytes += int64(len(buf))
	return
}

func (p *ResponseRecorder) WriteHeader(statusCode int) {
	p.Code = statusCode
}

func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{Code: 200, HeaderMap: make(http.Header)}
}

// -------------------------------------------------------------------------------

func New(h http.Handler) http.Handler {
	return NewWith(h, log.Default())
}

func NewWith(h http.Handler, log *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := NewRecorder()
		tee := Tee(w, recorder)
		log.Printf("%s %v\n", r.Method, r.URL)
		start := time.Now()
		h.ServeHTTP(tee, r)
		dur := time.Since(start)
		bytes := humanize.Comma(recorder.Bytes)
		log.Printf("Returned %d of %s %v with %s bytes in %d ms\n", recorder.Code, r.Method, r.URL, bytes, dur.Milliseconds())
	})
}

// -------------------------------------------------------------------------------
