package MainProcess

import (
	zarinpal "ZarinPalTest"
	"cons"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mysql"
	"net/http"
	"security"
	"strconv"
	"strings"
	"structs"
	"useful"
)

func MainProcess(res http.ResponseWriter, req *http.Request, db *sql.DB,userId string) (string, string) {

	var musicName, folderName string

	qS := ` Select ` + cons.C_USER_LINE + ` From ` + cons.TA_LINE

	row := db.QueryRow(qS)
	var line int
	err := row.Scan(&line)
	if err != nil {
		fmt.Println(err)
	}

	if line == 0 {
		defer UpdateLine(`0`, db)

		UpdateLine(`1`, db)



		var userFolderNameTempOrPerm string
		var userName string
		var coin int
		if userId != "" {
			qS := ` Select ` + cons.C_FOLDER_NAME +`,`+cons.C_USER_NAME +`,`+cons.C_COIN+` From ` + cons.TA_USER + ` Where ` + cons.C_USER_ID + `=?`

			qSRow, err := db.Query(qS, userId)
			defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 68)
			if err != nil {
				fmt.Println(err)
			}
			var userFolderName string
			if qSRow.Next() {
				err = qSRow.Scan(&userFolderName,&userName,&coin)
				if err != nil {
					fmt.Println(err)
				}
				userFolderNameTempOrPerm =userFolderName+`/`
			}
		}

		inputMainMusic := cons.MAIN_MUSIC_PATH + userFolderNameTempOrPerm

		//--Upload files
		musicName, folderName = useful.UploadAudio(req, inputMainMusic)
		fmt.Println(folderName)
		if musicName == cons.R_EXTENTION_NOT_ALLOWED {
			fmt.Fprintln(res, "File format not allowed")
			return "", ""
		} else if musicName == cons.R_SIZE_IS_BIG {
			fmt.Fprintln(res, "File size must smaller then 30 mg")
			return "", ""
		}
		fmt.Println("Upload is Successfully ")
		//return ,""
		//-- Get audio duration (second) and then divide by 10
		audioTimeDuration := useful.SpecifyAudioTimeDuration(inputMainMusic + musicName)
		fmt.Println("Time duration is Successfully(second) : ", audioTimeDuration)

		outPutPathOfSpleeter := cons.OUT_PUT_SPLEETER_PATH + `/` + userFolderNameTempOrPerm

		//-- 40 is second here
		if audioTimeDuration <= 40 {
			err := useful.Spleeter(inputMainMusic+musicName, outPutPathOfSpleeter /*+strconv.Itoa(i)*/, "5stems ", "", "")
			fmt.Println("spleeter is complate !!", err)
			//if err == nil {
			//	useful.AttachAudio("/" + userFolderNameTempOrPerm + `/` + cookies)
			//}
		} else {

			outPutPath := cons.OUT_PUT_SPLIT_MUSIC_PATH + "/" + userFolderNameTempOrPerm

			//-- Split music to a few part
			errCondition, splitOutPutPath := useful.SplitAudio(outPutPath, musicName, inputMainMusic, 10)
			if errCondition == cons.R_FAILED {
				fmt.Println(cons.R_FAILED)
				return "", ""
			}

			sliceOfMusicParts := useful.ListOfFiles(splitOutPutPath)
			fmt.Println("spliter is succesfull : ", len(sliceOfMusicParts))
			//fmt.Fprintln(res, "separate music to : ", len(sliceOfMusicParts), " part")

			fmt.Println(` زمان تقریبی `, len(sliceOfMusicParts)*30, `ثانیه`)
			//fmt.Fprintln(res, ` زمان تقریبی `, len(sliceOfMusicParts)*50, `ثانیه`)

			removExtentionMusicName:=strings.ReplaceAll(musicName, ".mp3","")


		for i, v := range sliceOfMusicParts {
				fmt.Println(v)
				fmt.Println(splitOutPutPath)
				err := useful.Spleeter(splitOutPutPath+v, outPutPathOfSpleeter+removExtentionMusicName /*+strconv.Itoa(i)*/, "5stems ",
					"", "")
				if err != nil {
					fmt.Println(err)
					return "", ""
				}
				fmt.Println("One of th Spleeter Part is successfully : is on ", i)
				//fmt.Fprintln(res, "One of th Spleeter Part is successfully : is on ", i)

			}


			useful.AttachAudio("/" + userFolderNameTempOrPerm+ removExtentionMusicName)
			fmt.Println("Attach files is successfully")
			//fmt.Fprintln(res, "Attach files is successfully")


			qU :=` Update `+cons.TA_USER+` Set `+cons.C_COIN +`=`+cons.C_COIN+`-1`+` Where `+
				cons.C_USER_ID +`=?`

			_,err:=db.Exec(qU,userId)
			if err != nil {
				fmt.Println(err)
			}


		command := `rm -r ` + splitOutPutPath
		err, r, n := useful.ShellOut(command)
		fmt.Println(command)
		if err != nil {
			fmt.Println(err)
			fmt.Println(r)
			fmt.Println(n)
		}

		command2 := `rm -r `+outPutPathOfSpleeter+removExtentionMusicName+`/`
		err, _, _ = useful.ShellOut(command2)
		if err != nil {
			fmt.Println(err)

		}
		useful.ParsHtmFiles(res,"MyFiles",structs.EmailPass{UserName:userName,Coin:coin})
		}
	} else {
		return "", "busy"
	}
	return musicName, ""
}

