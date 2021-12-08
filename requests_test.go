package requests

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {
	req := NewRequest("https://api.mineor.xyz", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiNzA3NDYwMTItMWRmZS00N2E1LWEzYjAtMGYxYzM0ZmY0ZDZlIiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6Iui2hee6p-euoeeQhuWRmCIsIkF1dGhvcml0eUlkIjoiMTAwIiwiQnVmZmVyVGltZSI6ODY0MDAsImV4cCI6MTYzOTUzMDMxOSwiaXNzIjoicW1QbHVzIiwibmJmIjoxNjM4OTI0NTE5fQ.U0Bhtn6-VduEEF0MIPcCZPbFpXKHHRXh7_G1uh63qXE")
	req.SetProxy("99.0.85.1:808")
	code, err := req.Request("GET", "/project/listProjects", nil, nil)
	fmt.Println(err)
	fmt.Println(code)
}
