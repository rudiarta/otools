package outils

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

type HttpRequestParams struct {
	Context context.Context
	Method  string // e.g. POST, GET, etc
	URL     string // e.g. https://url.com
	Body    io.Reader
	Header  map[string]string
	Timeout time.Duration
}

type HttpRequestResponse struct {
	Body       []byte
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Header     http.Header
	Statistic  HttpRequestStatistic
}

type HttpRequestStatistic struct {
	startTimeTLSHandshake time.Time
	endTimeTLSHandshake   time.Time

	DurationTLSHandshake     time.Duration
	DurationTotalHTTPRequest time.Duration
}

// HTTPRequestJSON ...
func HTTPRequestJSON(params *HttpRequestParams) (*HttpRequestResponse, error) {
	req, err := http.NewRequestWithContext(params.Context, params.Method, params.URL, params.Body)
	if err != nil {
		return &HttpRequestResponse{
			Body:       nil,
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// iterate optional data of headers
	for key, value := range params.Header {
		req.Header.Set(key, value)
	}

	// calculate http statistic
	stat := HttpRequestStatistic{}
	clientTrace := &httptrace.ClientTrace{
		TLSHandshakeStart: func() {
			stat.startTimeTLSHandshake = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			stat.endTimeTLSHandshake = time.Now()
			stat.DurationTLSHandshake = stat.endTimeTLSHandshake.Sub(stat.startTimeTLSHandshake)
		},
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)

	// set timeout
	client := &http.Client{Timeout: params.Timeout}

	// send http request
	startTimeRequest := time.Now()
	r, err := client.Do(req)
	stat.DurationTotalHTTPRequest = time.Since(startTimeRequest)
	if err != nil {
		return &HttpRequestResponse{
			Body:       nil,
			StatusCode: http.StatusBadRequest,
			Statistic:  stat,
		}, err
	}

	defer func() {
		r.Body.Close()
	}()

	resp := StreamToByte(r.Body)

	return &HttpRequestResponse{
		Body:       resp,
		Status:     r.Status,
		StatusCode: r.StatusCode,
		Header:     r.Header,
		Statistic:  stat,
	}, nil
}

// StreamToByte ...
func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
