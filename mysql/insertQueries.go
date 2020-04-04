package mysql

import (
	"cons"
	"database/sql"
	"fmt"
)

func InsertUserInfo( userName string, encrptPasword []byte, cookies string, UserEmail string,db *sql.DB) error {
	qIn := ` Insert into ` + cons.TA_USER + ` (` + cons.C_USER_NAME + `,` + cons.C_PASSWORD + `,` + cons.C_COOKIES + `,` +
		cons.C_EMAIL + `) Values(?,?,?,?)`

	_, err := db.Exec(qIn, userName, string(encrptPasword), cookies, UserEmail)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