func UpdateLine(num string, db *sql.DB) {
	qU2 := ` Update ` + cons.TA_LINE + ` Set ` + cons.C_USER_LINE + `=` + num
	_, err := db.Exec(qU2)
	if err != nil {
		fmt.Println(err)
	}
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	useful.ParsHtmFiles(res, "main", nil)
}

// this function responsible to upload file and process of spleeter
// also if server was busy and the user was perm that music will go to waiting folder to process later
func ProcessSpleetPermUser(res http.ResponseWriter, req *http.Request) {

	reqUserId := req.FormValue(cons.C_USER_ID)
	reqMusicFolder := req.FormValue(cons.C_MUSIC_NAME)
	reqUserFolder := req.FormValue(cons.C_USER_FOLDER)

	if reqUserId != "" {
	fmt.Println(reqUserId)
	fmt.Println(reqMusicFolder)
	fmt.Println(reqUserFolder)
	if reqUserId != "" && reqMusicFolder != "" && reqUserFolder != "" {

		musicAddress := reqUserFolder + `/` + reqMusicFolder
		useful.ParsHtmFiles(res, "resultForAll", structs.SendValue{MusicName: musicAddress})
		return

	}

	var (
		musicName, busy, folderName string
	)


		db, err := useful.OpenDatabases()
		defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 212)
		if err != nil {
			fmt.Println(err)
		}

		qS := ` Select ` + cons.C_FOLDER_NAME + ` From ` + cons.TA_USER + ` Where ` + cons.C_USER_ID + `=?`

		qSRow, err := db.Query(qS, reqUserId)
		defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 221)
		if err != nil {
			fmt.Println(err)
		}

		if qSRow.Next() {
			err = qSRow.Scan(&folderName)
			if err != nil {
				fmt.Println(err)
			}
		}


	multfile, _, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}

	if multfile != nil {

		qS := ` Select ` + cons.C_COIN + ` From ` + cons.TA_USER +
			` Where ` + cons.C_USER_ID + `=?`

		qSRow, err := db.Query(qS, reqUserId)
		defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 188)
		if err != nil {
			fmt.Println(err)
		}

		var coin int
		if qSRow.Next() {
			err = qSRow.Scan(&coin)
			if err != nil {
				fmt.Println(err)
			}
		}

		if reqUserId != "" {
			if coin > 0 {
				musicName, busy = MainProcess(res, req, db, reqUserId)
			}else {
				fmt.Fprintln(res,"Your Coin is not enough !")
			}
		}
	}

	if busy != "" {

			fmt.Fprintln(res, " Server is busy and will take several minutes ! ")

		//Select for all waiting line musics if was bigger then 20 dont save any music
		allCountOfWaitingLine := mysql.SelectCountOfWaitingLine(db, "")
		if allCountOfWaitingLine >= 20 {
			fmt.Fprintln(res, "The Musics are a lot please try later !")
			return
		}

		//Select for user waiting line musics if was bigger then 3 dont save any music
		countOfWaitingLine := mysql.SelectCountOfWaitingLine(db, reqUserId)

		if countOfWaitingLine > 3 {
			fmt.Fprintln(res, "Your waiting line is more please wait for those")
			return
		}

		folderName := useful.RandomString(15)

		qIns := ` Insert into ` + cons.TA_WAITING_LINE + ` (` + cons.C_USER_ID + `,` + cons.C_DATE + `,` + cons.C_FOLDER_NAME + `) ` +
			` values(?,unix_timestamp(),'` + folderName + `')`

		_, err = db.Exec(qIns, reqUserId)
		if err != nil {
			fmt.Println(err)
		}

		// save file in waitingLine folder
		useful.SaveFile(req, folderName)
		return
	} else {

		removeFormatOfMusic := useful.RemoveExtension(musicName)
		var folderNamePlusSlash string
		if folderName != "" {
			folderNamePlusSlash = folderName + `/` + removeFormatOfMusic
		} else {
			folderNamePlusSlash = removeFormatOfMusic
		}
		useful.ParsHtmFiles(res, "resultForAll", structs.SendValue{MusicName: folderNamePlusSlash})
	}
}
	//fmt.Fprintln(res, "finsih the job")
}

