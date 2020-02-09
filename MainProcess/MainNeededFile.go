package MainProcess

import (
	"cons"
	"fmt"
	"html/template"
	"net/http"
	"useful"
)

func MainProcess(res http.ResponseWriter, req *http.Request) {

	inputMainMusic := cons.MAIN_MUSIC_PATH

	//--Upload files
	musicName := useful.UploadAudio(req, inputMainMusic)
	if musicName == cons.R_EXTENTION_NOT_ALLOWED {
		fmt.Fprintln(res, "File format not allowed")
		return
	} else if musicName == cons.R_SIZE_IS_BIG {
		fmt.Fprintln(res, "File size must smaller then 30 mg")
		return
	}

	//-- Get audio duration (second) and then divide by 10
	audioTimeDuration := useful.SpecifyAudioTimeDuration(inputMainMusic + musicName)


	outPutPathOfSpleeter := cons.OUT_PUT_SPLEETER_PATH


	//-- 40 is second here
	if audioTimeDuration <= 40 {
		err:=useful.Spleeter(inputMainMusic, outPutPathOfSpleeter /*+strconv.Itoa(i)*/, "5stems ",	"", "")
		if err == nil {
			useful.AttachAudio()
		}
	} else {

		outPutPath := cons.OUT_PUT_SPLIT_MUSIC_PATH

		//-- Split music to a few part
		errCondition, splitOutPutPath := useful.SplitAudio(outPutPath, musicName, inputMainMusic, 40)
		if errCondition == cons.R_FAILED {
			fmt.Println(cons.R_FAILED)
			return
		}

		sliceOfMusicParts := useful.ListOfFiles(splitOutPutPath)

		fmt.Fprintln(res,` زمان تقریبی `,len(sliceOfMusicParts)*40,`ثانیه`)
		for _, v := range sliceOfMusicParts {
			useful.Spleeter(splitOutPutPath+v, outPutPathOfSpleeter /*+strconv.Itoa(i)*/, "5stems ",
				"", "")
		}

		useful.AttachAudio()
	}
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
		MainProcess(res, req)
	}

	fmt.Fprintln(res, "There is no error !")
}

