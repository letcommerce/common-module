// Package middlewares contains gin middlewares
// Usage: router.Use(middlewares.Connect)
package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	requestid "github.com/sumit-tembe/gin-requestid"
	"io"
	"io/ioutil"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogErrorResponse(ctx *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = blw
	ctx.Next()
	statusCode := ctx.Writer.Status()
	if statusCode >= 400 {
		// Record the response body if there was an error
		requestId := requestid.GetRequestIDFromContext(ctx)
		log.Errorf("Got Error Response while handling URI: [%v] %v - Response Body is: [%v] %v. [%v]", ctx.Request.Method, ctx.Request.RequestURI, statusCode, blw.body.String(), requestId)
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		buf, _ := ioutil.ReadAll(ctx.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		log.Debugf("Got Reuqest for URI: [%v] [%v] - ", ctx.Request.Method, ctx.Request.RequestURI, readBody(rdr1)) // Print request body

		ctx.Request.Body = rdr2
		ctx.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
