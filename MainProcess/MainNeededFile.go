package MainProcess

import (
	zarinpal "ZarinPalTest"
	"cons"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"mysql"
	"net/http"
	"security"
	"strconv"
	"structs"
	"useful"
)

func MainProcess(res http.ResponseWriter, req *http.Request) (string,string){

	db, err := useful.OpenDatabases()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(res, " the Connection error")
		return "",""
	}


	var musicName,folderName,busy string

	qS := ` Select ` + cons.C_USER_LINE + ` From ` + cons.TA_LINE

	row := db.QueryRow(qS)
	var line int
	err = row.Scan(&line)
	if err != nil {
		fmt.Println(err)
	}

	if line == 0 {
		defer UpdateLine(`0`, db)

	UpdateLine(`1`,db)

		command :=`rm -r `+cons.MAIN_MUSIC_PATH+`/*`
		err,_,_=useful.ShellOut(command)
		if err != nil {
		fmt.Println(err)
		}

		command2 :=`rm -r /home/hamed/finalOutPut/*`
		err,_,_=useful.ShellOut(command2)
		if err != nil {
			fmt.Println(err)
		}


		inputMainMusic := cons.MAIN_MUSIC_PATH

		//for i:=0;i<10;i++ {
		//	time.Sleep(time.Second * 5)
		//	fmt.Println("salam")
		//}
		//	return ,""
		//--Upload files
		musicName, folderName = useful.UploadAudio(req, inputMainMusic, false)
		fmt.Println(folderName)
		if musicName == cons.R_EXTENTION_NOT_ALLOWED {
			fmt.Fprintln(res, "File format not allowed")
			return "",""
		} else if musicName == cons.R_SIZE_IS_BIG {
			fmt.Fprintln(res, "File size must smaller then 30 mg")
			return "",""
		}
		fmt.Println(res, "Upload is Successfully ")
		//return ,""
		//-- Get audio duration (second) and then divide by 10
		audioTimeDuration := useful.SpecifyAudioTimeDuration(inputMainMusic + musicName)
		fmt.Println(res, "Time duration is Successfully(second) : ", audioTimeDuration)

		outPutPathOfSpleeter := cons.OUT_PUT_SPLEETER_PATH

		//-- 40 is second here
		if audioTimeDuration <= 40 {
			err := useful.Spleeter(inputMainMusic+musicName, outPutPathOfSpleeter /*+strconv.Itoa(i)*/, "5stems ", "", "")
			fmt.Println("111", err)
			if err == nil {
				useful.AttachAudio()
			}
		} else {

			outPutPath := cons.OUT_PUT_SPLIT_MUSIC_PATH

			//-- Split music to a few part
			errCondition, splitOutPutPath := useful.SplitAudio(outPutPath, musicName, inputMainMusic, 10)
			if errCondition == cons.R_FAILED {
				fmt.Println(cons.R_FAILED)
				return "",""
			}

			sliceOfMusicParts := useful.ListOfFiles(splitOutPutPath)
			fmt.Println("spliter is succesfull : ", len(sliceOfMusicParts))
			fmt.Fprintln(res, "separate music to : ", len(sliceOfMusicParts), " part")

			fmt.Println(` زمان تقریبی `, len(sliceOfMusicParts)*30, `ثانیه`)
			fmt.Fprintln(res, ` زمان تقریبی `, len(sliceOfMusicParts)*50, `ثانیه`)

			for i, v := range sliceOfMusicParts {
				err := useful.Spleeter(splitOutPutPath+v, outPutPathOfSpleeter /*+strconv.Itoa(i)*/, "5stems ",
					"", "")
				if err != nil {
					fmt.Println(err)
					return "",""
				}
				fmt.Println("One of th Spleeter Part is successfully : is on ", i)
				fmt.Fprintln(res, "One of th Spleeter Part is successfully : is on ", i)

			}

			useful.AttachAudio()
			fmt.Println("Attach files is successfully")
			fmt.Fprintln(res, "Attach files is successfully")
			//TODO remove al directories

		}
	} else {
		busy = "busy"
		return "",busy
	}
	return musicName,""
}

