package main

import (
	"fmt"
	"net/http"
	"useful"
)

func GetInitiateMainPage(req *http.Request)  {
useful.PrintReqForm("in MainPage",req)

getInitiateMainPage(req)

}

func getInitiateMainPage(req *http.Request) {

	reqSearch :=req.FormValue("search")

	fmt.Println(reqSearch)

}
