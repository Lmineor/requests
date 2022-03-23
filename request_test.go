package requests

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {

	req := NewRequest(GET, "https://www.baidu.com", nil)
	//req.SetProxy("99.0.85.1:808")

	body, _, err := req.Do()

	fmt.Println(string(body))
	fmt.Println(err)
}
