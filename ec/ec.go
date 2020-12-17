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
)

const (
	logName = "msg.log"
	logDir  = "log"
	timeFMT = "2006-01-02 15:04:05"
	debug   = "[debug]"
	_info   = "[info]"
	dirSep  = "/"
	space   = " "
	lb      = "\n"
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

	f, err := os.OpenFile(dir+dirSep+logDir+dirSep+logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if nil != err {

	}

	eci = &info{
		file:    f,
		status:  true,
		RWMutex: &sync.RWMutex{},
	}
}

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
		io.WriteString(eci.file, debug+space+
			time.Now().Format(timeFMT)+space+
			path.Base(file)+space+
			runtime.FuncForPC(pc).Name()+space+
			strconv.Itoa(line)+space+fmt.Sprintf(fmtstr, v...)+lb)
		// io.WriteString(eci.file, )

	} else {
		io.WriteString(os.Stderr, debug+space+
			time.Now().Format(timeFMT)+space+
			path.Base(file)+space+
			runtime.FuncForPC(pc).Name()+space+
			strconv.Itoa(line)+space+fmt.Sprintf(fmtstr, v...)+lb)
	}
}

//Info 出错信息t
func Info(msg string) {
	_, file, line, _ := runtime.Caller(1)

	eci.Lock()
	defer eci.Unlock()
	if eci.status {
		io.WriteString(eci.file, _info+time.Now().Format(timeFMT)+space+
			path.Base(file)+space+
			"line"+strconv.Itoa(line)+space+
			msg+lb)
	} else {
		io.WriteString(os.Stderr, _info+msg)
	}
}

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
