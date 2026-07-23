package proxy

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xhttp"
	"github.com/Loe1210/personal-site/internal/xresilience"
)

type Options struct {
	TargetBaseURL string
	StripPrefix   string
	Timeout       time.Duration
	MaxAttempts   int
	Backoff       time.Duration
	Client        *http.Client
	Stats         *xresilience.FailureStats
	Breaker       *xresilience.CircuitBreaker
}

type retryableStatusError struct{ status int }

func (e retryableStatusError) Error() string { return http.StatusText(e.status) }

func RewritePath(path string, stripPrefix string) string {
	rewritten := strings.TrimPrefix(path, stripPrefix)
	if rewritten == "" {
		return "/"
	}
	if !strings.HasPrefix(rewritten, "/") {
		return "/" + rewritten
	}
	return rewritten
}

func NewReverseProxy(targetBaseURL string, stripPrefix string) app.HandlerFunc {
	return NewReverseProxyWithOptions(Options{TargetBaseURL: targetBaseURL, StripPrefix: stripPrefix})
}

func NewReverseProxyWithOptions(opts Options) app.HandlerFunc {
	baseURL := strings.TrimRight(opts.TargetBaseURL, "/")
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = xresilience.DefaultGatewayProxyTimeout
	}
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	maxAttempts := opts.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = xresilience.DefaultReadMaxAttempts
	}
	backoff := opts.Backoff
	if backoff <= 0 {
		backoff = xresilience.DefaultRetryBackoff
	}
	stats := opts.Stats
	if stats == nil {
		stats = &xresilience.FailureStats{}
	}
	breaker := opts.Breaker
	if breaker == nil && baseURL != "" {
		breaker = xresilience.NewCircuitBreaker(xresilience.CircuitBreakerConfig{Name: opts.StripPrefix})
	}

	return func(ctx context.Context, c *app.RequestContext) {
		if baseURL == "" {
			xhttp.Fail(c, xerrors.New(xerrors.CodeGatewayUpstreamMissing, "upstream is not configured"))
			return
		}
		path := RewritePath(string(c.Path()), opts.StripPrefix)
		if query := string(c.QueryArgs().QueryString()); query != "" {
			path += "?" + query
		}
		method := string(c.Method())
		body := append([]byte(nil), c.Request.Body()...)
		requestCtx, cancel := xresilience.WithTimeout(ctx, timeout)
		defer cancel()

		var resp *http.Response
		err := breaker.Run(func() error {
			policy := xresilience.RetryPolicy{MaxAttempts: 1}
			if xresilience.IsIdempotentReadMethod(method) {
				policy = xresilience.RetryPolicy{
					MaxAttempts: maxAttempts,
					Backoff:     backoff,
					Retryable: func(err error) bool {
						if err == nil {
							return false
						}
						if xresilience.IsTimeoutError(err) {
							return true
						}
						var statusErr retryableStatusError
						return errors.As(err, &statusErr)
					},
				}
			}
			return xresilience.DoWithRetry(requestCtx, policy, func(attemptCtx context.Context) error {
				attemptReq, reqErr := http.NewRequestWithContext(attemptCtx, method, baseURL+path, bytes.NewReader(body))
				if reqErr != nil {
					return reqErr
				}
				c.Request.Header.VisitAll(func(key, value []byte) {
					attemptReq.Header.Set(string(key), string(value))
				})
				var callErr error
				resp, callErr = client.Do(attemptReq)
				if callErr != nil {
					return callErr
				}
				if resp.StatusCode >= http.StatusInternalServerError {
					_, _ = io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()
					return retryableStatusError{status: resp.StatusCode}
				}
				return nil
			})
		})
		if err != nil {
			stats.RecordFailure()
			logBreaker("gateway_proxy", breaker, err)
			if errors.Is(err, xresilience.ErrCircuitOpen) {
				xhttp.Fail(c, xerrors.New(xerrors.CodeGatewayCircuitOpen, "upstream circuit open"))
				return
			}
			if xresilience.IsTimeoutError(err) || errors.Is(err, context.DeadlineExceeded) {
				xhttp.Fail(c, xerrors.New(xerrors.CodeGatewayUpstreamTimeout, "upstream request timeout"))
				return
			}
			xhttp.Fail(c, xerrors.New(xerrors.CodeGatewayUpstreamFailed, "upstream request failed"))
			return
		}
		stats.RecordSuccess()
		logBreaker("gateway_proxy", breaker, nil)
		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			xhttp.Fail(c, xerrors.New(xerrors.CodeGatewayUpstreamFailed, "read upstream response failed"))
			return
		}
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response.Header.Add(key, value)
			}
		}
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
	}
}

func logBreaker(component string, breaker *xresilience.CircuitBreaker, err error) {
	if breaker == nil {
		return
	}
	snapshot := breaker.Snapshot()
	log.Printf("component=%s breaker=%s breaker_state=%s breaker_rejected=%d breaker_failures=%d breaker_recoveries=%d err=%v", component, snapshot.Name, snapshot.State, snapshot.RejectedCalls, snapshot.FailureCount, snapshot.RecoveryCount, err)
}
