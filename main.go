package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	dd := "123 34 100 101 118 95 116 121 112 101 34 58 34 240 16 34 44 34 101 118 101 110 116 95 116 121 112 101 34 58 34 128 17 34 44 34 105 100 34 58 48 125"
	ds := strings.Split(dd, " ")
	var tmp []byte
	for _, c := range ds {
		c0, _ := strconv.Atoi(c)
		fmt.Printf("%c", c0)
		tmp = append(tmp, byte(c0))
	}
	fmt.Println(string(tmp))
	// go tt(nil)
	// time.Sleep(time.Second)
	// fmt.Println("main exit")
}

// func tt(v interface{}) {
// 	defer fmt.Println("exit")
// 	kill.IfNil(v, "test exist")

// 	fmt.Println("am I run?")
// }
