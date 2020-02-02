package main

import (
	"bytes"
	"cons"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type Field struct {
	Firstname  string
	Secondname string
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	//tmpl, err := template.ParseFiles("/home/hamed/html/main.html")
	tmpl, err := template.ParseFiles("/home/hamed/Spleeter/src/html/main.html")
	//fmt.Print(os.Getwd())
	if err != nil {
		fmt.Println("Index Template Parse Error: ", err)
	}
	err = tmpl.Execute(res, nil)
	if err != nil {
		fmt.Println("Index Template Execution Error: ", err)
	}
	multfile, _, err := req.FormFile("file")
	if multfile != nil {
		fmt.Println("hello")
		MainProcess(res, req)
	}

	fmt.Fprintln(res, "That was no error !")
}

func main() {
	//command :=`spleeter separate -i ~/Spleeter/src/music/audio_example.mp3 -p spleeter:2stems -o ~/Spleeter/src/outputMusic`
	// its main
	fmt.Println("Start Main! ")


	//fmt.Println((70*10)/400)

	////err:=exec.Command("sh","-c","source ~/.bashrc").Run()
	//err:=exec.Command("sh",command).Run()
	//if err != nil {
	//	fmt.Println("error herer 1 ",err)
	//}
	//	app := "bash"
	//
	//	arg0 := "-e"
	//	//arg1 := "Hello world"
	//
	//		arg1 := "source ~/.bashrc;"
	//	arg2 :=`spleeter separate -i ~/Spleeter/src/music/audio_example.mp3 -p spleeter:2stems -o ~/Spleeter/src/outputMusic`
	//	arg2 :=`spleeter separate -i ~/Spleeter/src/music/audio_example.mp3 ~/Spleeter/src/outputMusic`
	//	arg2 :=`spleeter separate -i audio_example.mp3 -p spleeter:2stems -o output`
	//	err := exec.Command(app, arg0,arg2).Run()
	////	stdout, err := cmd.Output()
	////fmt.Println(string(stdout))
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//fmt.Println("str")
	//	}

	//fmt.Println("successfully")
	//fs := http.FileServer(http.Dir("/home/hamed/output"))
	//http.Handle("/output/", http.StripPrefix("/output/", fs))

	http.HandleFunc("/", RootHandler) // sets router
	//http.HandleFunc("/file", WelcomeHandler)
	err := http.ListenAndServe(":4001", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

const ShellToUse = "bash"

func ShellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	//fmt.Println(cmd.Stdout, cmd.Stderr)
	return err, stdout.String(), stderr.String()
}

func MainProcess(res http.ResponseWriter, req *http.Request) {

	inputMainMusic :=cons.MAIN_MUSIC_PATH

	//--Upload files
	musicName := UploadAudio(req, inputMainMusic)
	if musicName == cons.R_EXTENTION_NOT_ALLOWED {
		fmt.Fprintln(res, "file format not allowed")
		return
	} else if musicName == cons.R_SIZE_IS_BIG {
		fmt.Fprintln(res, "file size is must smaller the n 30 mg")
		return
	}

	//-- Get audio duration (second) and then divide by 10
	audioTimeDuration := SpecifyAudioTimeDuration(inputMainMusic + musicName)
	var divideBy int
	if audioTimeDuration < 40 {
		divideBy = 1
	} else {
		divideBy = (audioTimeDuration * 10) / 400
	}
	if divideBy > 1 {
		//-- Split music to a few part
		outPutPath := `/home/hamed/outSplitMusic`
		_, splitOutPutPath := SplitAudio(outPutPath, musicName, inputMainMusic, divideBy)

		sliceOfMusicParts := ListOfFiles(splitOutPutPath)

		outPutPathOfSpleeter := `/home/hamed/outputSpleeter`
		for _, v := range sliceOfMusicParts {
			Spleeter(splitOutPutPath+`/`+v, outPutPathOfSpleeter, "2tems ", "", "")

			//sliceOfResultMusics :=ListOfFiles(outPutPathOfSpleeter+)
			//for _,v :=range sliceOfResultMusics {
			//
			//}
		}

		AttachAudio(splitOutPutPath)
	}
	//strings.Split()
	//

}

func UploadAudio(req *http.Request, inputMainMusic string) string {

	mf, fh, err := req.FormFile("file")
	defer mf.Close()
	if err != nil {
		fmt.Print(err)
	}

	//-- Check file size and format
	if fh.Size <= 30000000 {


		legalExtension := []string{"mp3", "aac", "wma", "flac", "wav", "aiff"}

		if HasElement(legalExtension, FindAudioFormat(inputMainMusic + fh.Filename)) {

			//create sha for file name
			//ext := strings.Split(fh.Filename, ".")[1]
			//h := sha1.New()
			//_,err=io.Copy(h, mf)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//	fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
			//fmt.Println(fname)

			//u, _ := url.Parse("/home/hamed/music/")
			//u, _ := url.Parse("/home/hamed/tmusic/")
			u, _ := url.Parse(inputMainMusic)
			u.Path = path.Join(u.Path, fh.Filename)
			//fmt.Println("main (537):::", fh.Size)
			//fmt.Println("main (537):::", u.Path)
			s := u.String()
			nf, err := os.Create(s)
			if err != nil {
				fmt.Println(err)
			}
			//mf.Seek(0, 0)
			_, err = io.Copy(nf, mf)
			if err != nil {
				fmt.Println(err)
			}
			return fh.Filename
		}
		return cons.R_EXTENTION_NOT_ALLOWED
	}
	return cons.R_SIZE_IS_BIG
}

//-- Split audio to parts
func SplitAudio(outPutPath, musicName, inputPath string, separateTime int) (int, string) {
	removeSpace := strings.ReplaceAll(musicName, " ", "")
	finalMusicName := strings.ReplaceAll(removeSpace, ".mp3", fmt.Sprint(RandomNumber(1, 500000))+"%03d.mp3")
	folderName := strings.ReplaceAll(strings.ReplaceAll(finalMusicName, ".mp3", ""), "%", "")

	fmt.Println(finalMusicName)
	mkFolderCommand := ` mkdir ` + outPutPath + `/` + folderName
	err, _, _ := ShellOut(mkFolderCommand)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}
	command := ` ffmpeg -i ` + inputPath + musicName + ` -f segment -segment_time ` + strconv.Itoa(separateTime) +
		` -c copy ` + outPutPath + `/` + folderName + `/` + finalMusicName
	fmt.Println(command)
	err, _, _ = ShellOut(command)
	if err != nil {
		return 0, ""
	}
	return 1, outPutPath + `/` + folderName + `/`
}

//-- Get Spleeter of python library for split audios to vocal, base , etc... proposes
func Spleeter(inputPath, outputPath, stemsKind, timeDuration, offSet string) {
	command := `spleeter separate -i ` + inputPath + ` -p spleeter:` + stemsKind +
		`-o /home/hamed/` + outputPath /*+ ` -s ` + offSet + ` -d ` + timeDuration*/
	err, _, _ := ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}
}

