package customlogger

import (
	"fmt"
	"go-mysql/config"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Customlog struct {
	filename string
	*log.Logger
}

var logger *Customlog
var once sync.Once

// start log
func GetInstance() *Customlog {
	// once.Do(func() {
	//dt := time.Now()
	var logFile string

	if os.Getenv("log-files") == ""{
		logFile = GetPathLog()
	}else{
		err := PrintSize(os.Getenv("log-files"))
		if err != nil {
			logFile = GetPathLog()
			//panic(err)
		}else {
			logFile = os.Getenv("log-files")
		}
	}

	fname := logFile
	logger = createLogger(fname)
	// })
	return logger
}

func PrintSize(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	x := info.Size()
	if x > 10485760{
		return fmt.Errorf(`File Sudah >= 10 MB.`)
	}
	return nil
}

func createLogger(fname string) *Customlog {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return &Customlog{
		filename: fname,
		Logger:   log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func GetPathLog() string {
	var result string =""
	dt := time.Now()

	dir := config.GetValueFromConfig("", "log.filepathApi")
	if dir == "" || dir == "<nil>"{
		tempDir, errtempDir := filepath.Abs(filepath.Dir(os.Args[0]))
		if errtempDir != nil {
			log.Fatal(errtempDir)
		}
		tempDir = fmt.Sprintf(`%s/logs`,tempDir)
		dir = tempDir
	}
	dir = checkPathFolder(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println(err.Error())
	}else {
		result = config.GetValueFromConfig("", "log.fileApi")
		result = dir + result + "_" + fmt.Sprintf("%s_%s", dt.Format("2006-01-02"), dt.Format("15-04-05")) + ".log"
		os.Setenv("log-files", result)
		return result
	}
	return result
}

func checkPathFolder(path string) string {
	var err error
	err = nil
	osPC := runtime.GOOS
	var checkPath string = ""
	if osPC == "windows" {
		path = strings.Replace(path, `/`, `\`, -1)
		s := strings.Split(path, `\`)
		for i := 0; i < len(s); i ++  {
			if checkPath == ""{
				checkPath = s[i]
				if _, err = os.Stat(checkPath); os.IsNotExist(err) {
					os.Mkdir(checkPath, os.ModePerm)
				}
			}else {
				checkPath = fmt.Sprintf(`%s\%s`, checkPath, s[i])
				if _, err = os.Stat(checkPath); os.IsNotExist(err) {
					os.Mkdir(checkPath, os.ModePerm)
				}
			}
		}
		checkPath = fmt.Sprintf(`%s\`, checkPath)
		checkPath = strings.Replace(checkPath, `\\`, `\`, -1)
	} else if osPC == "linux" {
		s := strings.Split(path, `/`)
		for i := 0; i < len(s); i ++  {
			if checkPath == ""{
				checkPath = fmt.Sprintf(`/%s`, s[i])
				if _, err := os.Stat(checkPath); os.IsNotExist(err) {
					os.Mkdir(checkPath, os.ModePerm)
				}
			}else {
				checkPath = fmt.Sprintf(`%s/%s`, checkPath, s[i])
				if _, err := os.Stat(checkPath); os.IsNotExist(err) {
					os.Mkdir(checkPath, os.ModePerm)
				}
			}
		}
		checkPath = fmt.Sprintf(`%s/`, checkPath)
		checkPath = strings.Replace(checkPath, `//`, `/`, -1)
	} else {
		s := strings.Split(path, `/`)
		for i := 0; i < len(s); i ++  {
			if checkPath == ""{
				checkPath = fmt.Sprintf(`/%s`, s[i])
				if _, err := os.Stat(checkPath); os.IsNotExist(err) {
					os.Mkdir(checkPath, os.ModePerm)
				}
			}else {
				checkPath = fmt.Sprintf(`%s/%s`, checkPath, s[i])
				if _, err := os.Stat(checkPath); os.IsNotExist(err) {
					os.Mkdir(checkPath, os.ModePerm)
				}
			}
		}
		checkPath = fmt.Sprintf(`%s/`, checkPath)
		checkPath = strings.Replace(checkPath, `//`, `/`, -1)
	}
	return checkPath
}