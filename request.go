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
	"strings"
	"sync"
	"time"
)

type Request struct {
	http.Client
	req       *http.Request
	l         sync.Mutex
	https     bool
	UserAgent string
	transport *http.Transport
	ProxyFunc func(*http.Request) (*url.URL, error)
}

func NewRequest(method string, url string, body interface{}) (*Request, error) {
	normalizedMethod, normalizedUrl, err := validateParams(method, url)
	if err != nil {
		return nil, err
	}
	r := &Request{
		Client: http.Client{
			Timeout: 20 * time.Second,
		},
		UserAgent: DefaultUserAgent,
	}
	r.Client.Jar, _ = cookiejar.New(nil)
	r.lazyInit()
	r.prepareRequest(normalizedMethod, normalizedUrl, body)
	return r, nil
}

func (r *Request) Do() ([]byte, *http.Response, error) {
	respBy, resp, err := r.do()
	if err != nil {
		return respBy, resp, err
	}
	return respBy, resp, nil
}

func (r *Request) lazyInit() {
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
func (r *Request) SetUserAgent(ua string) {
	r.l.Lock()
	defer r.l.Unlock()
	r.UserAgent = ua
}

// SetProxy 设置代理
func (r *Request) SetProxy(proxyAddr string) {
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
func (r *Request) SetCookiejar(jar http.CookieJar) {
	r.l.Lock()
	defer r.l.Unlock()
	r.Client.Jar = jar
}

// ResetCookiejar 清空cookie
func (r *Request) ResetCookiejar() {
	r.l.Lock()
	defer r.l.Unlock()

	r.Jar, _ = cookiejar.New(nil)
}

// SetHttpSecure 是否启用 https 安全检查, 默认不检查
func (r *Request) SetHttpSecure(b bool) {
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
func (r *Request) SetKeepAlive(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.DisableKeepAlives = !b
}

// SetGzip 是否启用Gzip
func (r *Request) SetGzip(b bool) {
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.DisableCompression = !b
}

// SetResponseHeaderTimeout 设置目标服务器响应超时时间，单位为s
func (r *Request) SetResponseHeaderTimeout(t int64) {
	if t == 0 {
		panic("ResponseHeaderTimeout should not be zero")
	}
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.ResponseHeaderTimeout = time.Duration(t) * time.Second
}

// SetTLSHandshakeTimeout 设置tls握手超时时间
func (r *Request) SetTLSHandshakeTimeout(t int64) {
	if t == 0 {
		panic("HandshakeTimeout should not be zero")
	}
	r.l.Lock()
	defer r.l.Unlock()
	r.transport.TLSHandshakeTimeout = time.Duration(t) * time.Second
}

// SetTimeout 设置 http 请求超时时间，单位为s
func (r *Request) SetTimeout(t int64) {
	if t == 0 {
		panic("request timeout should not be zero")
	}
	r.l.Lock()
	defer r.l.Unlock()
	r.Client.Timeout = time.Duration(t) * time.Second
}

// SetHeader 设置 http 请求header
func (r *Request) SetHeader(key, value string) {
	if key == "" {
		panic("request header's key should not be empty")
	}
	r.l.Lock()
	defer r.l.Unlock()
	r.req.Header.Set(key, value)
}

func (r *Request) setDefaultHeader() {
	r.SetHeader("Content-Type", MimeJSON)
	r.SetHeader("Accept", MimeJSON)
}

func (r *Request) prepareRequest(method string, url string, body interface{}) error {
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

func (r *Request) do() ([]byte, *http.Response, error) {
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

func validateParams(method string, url string) (normalizedMethod, normalizedUl string, err error) {
	if method == "" {
		err = fmt.Errorf("empty method")
	}
	if url == "" {
		err = fmt.Errorf("empty url")
	}

	normalizedMethod = strings.ToUpper(method)
	switch normalizedMethod {
	case GET, POST, PUT, DELETE, OPTIONS, PATCH:
	default:
		err = fmt.Errorf("UnSupported method %s", url)
	}

	splitUrl := strings.Split(url, ":")
	if len(splitUrl) < 2 {
		normalizedUl = fmt.Sprintf("%s://%s", Http, url)
	} else {
		protocol := splitUrl[0]
		switch protocol {
		case Http, Https:
			normalizedUl = url
		default:
			err = fmt.Errorf("unsupport protocol <%s>, expect [%s]", protocol, strings.Join([]string{Http, Https}, ","))
		}
	}

	return normalizedMethod, normalizedUl, err
}
