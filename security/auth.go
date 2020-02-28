package security

import (
	"fmt"
	"net/http"
	"useful"
)

func CheckUserAuth(req *http.Request) {

	db,err:=useful.OpenDatabases()
	defer useful.HandleCloseableErrorClient(db,"auth.go", 11)
	if err != nil {
	fmt.Println(err)
	}






}