func ProcessSpleetTempUser(res http.ResponseWriter,req *http.Request)  {



	cookies := req.Cookies()[2].Value

	if cookies == "" {
		cookies = useful.RandomString(15)
	}


	//--Upload files
	inputPath :=cons.TEMPO_INPUT_PATH +`/`+cookies+`/`
	musicName, _ := useful.UploadAudio(req, inputPath)

	removExtentionMusicName:=strings.ReplaceAll(musicName, ".mp3","")

	outputPath :=cons.TEMPO_OUTPUT_PATH +`/`+cookies+`/`
	err := useful.Spleeter(inputPath+musicName, outputPath /*+strconv.Itoa(i)*/, "5stems ",
		"5", "1")
	if err != nil {
		fmt.Println(err)
	}

fmt.Println("Finished -----")

useful.ParsHtmFiles(res,"resultForAll",structs.SendValue{MusicName:"temporary/"+cookies+`/`+removExtentionMusicName})
}




func SingUp(res http.ResponseWriter, req *http.Request) {
	err, haveCookies, userId := security.CheckCookies(req, res)
	if err != nil {
		fmt.Println(err)
		return
	}

	if haveCookies {
		db, err := useful.OpenDatabases()
		defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 185)
		if err != nil {
			fmt.Println(err)
		}

		qS := ` Select ` + cons.C_USER_NAME + ` From ` + cons.TA_USER +
			` Where ` + cons.C_USER_ID + `=` + userId

		qSRow, err := db.Query(qS)
		defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 197)
		if err != nil {
			fmt.Println(err)
		}

		var userName string
		if qSRow.Next() {
			err = qSRow.Scan(&userName)
			if err != nil {
				fmt.Println(err)
			}
		}

		qS2 := ` Select ` + cons.C_COIN + ` From ` + cons.TA_USER_BOUGHT + ` Where ` + cons.C_USER_ID + `=` + userId

		qS2Row, err := db.Query(qS2)
		defer useful.HandleCloseableErrorClient(qS2Row, "MainNeededFile.go", 214)
		if err != nil {
			fmt.Println(err)
		}

		var coin int
		if qS2Row.Next() {
			err = qS2Row.Scan(&coin)
			if err != nil {
				fmt.Println(err)
			}
		}

		useful.ParsHtmFiles(res, "MyFiles", structs.EmailPass{UserName: userName, UserId: userId, Coin: coin})
		return
	}
	useful.ParsHtmFiles(res, "singUp", nil)
}

