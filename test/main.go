package main

import (
	"fmt"
	"useful"
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
	//alphabet:=sha256.Sum256([]byte("hi"))
	//var slicc []string
	//for _,v :=range alphabet {
	//	if string(v)!= "" {
	//		slicc = append(slicc, string(v))
	//	}
	//
	//}
	//	fmt.Println(slicc)

	//fmt.Println("ss "+os.Getenv("DB_USERNAME"))

	//res,err:=http.Get("https://pouland.ir/login/?username=ha'med';&password=123#")
	//if err != nil {
	//fmt.Println("err : ",err)
	//}
	//
	//io,err:=ioutil.ReadAll(res.Request.Form)
	//if err != nil {
	//fmt.Println(err)
	//}

	//fmt.Println(string(io))
	//fmt.Println(res.Request.Form)
	//fmt.Println(res.Request.Body)
	//fmt.Println("-----------",res.Header)

	//fmt.Println(string(""))
	//
	//	fmt.Println("frist -----",os.Getenv("CONDA_DEFAULT_ENV"))
	//
	//
	//	fmt.Println("aaa",err)
	//
	//	fmt.Println("second ----- ",os.Getenv("CONDA_DEFAULT_ENV"))
	//
	//	err=useful.Spleeter(cons.MAIN_MUSIC_PATH+"audio_example.mp3", cons.OUT_PUT_SPLEETER_PATH,"5stems ","","")
	//	fmt.Println("out error  ",err)

	//str := `audio_example.msdasjkuiodp3`
	//var sv = true
	//var j = 1
	//for i := 1; i < len(str); i++ {
	//
	//	if string(str[len(str)-i]) != `.` && sv {
	//		j++
	//	}
	//	if string(str[len(str)-i]) == `.` {
	//		sv = false
	//	}
	//
	//}
	//fmt.Println(str[:len(str)-j])
	//fmt.Println(useful.RandomString(20))
//os.Chmod("/etc/hosts",os.``)
//	err:=ioutil.WriteFile("/etc/hosts", []byte("helloo"), os.ModeExclusive)
//
//	fmt.Println(err)
//	e,err:=os.Stat("/etc/hosts")
//	fmt.Println(e.Mode(),err)

//f:=`/home/hamed.txt`
f:=`nano /etc/hosts`
//file,err:=os.OpenFile(f,os.O_WRONLY,os.ModePerm)
//fmt.Println(err)
//	reader,err:=ioutil.ReadAll(file)
//fmt.Println(string(reader),err)
//
//	writer,err :=file.WriteAt([]byte("chetory"),10)
//	fmt.Println(writer,err)


err,r,_:=useful.ShellOut(f)
fmt.Println(err,r)

}


