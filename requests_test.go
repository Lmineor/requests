package requests

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {

	req := NewRequest("https://api.mineor.xyz", "")
	//req.SetProxy("99.0.85.1:808")

	code, err := req.Request("GET", "/project/listProjects", nil, nil)
	fmt.Println(err)
	fmt.Println(code)
}