func MyFiles(res http.ResponseWriter, req *http.Request) {
	fmt.Println("start myFiles")
	UserEmail := req.FormValue(cons.C_EMAIL)
	password := req.FormValue(cons.C_PASSWORD)
	//reqUserName := req.FormValue(cons.C_USER_NAME)
	//reqUserId := req.FormValue(cons.C_USER_ID)
	//var cookies string
	//
	//if len(req.Cookies())!=0{
	//	cookies = req.Cookies()[0].String()
	//}

	//fmt.Println("-----", UserEmail, password)
	//fmt.Println("1++++", req.URL)
	//fmt.Println("2++++", req.URL.Query().Get("Authority"))
	//fmt.Println("3++++", req.URL.Query().Get("Status"))

	UserEmail = "hamed@gmail.com"
	password = "123"

	if UserEmail != "" && password != "" {
		db, err := useful.OpenDatabases()
		defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 181)
		if err != nil {
			fmt.Println(err)
		}

		err, emlPass := mysql.SelectUserAuth(UserEmail, "", db)
		if err != nil {
			fmt.Println(err)
			return
		}

		if emlPass.Email == UserEmail {

			//err = bcrypt.CompareHashAndPassword([]byte(emlPass.Password), []byte(password))
			//if err != nil {
			//	fmt.Println(err)
			//	useful.ParsHtmFiles(res, "login", nil)
			//	return
			//}

			musics := useful.ListOfFiles(cons.FINAL_OUT_PUT + `/` + emlPass.FolderName)
			emlPass.Folders = musics
			fmt.Println("---", req.FormValue("giveResult"))
			if req.FormValue("giveResult") != "" {
				encode := json.NewEncoder(res)
				if err := encode.Encode(emlPass); err != nil {
					fmt.Println(err)
				}
			} else {
				qS := ` Select ` + cons.C_COIN + ` From ` + cons.TA_USER +
					` Where ` + cons.C_USER_ID + `=` + emlPass.UserId

				qSRow, err := db.Query(qS)
				defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 288)
				if err != nil {
					fmt.Println(err)
				}

				var coin int
				if qSRow.Next() {
					err = qSRow.Scan(&coin)
					if err != nil {
						fmt.Println(err)
					}
				}

				emlPass.Status = true
				emlPass.Coin = coin
				useful.ParsHtmFiles(res, "MyFiles", emlPass)
				return
			}
			//templates, err := template.ParseFiles(cons.HTLM_FOLDER + "MyFiles.html")
			//if err != nil {
			//	fmt.Println(err)
			//}
			//err = templates.ExecuteTemplate(res, "MyFiles.html", nil)
			////err=templates.ExecuteTemplate(res,"PhonicMind _ My files.html",nil)
			//if err != nil {
			//	fmt.Println(err)
			//}

		} else {
			fmt.Fprintln(res, " noting here !!")
			fmt.Println(" noting here !!")
		}
	} else {
		fmt.Fprintln(res, "eml pass is empty")
	}
}

