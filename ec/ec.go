package ec

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//LEVEL
const (
	DEBUG = iota
	ERROR
	WARN
)

const (
	debugLog = "debug.log"
	warnLog  = "warn.log"
	errorLog = "error.log"
	logName  = "msg.log"
	logDir   = "log"
	timeFMT  = "2006-01-02 15:04:05"
	_debug   = "[-debug-]"
	_info    = "[-warn-]"
	_error   = "[-error-]"
	dirSep   = "/"
	space    = " "
	lb       = "\n"
)

type info struct {
	file   *os.File
	status bool
	*sync.RWMutex
}

var eci *info

func init() {
	dir, _ := os.Getwd()
	_, err := os.Stat(dir + dirSep + logDir)
	if nil != err {
		err = os.Mkdir(dir+dirSep+logDir, 0666)
		if nil != err {
			return
		}
	}

	fullp := dir + dirSep + logDir + dirSep + logName
	if _, err := os.Stat(fullp); nil == err {
		os.Rename(fullp, fullp+".old")
	}

	f, err := os.OpenFile(fullp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if nil != err {
		fmt.Println(err)
	}

	eci = &info{
		file:    f,
		status:  true,
		RWMutex: &sync.RWMutex{},
	}
}

// func Switch(level string, outto string) {
// 	switch level {
// 	case "debug":

// 	}
// }

//Debug1 ...用来打印输出debug信息
// func Debug1(msg interface{}) {
// 	eci.Lock()
// 	defer eci.Unlock()
// 	pc, file, line, _ := runtime.Caller(1)
// 	if eci.status {
// 		fmt.Fprintln(eci.file, debug, time.Now().Format(timeFMT),
// 			file,
// 			runtime.FuncForPC(pc).Name(),
// 			line,
// 			msg)
// 	} else {
// 		fmt.Fprintln(os.Stderr, debug, time.Now().Format(timeFMT),
// 			file,
// 			runtime.FuncForPC(pc).Name(),
// 			line,
// 			msg)
// 	}
// }

// BenchmarkDebug-6          307542              3754 ns/op
// BenchmarkDebug1-6         315781              3855 ns/op
// BenchmarkDebug-6          292567              3836 ns/op
// BenchmarkDebug1-6         300013              4005 ns/op
// BenchmarkDebug-6          272625              3943 ns/op
// BenchmarkDebug1-6         315565              3922 ns/op

//Debug 用来打印输出debug信息
func Debug(fmtstr string, v ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)

	eci.Lock()
	defer eci.Unlock()

	if eci.status {
		io.WriteString(eci.file, _debug+space+
			time.Now().Format(timeFMT)+space+
			path.Base(file)+space+
			runtime.FuncForPC(pc).Name()+space+
			strconv.Itoa(line)+space+fmt.Sprintf(fmtstr, v...)+lb)
		// io.WriteString(eci.file, )

	} else {
		io.WriteString(os.Stdout, _debug+space+
			time.Now().Format(timeFMT)+space+
			path.Base(file)+space+
			runtime.FuncForPC(pc).Name()+space+
			strconv.Itoa(line)+space+fmt.Sprintf(fmtstr, v...)+lb)
	}
}

//Errorf ...
func Errorf(fmtStr string, v ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(fmtStr, v...)
	eci.Lock()
	defer eci.Unlock()
	if eci.status {
		io.WriteString(eci.file,
			fmt.Sprintf("%s %s %s %s %d %s\n",
				_error,
				time.Now().Format(timeFMT),
				path.Base(file),
				runtime.FuncForPC(pc).Name(),
				line,
				content))
	} else {
		io.WriteString(os.Stderr,
			fmt.Sprintf("%s %s %s %s %d %s\n",
				_error,
				time.Now().Format(timeFMT),
				path.Base(file),
				runtime.FuncForPC(pc).Name(),
				line,
				content))
	}

}

//Warnf ...
func Warnf(fmtStr string, v ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(fmtStr, v...)
	eci.Lock()
	defer eci.Unlock()
	if eci.status {
		io.WriteString(eci.file,
			fmt.Sprintf("%s %s %s %s %d %s\n",
				_info,
				time.Now().Format(timeFMT),
				path.Base(file),
				runtime.FuncForPC(pc).Name(),
				line,
				content))
	} else {
		io.WriteString(os.Stdout,
			fmt.Sprintf("%s %s %s %s %d %s\n",
				_info,
				time.Now().Format(timeFMT),
				path.Base(file),
				runtime.FuncForPC(pc).Name(),
				line,
				content))
	}
}

//PostLogInfo ...
func PostLogInfo(c *gin.Context) {

}

func GetLogInfo(c *gin.Context) {

}

//Info 出错信息t
// func Info(msg string) {
// 	_, file, line, _ := runtime.Caller(1)

// 	eci.Lock()
// 	defer eci.Unlock()
// 	if eci.status {
// 		io.WriteString(eci.file, _info+time.Now().Format(timeFMT)+space+
// 			path.Base(file)+space+
// 			"line"+strconv.Itoa(line)+space+
// 			msg+lb)
// 	} else {
// 		io.WriteString(os.Stderr, _info+msg)
// 	}
// }

//Close 关掉ec
func Close() {
	eci.file.Close()
	eci = nil
}

//GetECStatus 返回ec是否可用
func GetECStatus() bool {
	return eci.status
}

//ResetECStatus ...
func ResetECStatus() {
	eci.status = false
}

//SetECStatus ...
func SetECStatus() {
	eci.status = true
}
