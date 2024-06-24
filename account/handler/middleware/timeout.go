package middleware

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"memory-app/account/models/apprerrors"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration, errTimedOut *apprerrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {

		tw := &timeoutWriter{
			ResponseWriter: c.Writer,
			header:         make(http.Header),
		}
		c.Writer = tw
		withTimeoutCtx, cancel := context.WithTimeout(c, timeout)
		defer cancel()
		c.Request = c.Request.WithContext(withTimeoutCtx)

		finished := make(chan struct{})
		panickable := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panickable <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {

		case <-panickable:
			{

				e := apprerrors.NewInternal()
				eJSON, _ := json.Marshal(gin.H{
					"error": e,
				})
				tw.ResponseWriter.WriteHeader(e.Status())
				tw.ResponseWriter.Write(eJSON)
			}
		case <-finished:
			{
				tw.mu.Lock()
				defer tw.mu.Unlock()

				dst := tw.ResponseWriter.Header()
				for key, value := range tw.Header() {
					dst[key] = value
				}
				tw.ResponseWriter.WriteHeader(tw.code)
				tw.ResponseWriter.Write(tw.wriBuf.Bytes())

			}
		case <-withTimeoutCtx.Done():
			{

				tw.mu.Lock()
				defer tw.mu.Unlock()

				tw.ResponseWriter.Header().Set("Content-Type", "application/json")
				tw.ResponseWriter.WriteHeader(errTimedOut.Status())
				eJSON, _ := json.Marshal(gin.H{
					"error": errTimedOut,
				})
				tw.ResponseWriter.Write(eJSON)
				c.Abort()
				tw.timedOut = true
			}
		}
	}

}