func UpdateLine(num string,db *sql.DB) {
	qU2 := ` Update ` + cons.TA_LINE + ` Set ` + cons.C_USER_LINE + `=`+num
	_, err := db.Exec(qU2)
	if err != nil {
		fmt.Println(err)
	}
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	useful.ParsHtmFiles(res, "main", nil)
}

func ProccessForSpleet(res http.ResponseWriter, req *http.Request) {

var musicName,busy string
	multfile, _, err := req.FormFile("file")
	if err != nil {
	fmt.Println(err)
	}
	if multfile != nil {
		musicName,busy=MainProcess(res, req)
		fmt.Println("OKKK")
	}
	type SendValue struct {
		MusicName string
	}
	if busy != "" {
		fmt.Fprintln(res, " Server is busy and will take several minutes ! ")
	}else {
		removeFormatOfMusic :=useful.RemoveExtention(musicName)
		useful.ParsHtmFiles(res, "resultForAll", SendValue{MusicName: removeFormatOfMusic})
	}
	//fmt.Fprintln(res, "finsih the job")
}

func SingUp(res http.ResponseWriter, req *http.Request) {
	err, haveCookies := security.CheckCookies(req, res)
	if err != nil {
		fmt.Println(err)
		return
	}

	if haveCookies {
		useful.ParsHtmFiles(res, "MyFiles", nil)
		return
	}
	useful.ParsHtmFiles(res, "singUp", nil)
}

func MyFiles(res http.ResponseWriter, req *http.Request) {
	UserEmail := req.FormValue(cons.C_EMAIL)
	password := req.FormValue(cons.C_PASSWORD)
	//var cookies string
	//
	//if len(req.Cookies())!=0{
	//	cookies = req.Cookies()[0].String()
	//}

	//fmt.Println("-----", UserEmail, password)
	//fmt.Println("1++++", req.URL)
	//fmt.Println("2++++", req.URL.Query().Get("Authority"))
	//fmt.Println("3++++", req.URL.Query().Get("Status"))

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

		err = bcrypt.CompareHashAndPassword([]byte(emlPass.Password), []byte(password))
		if err != nil {
			fmt.Println(err)
			useful.ParsHtmFiles(res, "login", nil)
			return
		}

		templates, err := template.ParseFiles(cons.HTLM_FOLDER + "MyFiles.html")
		if err != nil {
			fmt.Println(err)
		}
		err = templates.ExecuteTemplate(res, "MyFiles.html", nil)
		//err=templates.ExecuteTemplate(res,"PhonicMind _ My files.html",nil)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Fprintln(res, " noting here !!")
		fmt.Println(" noting here !!")

	}

}

func ReqForBuy(res http.ResponseWriter, req *http.Request) {
	zarin, err := zarinpal.NewZarinpal("45bfa1e0-30e4-11e9-869d-005056a205be", false)
	fmt.Println(zarin, "  ", err)
	paymentUrl, authority, statusCod, err := zarin.NewPaymentRequest(1000, "http://localhost:4002/login", "خرید موزیک", "hamed.m7100@gmail.com", "")
	fmt.Println("1---", paymentUrl, authority, statusCod, err)

	verification, refID, statusCod2, err := zarin.PaymentVerification(1000, authority)
	fmt.Println("2----", verification, refID, statusCod2, err)

	if statusCod2 != 100 {
		http.Redirect(res, req, paymentUrl, http.StatusTemporaryRedirect)

	} else {
		fmt.Fprintln(res, "3----- authority is successfully ")
		fmt.Println("3-----  authority is successfully ")
	}
}

func LoginPage(res http.ResponseWriter, req *http.Request) {

	err, haveCookies := security.CheckCookies(req, res)
	if err != nil {
		fmt.Println(err)
		return
	}

	if haveCookies {
		useful.ParsHtmFiles(res, "MyFiles", nil)
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

	err = mysql.InsertUserInfo(userName, encrptPasword, cookies, UserEmail, db)
	if err != nil {
		fmt.Println(err)
	}

	var sendEmlPass = structs.EmailPass{Email: UserEmail, Password: password}
	useful.ParsHtmFiles(res, "welcome", sendEmlPass)
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
	musicName, folderName := useful.UploadAudio(req, inputMainMusic, true)
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
