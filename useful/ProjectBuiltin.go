package useful

import (
	"bytes"
	"cons"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"interfaces"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func PrintReqForm(PageName string, req *http.Request) {
	fmt.Println("======>", PageName, "\n", req.Form)
}

func ShellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(cons.ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	//fmt.Println(cmd.Stdout, cmd.Stderr)
	return err, stdout.String(), stderr.String()
}

func UploadAudio(req *http.Request, inputMainMusic string) (string, string) {

	mf, fh, err := req.FormFile("file")
	defer mf.Close()
	if err != nil {
		fmt.Print(err)
	}

	//-- Check file size and format
	fh.Filename = strings.ReplaceAll(fh.Filename, " ", "")

	if fh.Size <= 30000000 {
		audioExtesnion := fh.Header.Get("Content-Type")

		legalExtensions := []string{"audio/mp3", "audio/aac", "audio/wma", "audio/flac", "audio/wav", "audio/aiff"}

		//TODO 2
		//extention :=FindAudioFormat(inputMainMusic+fh.Filename)

		if HasElement(legalExtensions, audioExtesnion) {

			err = os.MkdirAll(inputMainMusic, os.ModePerm)

			if err != nil {
				fmt.Println(err)
			}
			nf, err := os.Create(inputMainMusic + fh.Filename)
			if err != nil {
				fmt.Println(err, "--", inputMainMusic+fh.Filename)
			}
			//mf.Seek(0, 0)
			_, err = io.Copy(nf, mf)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("upload is succesfull")
			return fh.Filename, inputMainMusic + fh.Filename
		}
		return cons.R_EXTENTION_NOT_ALLOWED, ""
	}
	return cons.R_SIZE_IS_BIG, ""
}

//-- Split audio to parts
func SplitAudio(outPutPath, musicName, inputPath string, separateTime int) (string, string) {
	removeSpace := strings.ReplaceAll(musicName, " ", "")
	finalMusicName := strings.ReplaceAll(removeSpace, ".mp3", fmt.Sprint(RandomNumber(1, 500000))+"%03d.mp3")
	folderName := strings.ReplaceAll(strings.ReplaceAll(finalMusicName, ".mp3", ""), "%", "")

	//fmt.Println(finalMusicName)
	//mkFolderCommand := ` mkdir ` + outPutPath+ folderName

	err := os.MkdirAll(outPutPath+folderName, os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
	//err, r, n := ShellOut(mkFolderCommand)
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println(r)
	//	fmt.Println(n)
	//	return cons.R_FAILED, ""
	//}
	command := ` ffmpeg -i ` + inputPath + musicName + ` -f segment -segment_time ` + strconv.Itoa(separateTime) +
		` -c copy ` + outPutPath+ folderName + `/` + finalMusicName
	fmt.Println(command)
	err, _, _ = ShellOut(command)
	if err != nil {
		return cons.R_FAILED, ""
	}
	return cons.R_SUCCESSFUL, outPutPath + folderName + `/`
}

//-- Get Spleeter of python library for split audios to vocal, base , etc... proposes
func Spleeter(inputPath, outputPath, stemsKind, timeDuration, offSet string) error {
	//command :=`conda activate spleeter`
	//
	//err,_,_:=ShellOut(command)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//err := os.Setenv("CONDA_DEFAULT_ENV", "spleeter")
	err := os.Setenv("CONDA_DEFAULT_ENV", "base")

	if err != nil {
		fmt.Println(err)
	}

	var timedurationAndOffset string
	if timeDuration != ""{
		timedurationAndOffset=` -s `+offSet+` -d `+timeDuration
	}

	command2 := `spleeter separate -i ` + inputPath + ` -p spleeter:` + stemsKind +
		`-o ` + outputPath +timedurationAndOffset
	err, r, n := ShellOut(command2)
	fmt.Println(command2)

	if err != nil {
		fmt.Println(err)
		fmt.Println(r)
		fmt.Println(n)
	}
	return err
}

//-- Attach portion audios to one
func AttachAudio(folderName string) {
	//output type of music
	//outPutTypes :=[]string{"bass","drum","other","piano","vocals"}
	command := `ls ` + cons.OUT_PUT_SPLEETER_PATH + folderName
	err, a, _ := ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}

	splitSeparatedMusicsToArray := strings.Split(a, "\n")
	fmt.Println("111.222", splitSeparatedMusicsToArray)
	var typeOfOutputBass []string
	var typeOfOutputDrums []string
	var typeOfOutputOther []string
	var typeOfOutputPiano []string
	var typeOfOutputVocals []string
	for _, v := range splitSeparatedMusicsToArray {
		if v != "" {
			command2 := `ls ` + cons.OUT_PUT_SPLEETER_PATH + `/` + folderName + `/` + v
			err, outPutTypes, _ := ShellOut(command2)
			if err != nil {
			}

			splitSeparatedTypeMusicsToArray := strings.Split(outPutTypes, "\n")
			for i := 0; i < len(splitSeparatedTypeMusicsToArray); i++ {
				if splitSeparatedTypeMusicsToArray[i] != "" {

					var s = cons.OUT_PUT_SPLEETER_PATH + folderName + `/` + v + `/` + splitSeparatedTypeMusicsToArray[i]

					switch splitSeparatedTypeMusicsToArray[i] {
					case "bass.wav":
						typeOfOutputBass = append(typeOfOutputBass, s)
					case "drums.wav":
						typeOfOutputDrums = append(typeOfOutputDrums, s)
					case "other.wav":
						typeOfOutputOther = append(typeOfOutputOther, s)
					case "piano.wav":
						typeOfOutputPiano = append(typeOfOutputPiano, s)
					case "vocals.wav":
						typeOfOutputVocals = append(typeOfOutputVocals, s)
					}
				}
			}
		}
	}
	err = os.MkdirAll(cons.FINAL_OUT_PUT+folderName, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("folder Name : ",cons.FINAL_OUT_PUT+folderName, os.ModePerm)

	if len(typeOfOutputBass) > 0 {
		bass := strings.Join(typeOfOutputBass, " ")
		command := `sox ` + bass + ` ` + cons.FINAL_OUT_PUT + folderName + `/` + "bass.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputDrums) > 0 {
		drum := strings.Join(typeOfOutputDrums, " ")
		command := `sox ` + drum + ` ` + cons.FINAL_OUT_PUT + folderName + `/` + "drums.wav"
		fmt.Println("---", folderName)
		fmt.Println("---", command)
		err, r, n := ShellOut(command)
		if err != nil {
			fmt.Println("===", err)
			fmt.Println("===", r)
			fmt.Println("===", n)
		}
	}

	if len(typeOfOutputOther) > 0 {
		other := strings.Join(typeOfOutputOther, " ")
		command := `sox ` + other + ` ` + cons.FINAL_OUT_PUT + folderName + `/` + "other.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputPiano) > 0 {
		piano := strings.Join(typeOfOutputPiano, " ")
		command := `sox ` + piano + ` ` + cons.FINAL_OUT_PUT + folderName + `/` + "piano.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputVocals) > 0 {
		vocals := strings.Join(typeOfOutputVocals, " ")
		command := `sox ` + vocals + ` ` + cons.FINAL_OUT_PUT + folderName + `/` + "vocals.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("222", err)
}

//-- This function responsible for find out duration of audio time (second)
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
	secondToInt, err := strconv.Atoi(splitDuration[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("time duration is succesfully : ", minuteToInt*60+secondToInt)

	return minuteToInt*60 + secondToInt
}

//-- Give List of files (in this case musics parts which did split)
func ListOfFiles(path string) (array []string) {
	command := `ls ` + path + `/`
	err, echo, _ := ShellOut(command)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range strings.Split(echo, "\n") {
		if v != "" {
			array = append(array, v)
		}
	}
	return array
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

func OpenDatabases() (db *sql.DB, err error) {

	dbUserName := /*os.Getenv("hamed")*/ "hamed"
	dbUserPass := `*7yH09&^mn7!`
	dbDatabases := `spleeter`
	addressIp := `194.5.195.203`
	port := `3306`

	os.Getenv(".env")
	dbUrl := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci", dbUserName, dbUserPass,
		addressIp, port, dbDatabases)
	db, err = sql.Open("mysql", dbUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func HandleCloseableErrorClient(value interfaces.CloseAllDb, fileName string, lineNumber int) {
	//ClientPrintErr("BuiltinHelpersOrigin.go", 1199, fileName," ",lineNumber," ",value)
	if value == nil {
		return
	}
	err := value.Close()
	if err != nil {
		fmt.Println()
		//ClientPrintErr(fileName, lineNumber, err)
	}
}

func ParsHtmFiles(res http.ResponseWriter, htmlFileName string, data interface{}) {
	templates, err := template.ParseFiles(cons.HTLM_FOLDER + htmlFileName + ".html")
	if err != nil {
		fmt.Println(err)
	}
	err = templates.ExecuteTemplate(res, htmlFileName+".html", data)
	if err != nil {
		fmt.Println(err)
	}
}

func RemoveExtension(str string) string {
	var sv = true
	var j = 1
	for i := 1; i < len(str); i++ {

		if string(str[len(str)-i]) != `.` && sv {
			j++
		}
		if string(str[len(str)-i]) == `.` {
			sv = false
		}
	}
	return str[:len(str)-j]
}

func Neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/") || req.URL.Path == "" {
			http.Redirect(res, req, "http://http://194.5.195.203:4002/main", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(res, req)
	})
}

func RandomString(max int) string {
	var randomString string
	var alphabet []int
	for i := 65; i < 120; i++ {
		if i < 91 || i >= 97 {
			alphabet = append(alphabet, i)
		}
	}
	for i := 0; i <= max; i++ {
		randomString += string(alphabet[RandomNumber(0, len(alphabet))])
	}
	return randomString
}

func SaveFile(req *http.Request, folderName string) {

	mf, fh, err := req.FormFile("file")
	defer mf.Close()
	if err != nil {
		fmt.Print(err)
	}

	//-- Check file size and format
	fh.Filename = strings.ReplaceAll(fh.Filename, " ", "")

	if fh.Size <= 30000000 {
		audioExtesnion := fh.Header.Get("Content-Type")

		legalExtensions := []string{"audio/mp3", "audio/aac", "audio/wma", "audio/flac", "audio/wav", "audio/aiff"}

		//TODO 2
		//extention :=FindAudioFormat(inputMainMusic+fh.Filename)

		waitingLine := cons.WAITING_LINE_FOLDER + `/` + folderName + `/`

		if HasElement(legalExtensions, audioExtesnion) {
			err = os.MkdirAll(waitingLine, os.ModePerm)

			if err != nil {
				fmt.Println(err)
			}
			nf, err := os.Create(waitingLine + fh.Filename)
			if err != nil {
				fmt.Println(err, "--", waitingLine+fh.Filename)
			}
			//mf.Seek(0, 0)
			_, err = io.Copy(nf, mf)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("upload is succesfull")
		}
	}
}
