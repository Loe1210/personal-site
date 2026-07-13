package xtrace

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	HeaderRequestID = "X-Request-Id"
	HeaderTraceID   = "X-Trace-Id"
)

func EnsureTraceID(c *app.RequestContext) string {
	if traceID := string(c.Request.Header.Peek(HeaderTraceID)); traceID != "" {
		c.Response.Header.Set(HeaderTraceID, traceID)
		c.Response.Header.Set(HeaderRequestID, traceID)
		return traceID
	}
	if requestID := string(c.Request.Header.Peek(HeaderRequestID)); requestID != "" {
		c.Response.Header.Set(HeaderRequestID, requestID)
		c.Response.Header.Set(HeaderTraceID, requestID)
		return requestID
	}

	traceID := newTraceID()
	c.Response.Header.Set(HeaderRequestID, traceID)
	c.Response.Header.Set(HeaderTraceID, traceID)
	return traceID
}

func newTraceID() string {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "trace-unavailable"
	}
	return hex.EncodeToString(randomBytes)
}
