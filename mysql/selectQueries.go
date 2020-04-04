package mysql

import (
	"cons"
	"database/sql"
	"fmt"
	"structs"
	"useful"
)

func SelectUserAuth(UserEmail, userName string, db *sql.DB) (error, *structs.EmailPass) {

	var key, value string

	if UserEmail != "" {
		key = cons.C_EMAIL
		value = UserEmail
	} else if userName != "" {
		key = cons.C_USER_NAME
		value = userName
	}
	qS := ` Select ` + cons.C_USER_NAME + `,` + cons.C_EMAIL + `,` + cons.C_PASSWORD + ` From ` + cons.TA_USER +
		` Where ` + key + ` LIKE ?`

	qSRow, err := db.Query(qS, value)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 190)
	if err != nil {
		fmt.Println(err)
	}
	mkS := new(structs.EmailPass)
	if qSRow.Next() {
		err = qSRow.Scan(&mkS.UserName, &mkS.Email, &mkS.Password)
		if err != nil {
			fmt.Println(err)
			return err, nil
		}
	}
	return err, *&mkS
}

func SelectUserCookies(reqCookies string, db *sql.DB) (error, string) {
	qS := ` Select ` + cons.C_COOKIES + ` From ` + cons.TA_USER + ` Where ` + cons.C_COOKIES + ` LIKE ?`

	qSRow, err := db.Query(qS, reqCookies)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 207)
	if err != nil {
		fmt.Println(err)
		return err, ""
	}

	var dbCookies string
	if qSRow.Next() {
		err = qSRow.Scan(&dbCookies)
		if err != nil {
			fmt.Println(err)
			return err, ""
		}
	}
	return err, dbCookies
}
