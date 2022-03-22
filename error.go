package requests

import "fmt"

type JSONDecodingError struct {
	Data []byte
	V    interface{}
	Err  error
}

func (e JSONDecodingError) Error() string {
	return fmt.Sprintf("failed to perform json deserialization on %s to %v: %v", e.Data, e.V, e.Err)
}

func NewJsonDecodingError(data []byte, v interface{}, err error) JSONDecodingError {
	return JSONDecodingError{Data: data, V: v, Err: err}
}

type JSONEncodingError struct {
	V   interface{}
	Err error
}

func (e JSONEncodingError) Error() string {
	return fmt.Sprintf("fialed to perform json serialization on %v: %v", e.V, e.Err)
}

func NewJSONEncodingError(v interface{}, err error) JSONEncodingError {
	return JSONEncodingError{V: v, Err: err}
}

type BadResponseError struct {
	Method string
	Path   string
	Body   interface{}
	ErrMsg string
}

func (be BadResponseError) Error() string {
	return fmt.Sprintf("bad response form request method: %s, path: %s, body: %v: %s", be.Method, be.Path, be.Body, be.ErrMsg)
}
func NewBadResponseError(method string, path string, body interface{}, errMsg string) BadResponseError {
	return BadResponseError{
		Method: method,
		Path:   path,
		Body:   body,
		ErrMsg: errMsg,
	}
}

type ErrProxyAddrEmpty struct{}

func (ede ErrProxyAddrEmpty) Error() string {
	return fmt.Sprintf("proxy addr is empty")
}

func NewErrProxyAddrEmpty() ErrProxyAddrEmpty {
	return ErrProxyAddrEmpty{}
}
