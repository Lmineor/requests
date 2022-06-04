package requests

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {

	req, err := NewRequest(GET, "www.baidu.com", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetProxy("99.0.85.1:808")
	//req.SetTimeout(0)
	body, _, err := req.Do()

	fmt.Println(string(body))
	fmt.Println(err)
}

//
//func TestIllegalParams(t *testing.T) {
//	_, err := NewRequest("method", "url", nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
