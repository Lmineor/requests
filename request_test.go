package requests

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {

	req := NewRequest()
	//req.SetProxy("99.0.85.1:808")

	body, _, err := req.GET("http://www.baidu.com")

	fmt.Println(string(body))
	fmt.Println(err)
}
