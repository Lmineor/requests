package requests

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {

	req, err := NewRequest(GET, "https://www.baidu.com", nil)
	//req.SetProxy("99.0.85.1:808")

	body, _, err := req.Do()

	fmt.Println(string(body))
	fmt.Println(err)
}

func TestIllegalParams(t *testing.T) {
	_, err := NewRequest("method", "url", nil)
	if err != nil {
		fmt.Println(err)
	}
}
