package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	rerror "github.com/Prolht/requests/error"
)

type ContentTyper interface {
	ContentType() string
}

type ContentLengther interface {
	ContentLength() int64
}

type (
	Event        func()
	EventOnError func(err error)
)

func Request(r BaseReq, method string, url string, body interface{}, respObj interface{}) (code int, err error) {
	code = http.StatusInternalServerError
	req, err := r.PrepareRequest(method, url, body)
	resp, err := r.Do(req)
	if err != nil {
		return
	}

	resp, err = r.CheckResponse(method, url, body, resp)
	if err != nil {
		code = resp.StatusCode
		return
	}
	defer resp.Body.Close()

	if respObj != nil {
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if err = json.Unmarshal(respBody, respObj); err != nil {
			err = rerror.NewJsonDecodingError(respBody, respObj, err)
			return
		}
	}
	return resp.StatusCode, nil
}

func (r *Req) Request(method string, path string, body interface{}, respObj interface{}) (code int, err error) {
	r.lazyInit()
	url := r.URLJoin(path)
	return Request(r, method, url, body, respObj)
}

func (r *Req) Do(req *http.Request) (*http.Response, error) {
	r.SetHeader(req)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Req) PrepareRequest(method string, path string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body == nil {
		reqBody = nil
	} else {
		jsonStr, err := json.Marshal(body)
		if err != nil {
			return nil, rerror.NewJSONEncodingError(body, err)
		}
		reqBody = bytes.NewReader(jsonStr)
	}
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request: %s, %s, %v: %v", method, path, body, err)
	}
	return req, nil
}

func (r *Req) CheckResponse(method string, path string, body interface{}, resp *http.Response) (*http.Response, error) {
	if IsBadResponse(resp) {
		return resp, rerror.NewBadResponseError(method, path, body, http.StatusText(resp.StatusCode))
	}
	return resp, nil
}

func IsBadResponse(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices
}
