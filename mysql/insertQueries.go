package mysql

import (
	"cons"
	"database/sql"
	"fmt"
	"strconv"
)

func InsertUserInfo( userName string, encrptPasword []byte,
	cookies string, UserEmail , folderName string,db *sql.DB) (error,string) {
	qIn := ` Insert into ` + cons.TA_USER + ` (` + cons.C_USER_NAME + `,` + cons.C_PASSWORD + `,` + cons.C_COOKIES + `,` +
		cons.C_EMAIL +`,`+cons.C_FOLDER_NAME+ `) Values(?,?,?,?,?)`

	insertRow, err := db.Exec(qIn, userName, string(encrptPasword), cookies, UserEmail,folderName)
	if err != nil {
		fmt.Println(err)
	}

	id,err :=insertRow.LastInsertId()
	if err != nil {
	fmt.Println(err)
	}

	return err,strconv.FormatInt(id,10)
}

