package requests

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Prolht/requests/utils"
)

type BaseReq interface {
	Do(req *http.Request) (*http.Response, error)
	Request(method string, path string, body interface{}, respObj interface{}) (int, error)
	PrepareRequest(method string, path string, body interface{}) (*http.Request, error)
	CheckResponse(method string, path string, body interface{}, resp *http.Response) (*http.Response, error)
}

type Req struct {
	http.Client
	BaseUrl   string
	l         sync.Mutex
	https     bool
	UserAgent string
	Token     string
	transport *http.Transport
	ProxyFunc func(*http.Request) (*url.URL, error)
}

func NewRequest(url string, token string) *Req {
	h := &Req{
		BaseUrl: url,
		Token:   token,
		Client: http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: utils.DefaultUserAgent,
	}
	h.Client.Jar, _ = cookiejar.New(nil)
	return h
}

func (r *Req) lazyInit() {
	if r.transport == nil {
		r.transport = &http.Transport{
			Proxy:       proxyFunc,
			DialContext: nil,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !r.https,
			},
			TLSHandshakeTimeout:   10 * time.Second,
			DisableCompression:    false,
			DisableKeepAlives:     false,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		}
		r.Client.Transport = r.transport
	}
}

func (r *Req) URLJoin(path string) string {
	slash := "/"
	if path == "" {
		return r.BaseUrl
	}
	r.BaseUrl = strings.TrimSuffix(r.BaseUrl, slash)
	path = strings.TrimPrefix(path, slash)

	return fmt.Sprintf("%s/%s", r.BaseUrl, path)
}

// SetUserAgent 设置user agent 标识
func (r *Req) SetUserAgent(ua string) {
	r.l.Lock()
	defer r.l.Unlock()
	r.UserAgent = ua
}

// SetProxy 设置代理
func (r *Req) SetProxy(proxyAddr string) {
	r.l.Lock()
	defer r.l.Unlock()
	r.lazyInit()
	u, err := checkProxyAddr(proxyAddr)
	if err != nil {
		r.transport.Proxy = http.ProxyFromEnvironment
		return
	}
	r.transport.Proxy = http.ProxyURL(u)
}

// SetCookiejar 设置cookie
func (r *Req) SetCookiejar(jar http.CookieJar) {
	r.l.Lock()
	defer r.l.Unlock()
	r.Client.Jar = jar
}

// ResetCookiejar 清空cookie
func (r *Req) ResetCookiejar() {
	r.l.Lock()
	defer r.l.Unlock()
	r.Jar, _ = cookiejar.New(nil)
}

// SetHttpSecure 是否启用 https 安全检查, 默认不检查
func (r *Req) SetHttpSecure(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.https = b
	r.lazyInit()
	if b {
		r.transport.TLSClientConfig = nil
	} else {
		r.transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: !b,
		}
	}
}

// SetKeepAlive 设置 Keep-Alive
func (r *Req) SetKeepAlive(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.lazyInit()
	r.transport.DisableKeepAlives = !b
}

// SetGzip 是否启用Gzip
func (r *Req) SetGzip(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.lazyInit()
	r.transport.DisableCompression = !b
}

// SetResponseHeaderTimeout 设置目标服务器响应超时时间
func (r *Req) SetResponseHeaderTimeout(t time.Duration) {
	r.l.Lock()
	defer r.l.Unlock()
	r.lazyInit()
	r.transport.ResponseHeaderTimeout = t
}

// SetTLSHandshakeTimeout 设置tls握手超时时间
func (r *Req) SetTLSHandshakeTimeout(t time.Duration) {
	r.l.Lock()
	defer r.l.Unlock()
	r.lazyInit()
	r.transport.TLSHandshakeTimeout = t
}

// SetTimeout 设置 http 请求超时时间, 默认30s
func (r *Req) SetTimeout(t time.Duration) {
	r.l.Lock()
	defer r.l.Unlock()
	r.Client.Timeout = t
}

// SetTimeout 设置 http 请求超时时间, 默认30s
func (r *Req) SetHeader(req *http.Request) {
	r.l.Lock()
	defer r.l.Unlock()
	req.Header.Set("Content-Type", utils.MimeJSON)
	req.Header.Set("Accept", utils.MimeJSON)
	req.Header.Set("x-token", r.Token)
}

func (r *Req) setProxyFunc(f func(*http.Request) (*url.URL, error)) {
	r.l.Lock()
	defer r.l.Unlock()
	r.ProxyFunc = f
}

func (r *Req) SetToken(token string) {
	r.l.Lock()
	defer r.l.Unlock()
	r.Token = token
}
