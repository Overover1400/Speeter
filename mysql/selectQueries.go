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
	qS := ` Select ` +cons.C_USER_ID +`,`+ cons.C_USER_NAME + `,` +
		cons.C_EMAIL + `,` + cons.C_PASSWORD +`,`+cons.C_FOLDER_NAME+ ` From ` + cons.TA_USER +
		` Where ` + key + ` LIKE ?`

	qSRow, err := db.Query(qS, value)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 190)
	if err != nil {
		fmt.Println(err)
	}
	mkS := new(structs.EmailPass)
	if qSRow.Next() {
		err = qSRow.Scan(&mkS.UserId,&mkS.UserName, &mkS.Email, &mkS.Password,&mkS.FolderName)
		if err != nil {
			fmt.Println(err)
			return err, nil
		}
	}
	return err, *&mkS
}

func SelectUserCookies(reqCookies string, db *sql.DB) (error, string,string) {
	qS := ` Select ` + cons.C_COOKIES +`,`+cons.C_USER_ID+ ` From ` + cons.TA_USER + ` Where ` + cons.C_COOKIES + ` LIKE ?`

	qSRow, err := db.Query(qS, reqCookies)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 207)
	if err != nil {
		fmt.Println(err)
		return err, "",""
	}

	var dbCookies,userId string
	if qSRow.Next() {
		err = qSRow.Scan(&dbCookies,&userId)
		if err != nil {
			fmt.Println(err)
			return err, "",""
		}
	}
	return err, dbCookies,userId
}

func SelectCountOfWaitingLine(db *sql.DB, reqUserId string) int {
	var whereCondition string
	if reqUserId != "" {
		whereCondition = ` Where ` + cons.C_USER_ID + `=` + reqUserId
	}

	qS := ` Select count(` + cons.C_ID + `) From ` + cons.TA_WAITING_LINE + whereCondition

	qSRow, err := db.Query(qS)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 178)
	if err != nil {
		fmt.Println(err)
	}

	var allCountOfWaitingLine int
	if qSRow.Next() {
		err = qSRow.Scan(&allCountOfWaitingLine)
		if err != nil {
			fmt.Println(err)
		}
	}
	return allCountOfWaitingLine
}

