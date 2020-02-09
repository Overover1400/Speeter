package main

import (
	"C"
	MainProcess2 "MainProcess"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Start Main! ")

	http.HandleFunc("/", MainProcess2.RootHandler) // sets router
	err := http.ListenAndServe(":4001", nil)       // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


