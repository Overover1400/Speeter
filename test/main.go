package main

import (
	"fmt"
	"os"
)

func main() {
	// Create an HTTP server that listens on port 8000
	//http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	ctx := r.Context()
	//	// This prints to STDOUT to show that processing has started
	//	fmt.Fprint(os.Stdout, "processing request\n")
	//	// We use `select` to execute a peice of code depending on which
	//	// channel receives a message first
	//	select {
	//	case <-time.After(2 * time.Second):
	//		// If we receive a message after 2 seconds
	//		// that means the request has been processed
	//		// We then write this as the response
	//		w.Write([]byte("request processed"))
	//	case <-ctx.Done():
	//		// If the request gets cancelled, log it
	//		// to STDERR
	//		fmt.Fprint(os.Stderr, "request cancelled\n")
	//	}
	//}))

	//fmt.Println((400 * 10) / 400)
//s:=sha256.Sum256([]byte("hi"))
//var slicc []string
//for _,v :=range s {
//	if string(v)!= "" {
//		slicc = append(slicc, string(v))
//	}
//
//}
//	fmt.Println(slicc)


fmt.Println("ss "+os.Getenv("DB_USERNAME"))

}