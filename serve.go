package traefik_headers

import (
	"net/http"
	"sync/atomic"
)


func (h *Headers) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.next.ServeHTTP(&responseWriter{
		writer:  rw,
		headers: ghs.headers[int(atomic.LoadInt32(ghs.curheader))],
	}, req)
}

type responseWriter struct {
	writer  http.ResponseWriter
	headers *headers
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	h := r.writer.Header()
	for k, vv := range r.headers.headers {
		if _, ok := h[k]; !ok {
			for _, v := range vv { 
				r.writer.Header().Add(k, v)
			}
		}
	}
	r.writer.WriteHeader(statusCode)
}
