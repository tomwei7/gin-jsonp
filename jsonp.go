//gin jsonp middleware
package jsonp

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net"
)

const (
	noWritten = -1
)

type bResponseWriter struct {
	buffer *bytes.Buffer
	gin.ResponseWriter
	size     int
	isFlush  bool
	status   int
	callback string
}

func (rw *bResponseWriter) WriteHeader(code int) {
	rw.status = code
}

func (rw *bResponseWriter) WriteHeaderNow() {
	rw.ResponseWriter.WriteHeaderNow()
}

func (rw *bResponseWriter) Write(data []byte) (n int, err error) {
	n, err = rw.buffer.Write(data)
	rw.size += n
	return
}

func (rw *bResponseWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(rw.buffer, s)
	rw.size += n
	return
}

func (rw *bResponseWriter) Status() int {
	return rw.ResponseWriter.Status()
}

func (rw *bResponseWriter) Size() int {
	return rw.size
}

func (rw *bResponseWriter) Written() bool {
	return rw.size != noWritten
}

func (rw *bResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.Hijack()
}

func (rw *bResponseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.CloseNotify()
}

func (rw *bResponseWriter) Flush() {
	if rw.isFlush {
		return
	}
	rw.ResponseWriter.WriteHeader(rw.status)
	if rw.buffer.Len() > 0 {
		data := rw.buffer.Bytes()
		if data[len(data)-1] == 10 {
			data[len(data)-1] = 41
		} else {
			data = append(data, 41)
		}
		_, err := rw.ResponseWriter.Write(data)
		if err != nil {
			panic(err)
		}
		rw.buffer.Reset()
	}
	rw.isFlush = true
}

func (rw *bResponseWriter) start() {
	rw.buffer.Write([]byte(rw.callback + "("))
}

var _ gin.ResponseWriter = &bResponseWriter{}

func newbResponseWriter(rw gin.ResponseWriter, callback string) *bResponseWriter {
	bresp := &bResponseWriter{ResponseWriter: rw, buffer: &bytes.Buffer{}, status: 200, callback: callback}
	bresp.start()
	return bresp
}

func JsonP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var callback string
		if jsonp := c.DefaultQuery("jsonp", ""); jsonp != "" {
			callback = jsonp
		}
		if callbackStr := c.DefaultQuery("callback", ""); callbackStr != "" {
			callback = callbackStr
		}
		if callback == "" {
			c.Next()
		} else {
			brw := newbResponseWriter(c.Writer, callback)
			brw.Header().Set("Content-Type", "application/javascript")
			c.Writer = brw
			c.Next()
			brw.Flush()
		}
	}
}
