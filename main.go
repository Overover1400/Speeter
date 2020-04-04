package main

import (
	"C"
	"MainProcess"
	zarinpal "ZarinPalTest"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Start Main! ")

	http.HandleFunc("/main", MainProcess.RootHandler)               // sets router
	http.HandleFunc("/resultForAll", MainProcess.ProccessForSpleet) // sets router
	http.HandleFunc("/singUp", MainProcess.SingUp)                  // sets router
	http.HandleFunc("/register", MainProcess.Register)              // sets router
	http.HandleFunc("/loginPage", MainProcess.LoginPage)            // sets router
	http.HandleFunc("/myFiles", MainProcess.MyFiles)
	http.HandleFunc("/reqForBuy", MainProcess.ReqForBuy)
	http.HandleFunc("/UploadForPay", MainProcess.UploadForPay)
	http.HandleFunc("/y", Y)
	// sets router
	http.Handle("/", http.FileServer(http.Dir("/home/hamed/Spleeter/src/html/")))
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/home/hamed/Spleeter/src/html/"))))

	http.Handle("/dl/",
		http.StripPrefix("/dl", http.FileServer(http.Dir("/home/hamed"))),
	)
	err := http.ListenAndServe(":4002", nil) // set listen port
	if err != nil {
		//log.Fatal("ListenAndServe: ", err)
		fmt.Println(err)
	}
}

func Payment(res http.ResponseWriter, req *http.Request) {

	zarin, err := zarinpal.NewZarinpal("45bfa1e0-30e4-11e9-869d-005056a205be", false)
	fmt.Println(zarin, "  ", err)
	paymentUrl, authority, statusCod, err := zarin.NewPaymentRequest(1000, "http://localhost:4002/pay", "خرید موزیک", "hamed.m7100@gmail.com", "")
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

func Payed(res http.ResponseWriter, req *http.Request)  {

}

func Y(res http.ResponseWriter,req *http.Request)  {


	for _,cookie:=range req.Cookies(){
		fmt.Println(cookie.Value)
	}

}