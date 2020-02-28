package MainProcess

import (
	zarinpal "ZarinPalTest"
	"cons"
	"fmt"
	"html/template"
	"net/http"
	"structs"
	"time"
	"useful"
	"golang.org/x/crypto/bcrypt"
)

func MainProcess(res http.ResponseWriter, req *http.Request) {

	//	db,err :=useful.OpenDatabases()
	//	if err != nil {
	//	fmt.Println(err)
	//	fmt.Fprintln(res," the Connection error")
	//		return
	//	}
	//
	//
	//	qS :=` Select `+cons.C_CONDITION+` From `+cons.TA_USER_LINE+` Order By ` +cons.C_ID+` Desc Limit 1`
	//
	//	qSRow,err:=db.Query(qS)
	//defer qSRow.Close()

	inputMainMusic := cons.MAIN_MUSIC_PATH

	//for i:=0;i<10;i++ {
	//	time.Sleep(time.Second * 5)
	//	fmt.Println("salam")
	//}
	//	return
	//--Upload files
	musicName := useful.UploadAudio(req, inputMainMusic)
	if musicName == cons.R_EXTENTION_NOT_ALLOWED {
		fmt.Fprintln(res, "File format not allowed")
		return
	} else if musicName == cons.R_SIZE_IS_BIG {
		fmt.Fprintln(res, "File size must smaller then 30 mg")
		return
	}
	fmt.Fprintln(res, "Upload is Successfully ")
	return
	//-- Get audio duration (second) and then divide by 10
	audioTimeDuration := useful.SpecifyAudioTimeDuration(inputMainMusic + musicName)
	fmt.Fprintln(res, "Time duration is Successfully(second) : ", audioTimeDuration)

	outPutPathOfSpleeter := cons.OUT_PUT_SPLEETER_PATH

	//-- 40 is second here
	if audioTimeDuration <= 40 {
		err := useful.Spleeter(inputMainMusic, outPutPathOfSpleeter /*+strconv.Itoa(i)*/, "5stems ", "", "")
		if err == nil {
			useful.AttachAudio()
		}
	} else {

		outPutPath := cons.OUT_PUT_SPLIT_MUSIC_PATH

		//-- Split music to a few part
		errCondition, splitOutPutPath := useful.SplitAudio(outPutPath, musicName, inputMainMusic, 10)
		if errCondition == cons.R_FAILED {
			fmt.Println(cons.R_FAILED)
			return
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
				return
			}
			fmt.Println("One of th Spleeter Part is successfully : is on ", i)
			fmt.Fprintln(res, "One of th Spleeter Part is successfully : is on ", i)

		}

		useful.AttachAudio()
		fmt.Println("Attach files is successfully")
		fmt.Fprintln(res, "Attach files is successfully")
		//TODO remove al directories

	}
}

type Page struct {
	Value    string
	DbStatus bool
	Number   int
}

func RootHandler(res http.ResponseWriter, req *http.Request) {

	//templates, err := template.ParseFiles("/home/hamed/html/main.html")
	//templates, err := template.ParseFiles(cons.HTLM_FOLDER+"main.html")

	templates, err := template.ParseFiles(cons.HTLM_FOLDER+"main.html")
	//templates := template.Must(template.ParseFiles(cons.HTLM_FOLDER+"main2.html"))
	//fmt.Print(os.Getwd())
	if err != nil {
		fmt.Println("Index Template Parse Error: -------------", err)
	}

	n := false
	page := Page{Value: " How are u ?", DbStatus: n}
	//
	//if name:=req.FormValue("value");name != ""{
	//	page.Value=name
	//}

	//page.DbStatus=false
	//page.Number=150

	//err = tmpl.Execute(res, nil)
	if err := templates.ExecuteTemplate(res, "main.html", nil); err != nil {
		fmt.Println("Index Template Execution Error: +++++++++++++", err)
	}

	if req.FormValue("k") == "ok" {
		//http.HandleFunc("/resultForAll",ResultForAll)
		page.DbStatus = true
	}

	//page.Finish=true
	//template.New(cons.HTLM_FOLDER+"resultForAll.html")
	//
	//if err :=templates.ExecuteTemplate(res,"resultForAll.html",page);err != nil {
	//	fmt.Println("Index Template Execution Error: ", err)
	//}
	fmt.Fprintln(res, cons.HTLM_FOLDER+"resultForAll.html")
	fmt.Println("No err here  ")
}

func ResultForAll(res http.ResponseWriter, req *http.Request) {

	templates2, err := template.ParseFiles(cons.HTLM_FOLDER+"resultForAll.html")
	fmt.Println(err, "22222")

	if err := templates2.ExecuteTemplate(res, "resultForAll.html", nil); err != nil {
		fmt.Println("Index Template Execution Error: ", err)
	}
	multfile, _, err := req.FormFile("file")
	if multfile != nil {
		MainProcess(res, req)
		fmt.Println("OKKK")
	}
	time.Sleep(time.Second * 4)
	fmt.Println(err)
	fmt.Fprintln(res, "finsih the job")
}

func SingUp(res http.ResponseWriter, req *http.Request) {

	templates, err := template.ParseFiles(cons.HTLM_FOLDER+"singUp.html")
	if err != nil {
		fmt.Println(err)
	}

	err = templates.ExecuteTemplate(res, "singUp.html", nil)
	if err != nil {
		fmt.Println(err)
	}

}

