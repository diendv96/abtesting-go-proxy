package main

import (
	"net"

	"github.com/valyala/fasthttp"
)

// NewReverseProxy ...
func newReverseProxy(addr string) *ReverseProxy {
	client := &fasthttp.HostClient{
		Addr: addr,
	}

	return &ReverseProxy{
		client: client,
	}
}

// ReverseProxy reverse handler using fasthttp.HostClient
type ReverseProxy struct {
	client *fasthttp.HostClient
}

// ServeHTTP ReverseProxy to serve
// ref to: https://golang.org/src/net/http/httputil/reverseproxy.go#L169
func (p *ReverseProxy) serveHTTP(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	res := &ctx.Response

	// prepare request(replace headers and some URL host)
	if clientIP, _, err := net.SplitHostPort(ctx.RemoteAddr().String()); err == nil {
		req.Header.Add("X-Forwarded-For", clientIP)
		req.Header.Add("Content-Type", "text/html; charset=utf-8")
	}

	for _, h := range hopHeaders {
		req.Header.Del(h)
	}

	ctx.Logger().Printf("recv a requets to proxy to: %s", p.client.Addr)
	if err := p.client.Do(req, res); err != nil {
		ctx.Logger().Printf("could not proxy: %v\n", err)
		return
	}

	for _, h := range hopHeaders {
		res.Header.Del(h)
	}
}

// SetClient ...
func (p *ReverseProxy) SetClient(addr string) *ReverseProxy {
	p.client.Addr = addr
	return p
}

// Reset ...
func (p *ReverseProxy) Reset() {
	p.client.Addr = ""
}

// Close ...
func (p *ReverseProxy) Close() {
	p.client = nil
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}