//-- Attach portion audios to one
func AttachAudio(AttachInputPath string) {

	//command := `ls `+AttachInputPath
	//err, a, _ := ShellOut(command)
	//if err != nil {
	//}
	//splitSeparatedMusicsToArray := strings.Split(a, "\n")
	//var concatMusics
	//for _,v:=range splitSeparatedMusicsToArray {
	//
	//}
	//
	//command := `ffmpeg -i "concat:file1.mp3|file2.mp3" -acodec copy output.mp3`
	//err, _, _ := ShellOut(command)
	//if err != nil {
	//	fmt.Println(err)
	//}

}

//-- The target of this function is that find out duration time of audio
func SpecifyAudioTimeDuration(audioPath string) int {

	command := `ffmpeg -i ` + audioPath + ` 2>&1 | grep Duration`
	err, a, _ := ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}
	splitDurationString := strings.Split(strings.Split(a, ",")[0], ":")
	splitDuration := strings.Split(splitDurationString[3], ".")

	minuteToInt, err := strconv.Atoi(splitDurationString[2])
	if err != nil {
		fmt.Println(err)
	}
	secondToInt, err := strconv.Atoi(splitDuration[1])
	if err != nil {
		fmt.Println(err)
	}
	return minuteToInt*60 + secondToInt
}

//-- Give List of files (in this case musics parts which did split)
func ListOfFiles(path string) []string {
	command := `ls ` + path
	err, echo, _ := ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(echo, "\n")
}


// pull out file format like mp3, wav, wma and so on, using ffmpeg (linux software)
func FindAudioFormat(filePath string) string {
	command := `ffmpeg -i ` + filePath + ` 2>&1 | grep Input`
	err, echo, _ := ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(echo, ",")[1]
}

//func ChangeFileName(path,file string) string {
//	command := `mv `+path+file
//	err, echo, _ := ShellOut(command)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return
//
//}

//-- Get random number from n to n
func RandomNumber(min int, max int) int {
	if min == max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//-- Get find element from your slice
func HasElement(array []string, elm string) bool {
	for _, v := range array {
		if v == elm {
			return true
		}
	}
	return false
}
