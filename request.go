package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

type request struct {
	http.Client
	req       *http.Request
	l         sync.Mutex
	https     bool
	UserAgent string
	transport *http.Transport
	ProxyFunc func(*http.Request) (*url.URL, error)
}

func NewRequest() *request {
	h := &request{
		Client: http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: DefaultUserAgent,
	}
	h.Client.Jar, _ = cookiejar.New(nil)
	h.lazyInit()
	return h
}

func (r *request) GET(url string) ([]byte, *http.Response, error) {
	err := r.prepareRequest("GET", url, nil)
	if err != nil {
		return []byte{}, nil, err
	}
	respBy, resp, err := r.Do()
	if err != nil {
		return respBy, resp, err
	}
	return respBy, resp, nil
}

func (r *request) POST(url string, body interface{}) ([]byte, *http.Response, error) {
	err := r.prepareRequest("POST", url, body)
	if err != nil {
		return []byte{}, nil, err
	}
	respBy, resp, err := r.Do()
	if err != nil {
		return respBy, resp, err
	}
	return respBy, resp, nil
}
func (r *request) PUT(url string, body interface{}) ([]byte, *http.Response, error) {
	err := r.prepareRequest("PUT", url, body)
	if err != nil {
		return []byte{}, nil, err
	}
	respBy, resp, err := r.Do()
	if err != nil {
		return respBy, resp, err
	}
	return respBy, resp, nil
}
func (r *request) DELETE(url string, body interface{}) ([]byte, *http.Response, error) {
	err := r.prepareRequest("DELETE", url, body)
	if err != nil {
		return []byte{}, nil, err
	}
	respBy, resp, err := r.Do()
	if err != nil {
		return respBy, resp, err
	}
	return respBy, resp, nil
}

func (r *request) lazyInit() {
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

// SetUserAgent 设置user agent 标识
func (r *request) SetUserAgent(ua string) {
	r.l.Lock()
	defer r.l.Unlock()
	r.UserAgent = ua
}

// SetProxy 设置代理
func (r *request) SetProxy(proxyAddr string) {
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
func (r *request) SetCookiejar(jar http.CookieJar) {
	r.l.Lock()
	defer r.l.Unlock()
	r.Client.Jar = jar
}

// ResetCookiejar 清空cookie
func (r *request) ResetCookiejar() {
	r.l.Lock()
	defer r.l.Unlock()

	r.Jar, _ = cookiejar.New(nil)
}

// SetHttpSecure 是否启用 https 安全检查, 默认不检查
func (r *request) SetHttpSecure(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.https = b
	if b {
		r.transport.TLSClientConfig = nil
	} else {
		r.transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: !b,
		}
	}
}

// SetKeepAlive 设置 Keep-Alive
func (r *request) SetKeepAlive(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.DisableKeepAlives = !b
}

// SetGzip 是否启用Gzip
func (r *request) SetGzip(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.DisableCompression = !b
}

// SetResponseHeaderTimeout 设置目标服务器响应超时时间
func (r *request) SetResponseHeaderTimeout(t time.Duration) {
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.ResponseHeaderTimeout = t
}

// SetTLSHandshakeTimeout 设置tls握手超时时间
func (r *request) SetTLSHandshakeTimeout(t time.Duration) {
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.TLSHandshakeTimeout = t
}

// SetTimeout 设置 http 请求超时时间, 默认30s
func (r *request) SetTimeout(t time.Duration) {
	r.l.Lock()
	defer r.l.Unlock()
	r.Client.Timeout = t
}

// SetHeader 设置 http 请求header
func (r *request) SetHeader(key, value string) {
	if key == "" {
		panic("request header's key should not be empty")
	}
	r.l.Lock()
	defer r.l.Unlock()
	r.req.Header.Set(key, value)
}

func (r *request) setDefaultHeader() {
	r.SetHeader("Content-Type", MimeJSON)
	r.SetHeader("Accept", MimeJSON)
}
func (r *request) setProxyFunc(f func(*http.Request) (*url.URL, error)) {
	r.l.Lock()
	defer r.l.Unlock()
	r.ProxyFunc = f
}
func (r *request) prepareRequest(method string, url string, body interface{}) error {
	r.lazyInit()
	var reqBody io.Reader
	if body == nil {
		reqBody = nil
	} else {
		jsonStr, err := json.Marshal(body)
		if err != nil {
			return NewJSONEncodingError(body, err)
		}
		reqBody = bytes.NewReader(jsonStr)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create new http request: %s, %s, %v: %v", method, url, body, err)
	}
	r.req = req
	return nil
}

func (r *request) Do() ([]byte, *http.Response, error) {
	resp, err := r.Client.Do(r.req)
	if err != nil {
		return []byte{}, nil, err
	}
	defer resp.Body.Close()
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return respBody, resp, err
	}
	return respBody, resp, nil
}
