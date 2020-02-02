package useful

import (
	"fmt"
	"net/http"
)

func PrintReqForm(PageName string, req *http.Request) {
		fmt.Println("======>", PageName, "\n", req.Form)
}
