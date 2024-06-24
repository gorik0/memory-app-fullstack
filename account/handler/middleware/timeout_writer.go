package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type timeoutWriter struct {
	gin.ResponseWriter
	mu          sync.Mutex
	timedOut    bool
	wriBuf      bytes.Buffer
	headerWrote bool
	code        int
	header      http.Header
}

func (t *timeoutWriter) Header() http.Header {
	return t.header
}

// ::: GET sure there hasn't been timeout
func (t *timeoutWriter) Write(bytes []byte) (int, error) {

	t.mu.Lock()

	defer t.mu.Unlock()

	if t.timedOut {

		return 0, nil

	}
	return t.wriBuf.Write(bytes)

}

func (t *timeoutWriter) WriteHeader(statusCode int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	checkStatusCode(statusCode)
	if t.headerWrote || t.timedOut {
		return
	}
	t.wroteCode(statusCode)
}

func (t *timeoutWriter) wroteCode(code int) {
	t.code = code
	t.headerWrote = true

}

func checkStatusCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("Erro status code ::: %d ::: (out of bounds", code))
	}

}

var _ http.ResponseWriter = &timeoutWriter{}
