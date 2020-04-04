package security

import (
	"fmt"
	"mysql"
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

//-- Check User Cookies
func CheckCookies(req *http.Request, res http.ResponseWriter) (error, bool) {
	if len(req.Cookies()) != 0 {
		reqCookies := req.Cookies()[0].Value

		db, err := useful.OpenDatabases()
		defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 197)
		if err != nil {
			fmt.Println(err)
			return err, false
		}

		err, dbCookies := mysql.SelectUserCookies(reqCookies, db)
		if err != nil {
			fmt.Println(err)
			return err, false
		}

		if dbCookies == reqCookies {
			return nil, true
		}
	}
	return nil, false
}