func ReqForBuy(res http.ResponseWriter, req *http.Request) {
	fmt.Println("hello price")
	reqUserId := req.FormValue(cons.C_USER_ID)
	reqBuy := req.FormValue(cons.BUY)

	if reqUserId == "" {
		useful.ParsHtmFiles(res, "login", nil)
		return
	}

	var amount int

	switch reqBuy {
	case "1":
		amount = 20000
	case "5":
		amount = 100000
	case "10":
		amount = 200000
	case "20":
		amount = 400000
	}

	zarin, err := zarinpal.NewZarinpal("45bfa1e0-30e4-11e9-869d-005056a205be", false)
	fmt.Println(zarin, "  ", err)
	paymentUrl, authority, statusCod, err := zarin.NewPaymentRequest(amount, "http://localhost:4002/login", "خرید موزیک", "hamed.m7100@gmail.com", "")
	fmt.Println("1---", paymentUrl, authority, statusCod, err)

	verification, refID, statusCod2, err := zarin.PaymentVerification(amount, authority)
	fmt.Println("2----", verification, refID, statusCod2, err)

	if statusCod2 != 100 {
		http.Redirect(res, req, paymentUrl, http.StatusTemporaryRedirect)

	} else {
		fmt.Fprintln(res, "3----- authority is successfully ")
		fmt.Println("3-----  authority is successfully ")
	}
}

func LoginPage(res http.ResponseWriter, req *http.Request) {

	err, haveCookies, userId := security.CheckCookies(req, res)
	if err != nil {
		fmt.Println(err)
		return
	}

	if haveCookies {
		db, err := useful.OpenDatabases()
		defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 183)
		if err != nil {
			fmt.Println(err)
		}

		qS := ` Select ` + cons.C_COIN +`,`+cons.C_USER_NAME + ` From ` + cons.TA_USER +
			` Where ` + cons.C_USER_ID + `=?`

		qSRow, err := db.Query(qS, userId)
		defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 188)
		if err != nil {
			fmt.Println(err)
		}
		var sendValue structs.EmailPass
		var coin int
		var userName string
		if qSRow.Next() {
			err = qSRow.Scan(&coin,&userName)
			if err != nil {
				fmt.Println(err)
			} else {
				sendValue.Status = true
			}
		} else {
			coin = 0
		}

		sendValue.Coin = coin
		sendValue.UserName = userName
		useful.ParsHtmFiles(res, "MyFiles", sendValue)
		return
	}
	useful.ParsHtmFiles(res, "login", nil)
	return
}

func Register(res http.ResponseWriter, req *http.Request) {

	userName := req.FormValue(cons.C_USER_NAME)
	UserEmail := req.FormValue(cons.C_EMAIL)
	password := req.FormValue(cons.C_PASSWORD)
	password2 := req.FormValue(cons.C_PASSWORD2)

	db, err := useful.OpenDatabases()
	defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 206)
	if err != nil {
		fmt.Println(err)
	}

	err, emlPass1 := mysql.SelectUserAuth("", userName, db)
	if err != nil {
		fmt.Println(err)
		return
	}

	if emlPass1.UserName == userName {
		useful.ParsHtmFiles(res, "singUp",
			structs.SendMessage{Msg: "User Name is duplicate ", State: true, UserInfo: structs.EmailPass{Email: UserEmail, UserName: userName}})
		return
	}

	err, emlPass2 := mysql.SelectUserAuth(UserEmail, "", db)
	if err != nil {
		fmt.Println(err)
		return
	}

	if emlPass2.Email == UserEmail {
		useful.ParsHtmFiles(res, "singUp",
			structs.SendMessage{Msg: "Email is exsist", State: true, UserInfo: structs.EmailPass{Email: UserEmail, UserName: userName}})
		return
	}

	if password != password2 {
		var messages = structs.SendMessage{Msg: " Second password is not match! ", State: true}
		useful.ParsHtmFiles(res, "singUp", messages)
		return
	}

	var cookies string
	if len(req.Cookies()) != 0 {
		cookies = req.Cookies()[0].Value
	}

	encrptPasword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		fmt.Println(err)
	}

	folderName := useful.RandomString(20)

	err, userId := mysql.InsertUserInfo(userName, encrptPasword, cookies, UserEmail, folderName, db)
	if err != nil {
		fmt.Println(err)
	}

	command := `mkdir ` + cons.OUT_PUT_SPLEETER_PATH + `/` + folderName
	err, _, _ = useful.ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}

	useful.ParsHtmFiles(res, "welcome", structs.EmailPass{UserId: userId, UserName: userName,
		Email: UserEmail, Password: password, FolderName: folderName})
}

