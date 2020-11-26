package kill

import (
	"fmt"
	"runtime"
)

type Killer interface {
	KillIfNil(interface{}, ...string)
}

func killIfNil(v interface{}, s ...string) {
	if v0, ok := v.(Killer); ok {
		v0.KillIfNil(v, s...)
		return
	}

	if nil == v {
		if 0 == len(s) {
			fmt.Println("nil data exist")
		} else {
			fmt.Println(s)
		}
		runtime.Goexit()
	}
}

// func KillIfNil(v interface{}, s ...string) {

// }

//IfNil : if v is nil, then kill the current gorountine
func IfNil(v interface{}, s ...string) {
	if v0, ok := v.(Killer); ok {
		v0.KillIfNil(v, s...)
		return
	}

	fmt.Println(runtime.GOROOT(), runtime.NumCPU(),runtime.)

	if nil == v {
		if 0 == len(s) {
			fmt.Println("nil data exist")
		} else {
			fmt.Println(s)
		}
		runtime.Goexit()
	}
}
