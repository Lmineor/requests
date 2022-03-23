# requests

---

Go语言中的requests库

用法：

`go get github.com/Lmineor/requests`

```go

package main

import (
	"fmt"

	"github.com/Lmineor/requests"
)

func main() {
	req := requests.NewRequest(requests.GET, "https://www.baidu.com", nil)
	body, _, err := req.Do()

	fmt.Println(string(body))
	fmt.Println(err)
}


```