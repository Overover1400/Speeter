package useful

import (
	"bytes"
	"cons"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"html/template"
	"interfaces"
	"io"
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

func UploadAudio(req *http.Request, inputMainMusic string) string {

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
			fmt.Println("upload is succesfull")
			return fh.Filename
		}
		return cons.R_EXTENTION_NOT_ALLOWED
	}
	return cons.R_SIZE_IS_BIG
}

//-- Split audio to parts
func SplitAudio(outPutPath, musicName, inputPath string, separateTime int) (string, string) {
	removeSpace := strings.ReplaceAll(musicName, " ", "")
	finalMusicName := strings.ReplaceAll(removeSpace, ".mp3", fmt.Sprint(RandomNumber(1, 500000))+"%03d.mp3")
	folderName := strings.ReplaceAll(strings.ReplaceAll(finalMusicName, ".mp3", ""), "%", "")

	//fmt.Println(finalMusicName)
	mkFolderCommand := ` mkdir ` + outPutPath + `/` + folderName
	err, _, _ := ShellOut(mkFolderCommand)
	if err != nil {
		fmt.Println(err)
		return cons.R_FAILED, ""
	}
	command := ` ffmpeg -i ` + inputPath + musicName + ` -f segment -segment_time ` + strconv.Itoa(separateTime) +
		` -c copy ` + outPutPath + `/` + folderName + `/` + finalMusicName
	//fmt.Println(command)
	err, _, _ = ShellOut(command)
	if err != nil {
		return cons.R_FAILED, ""
	}
	return cons.R_SUCCESSFUL, outPutPath + `/` + folderName + `/`
}

//-- Get Spleeter of python library for split audios to vocal, base , etc... proposes
func Spleeter(inputPath, outputPath, stemsKind, timeDuration, offSet string) error{
	command := `spleeter separate -i ` + inputPath + ` -p spleeter:` + stemsKind +
		`-o ` + outputPath /*+ ` -s ` + offSet + ` -d ` + timeDuration*/
	err, _, _ := ShellOut(command)
	//fmt.Println(cmd.Stdout, cmd.Stderr)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

//-- Attach portion audios to one
func AttachAudio() {
	//output type of music
	//outPutTypes :=[]string{"bass","drum","other","piano","vocals"}
	command := `ls ` + cons.OUT_PUT_SPLEETER_PATH
	err, a, _ := ShellOut(command)
	if err != nil {
	}

	splitSeparatedMusicsToArray := strings.Split(a, "\n")

	var typeOfOutputBass []string
	var typeOfOutputDrums []string
	var typeOfOutputOther []string
	var typeOfOutputPiano []string
	var typeOfOutputVocals []string
	for _, v := range splitSeparatedMusicsToArray {
		if v != "" {
			command2 := `ls ` + cons.OUT_PUT_SPLEETER_PATH + `/` + v
			err, outPutTypes, _ := ShellOut(command2)
			if err != nil {
			}

			splitSeparatedTypeMusicsToArray := strings.Split(outPutTypes, "\n")
			for i := 0; i < len(splitSeparatedTypeMusicsToArray); i++ {
				if splitSeparatedTypeMusicsToArray[i] != "" {

					var s = cons.OUT_PUT_SPLEETER_PATH + `/` + v + `/` + splitSeparatedTypeMusicsToArray[i]

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

	if len(typeOfOutputBass) > 0 {
		bass := strings.Join(typeOfOutputBass, " ")
		command := `sox ` + bass + ` ` + cons.FINAL_OUT_PUT + `/` + "bass.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputDrums) > 0 {
		drum := strings.Join(typeOfOutputDrums, " ")
		command := `sox ` + drum + ` ` + cons.FINAL_OUT_PUT + `/` + "drums.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputOther) > 0 {
		other := strings.Join(typeOfOutputOther, " ")
		command := `sox ` + other + ` ` + cons.FINAL_OUT_PUT + `/` + "other.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputPiano) > 0 {
		piano := strings.Join(typeOfOutputPiano, " ")
		command := `sox ` + piano + ` ` + cons.FINAL_OUT_PUT + `/` + "piano.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(typeOfOutputVocals) > 0 {
		vocals := strings.Join(typeOfOutputVocals, " ")
		command := `sox ` + vocals + ` ` + cons.FINAL_OUT_PUT + `/` + "vocals.wav"
		err, _, _ := ShellOut(command)
		if err != nil {
			fmt.Println(err)
		}
	}

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
	fmt.Println("time duration is succesfully : ",minuteToInt*60+secondToInt)

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

	dbUserName :=/*os.Getenv("hamed")*/"hamed"
	dbUserPass :=`*7yH09hamed125&^mn7!`
	dbDatabases := `spleeter`
	addressIp :=`194.5.175.118`
	port := `3306`

	os.Getenv(".env")
	dbUrl :=fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci",dbUserName,dbUserPass,
		addressIp,port,dbDatabases)
	db,err=sql.Open("mysql",dbUrl)
	if err != nil {
	fmt.Println(err)
	return
	}

	err=db.Ping()
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

func ParsHtmFiles(res http.ResponseWriter, htmlFileName string,data struct{}) {
	templates,err:=template.ParseFiles(cons.HTLM_FOLDER+ htmlFileName+".html")
	if err != nil {
		fmt.Println(err)
	}

	err=templates.ExecuteTemplate(res,htmlFileName+".html", data)
	if err != nil {
		fmt.Println(err)
	}
}