func MyFiles(res http.ResponseWriter, req *http.Request)  {
	UserEmail :=req.FormValue(cons.C_EMAIL)
	password :=req.FormValue(cons.C_PASSWORD)
	//var cookies string
	//
    //if len(req.Cookies())!=0{
	//	cookies = req.Cookies()[0].String()
	//}
	
	

	fmt.Println("-----",UserEmail,password)
	fmt.Println("1++++",req.URL)
	fmt.Println("2++++",req.URL.Query().Get("Authority"))
	fmt.Println("3++++",req.URL.Query().Get("Status"))

	db,err:=useful.OpenDatabases()
	defer useful.HandleCloseableErrorClient(db,"MainNeededFile.go", 181)
	if err != nil {
		fmt.Println(err)
	}


	qS :=` Select `+cons.C_EMAIL +`,`+cons.C_PASSWORD+` From `+cons.TA_USER+` Where `+cons.C_EMAIL +` LIKE '`+UserEmail+`'`

	fmt.Println(qS)
qSRow,err :=db.Query(qS)
defer useful.HandleCloseableErrorClient(qSRow,"MainNeededFile.go", 190)
if err != nil {
fmt.Println(err)
}
var email string
var pasw string
if qSRow.Next(){
	err=qSRow.Scan(&email,&pasw)
	if err != nil {
	fmt.Println(err)
	}
}

if email != ""{
	if email == UserEmail {

		err=bcrypt.CompareHashAndPassword([]byte(pasw),[]byte(password))
		fmt.Println(err)
		if err != nil {
fmt.Fprintln(res,"pass is incorrect !")
			return
		}
		//fmt.Fprintln(res,"succesfully !",email,password)
		fmt.Println("succesfully !",email,password)

		templates,err:=template.ParseFiles(cons.HTLM_FOLDER+"MyFiles.html")
		if err != nil {
			fmt.Println(err)
		}
		err=templates.ExecuteTemplate(res,"MyFiles.html",nil)
		//err=templates.ExecuteTemplate(res,"PhonicMind _ My files.html",nil)
		if err != nil {
		fmt.Println(err)
		}



	}else {
		fmt.Fprintln(res," noting here !!",email)
		fmt.Println(" noting here !!",email)

	}
}else {
	fmt.Fprintln(res," email is empty !  !!",email)
	fmt.Println(" email is empty !  !!",email)
}

}
func ReqForBuy(res http.ResponseWriter, req *http.Request) {
	fmt.Println("hellooooo")
		zarin, err := zarinpal.NewZarinpal("45bfa1e0-30e4-11e9-869d-005056a205be", false)
		fmt.Println(zarin, "  ", err)
		paymentUrl, authority, statusCod, err := zarin.NewPaymentRequest(1000, "http://localhost:4002/login", "خرید موزیک", "hamed.m7100@gmail.com", "")
		fmt.Println("1---",paymentUrl, authority, statusCod, err)

		verification,refID,statusCod2,err:=zarin.PaymentVerification(1000,authority)
		fmt.Println("2----",verification,refID,statusCod2,err)

		if statusCod2 != 100{
			http.Redirect(res,req,paymentUrl,http.StatusTemporaryRedirect)

		}else {
			fmt.Fprintln(res,"3----- authority is successfully ")
			fmt.Println("3-----  authority is successfully ")
		}
}

func LoginPage(res http.ResponseWriter, req *http.Request)  {

	templates,err:=template.ParseFiles(cons.HTLM_FOLDER+"login.html")
	if err != nil {
		fmt.Println(err)
	}
	err=templates.ExecuteTemplate(res,"login.html", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func Register(res http.ResponseWriter, req *http.Request)  {

	userName:=req.FormValue(cons.C_USER_NAME)
	UserEmail :=req.FormValue(cons.C_EMAIL)
	password :=req.FormValue(cons.C_PASSWORD)
	password2 :=req.FormValue(cons.C_PASSWORD2)

	if password != password2 {
		templates,err:=template.ParseFiles(cons.HTLM_FOLDER+"singUp.html")
		if err != nil {
			fmt.Println(err)
		}

		var messages =structs.SendMessage{Msg:" Second password is incorrect",State:true}

		err=templates.ExecuteTemplate(res,"singUp.html",messages)
		//err=templates.ExecuteTemplate(res,"PhonicMind _ My files.html",nil)
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	var cookies string


	if len(req.Cookies())!=0{
		cookies = req.Cookies()[0].String()
	}

	db,err:=useful.OpenDatabases()
defer useful.HandleCloseableErrorClient(db,"MainNeededFile.go", 206)
	if err != nil {
	fmt.Println(err)
	}

	encrptPasword,err:=bcrypt.GenerateFromPassword([]byte(password),7)
	if err != nil {
	fmt.Println(err)
	}
	qIn :=` Insert into `+cons.TA_USER+` (`+cons.C_USER_NAME +`,`+cons.C_PASSWORD +`,`+cons.C_COOKIES +`,`+
		cons.C_EMAIL+`) Values('`+userName +`','`+string(encrptPasword) +`','`+cookies+`','`+UserEmail+`')`

	_,err=db.Exec(qIn)
	if err != nil {
	fmt.Println(err)
	}

templates,err:=template.ParseFiles(cons.HTLM_FOLDER+"welcome.html")
if err != nil {
fmt.Println(err)
}

type EmailPass struct {
	Email string
	Password string
}


var sendEmlPass EmailPass
	if UserEmail == ""{
		sendEmlPass.Email =""
		sendEmlPass.Password =""
	} else {
		sendEmlPass.Email =UserEmail
		sendEmlPass.Password=password
	}
err=templates.ExecuteTemplate(res,"welcome.html", sendEmlPass)
if err != nil {
fmt.Println(err)
}


}




func Cookies(res http.ResponseWriter, req *http.Request){
	fmt.Println(req.Cookies()[0].String())
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