func UploadForPay(res http.ResponseWriter, req *http.Request) {

	db, err := useful.OpenDatabases()
	defer useful.HandleCloseableErrorClient(db, "MainNeededFile.go", 285)
	if err != nil {
		fmt.Println(err)
		return
	}

	reqCookies := req.Cookies()[0].Value

	qS := ` Select ` + cons.C_USER_ID + ` From ` + cons.TA_USER + ` Where ` + cons.C_COOKIES + ` LIKE ?`

	qSRow, err := db.Query(qS, reqCookies)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 297)
	if err != nil {
		fmt.Println(err)
		return
	}

	var userID string
	if qSRow.Next() {
		err = qSRow.Scan(&userID)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

	randNumber := useful.RandomNumber(100000, 10000000)

	inputMainMusic := cons.MAIN_MUSIC_PATH + strconv.Itoa(randNumber) + `/` + userID + `/`

	//for i:=0;i<10;i++ {
	//	time.Sleep(time.Second * 5)
	//	fmt.Println("salam")
	//}
	//	return
	//--Upload files
	musicName, folderName := useful.UploadAudio(req, inputMainMusic)
	if musicName == cons.R_EXTENTION_NOT_ALLOWED {
		fmt.Fprintln(res, "File format not allowed")
		return
	} else if musicName == cons.R_SIZE_IS_BIG {
		fmt.Fprintln(res, "File size must smaller then 30 mg")
		return
	}

	qIns := ` insert into ` + cons.TA_USER_BOUGHT + ` (` + cons.C_USER_ID + `,` + cons.C_ZARINPAL_AUTHORITY + `,` +
		cons.C_STATUS_INT + `,` + cons.C_STATUS_STRING + `,` + cons.C_AMOUNT + `,` + cons.C_MUSIC_NAME + `) Values(?,?,?,?,?,?)`

	_, err = db.Exec(qIns, userID, `''`, 0, `''`, 0, musicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	qS2 := ` Select ` + cons.C_STATUS_STRING + ` From ` + cons.TA_USER_BOUGHT + ` Where ` + cons.C_USER_ID + `=?`

	qSRow2, err := db.Query(qS2, userID)
	defer useful.HandleCloseableErrorClient(qSRow, "MainNeededFile.go", 342)
	if err != nil {
		fmt.Println(err)
		return
	}

	var status string
	if qSRow2.Next() {
		err = qSRow2.Scan(&status)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var showMusic structs.SendMessage
	showMusic.State = true

	if status == "OK" {
		showMusic.Msg = ""
	} else {
		showMusic.Msg = folderName
		showMusic.FolderName = musicName
	}
	fmt.Println(showMusic)
	useful.ParsHtmFiles(res, "MyFiles", showMusic)
	return
}

func Price(res http.ResponseWriter, req *http.Request) {
	useful.ParsHtmFiles(res, "price", nil)
}

type Savary struct {
	Tier string
	Door int
}

type Bary struct {
	Tier string
	Door int
}

type Motory struct {
	Tier string
	Door int
}

func (s Savary) ChangeTiers(summer bool) string {
	if summer {
		s.Tier = "tabestani"
		return "dont changed"
	} else {
		s.Tier = "zemestani"
		return "changed"
	}
}

func (b Bary) ChangeTiers(summer bool) string {
	if summer {
		b.Tier = "tabestani"
		return "dont changed"
	} else {
		b.Tier = "zemestani"
		return "changed"
	}
}

type TiereFasly interface {
	ChangeTiers(summer bool) string
}

func ZemestanOrtabbestan(fasle string) TiereFasly {

	var SavaryCondetion Savary
	var BaryCondition Bary
	//var MotoryCondition Motory

	if fasle == "zemestan" {
		return SavaryCondetion
	} else if fasle == "tabestan" {
		return BaryCondition
	} else {
		//return Motor``yCondition
	}
	return nil
}
