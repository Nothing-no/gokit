package safe

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const (
	VOID = iota
	INT
	BOOL
	FLOAT
	STRING
	ARRAY
)

type (
	//Global 只用于全局变量，考虑到全局变量可能存在竞态，所以加锁保护
	//全局变量的存在，可能导致同时写，或者一写一读的情况
	Global struct {
		v  interface{}
		rw *sync.RWMutex
	}
	//Map 自定义map[string]interface{}类型，使其可以拥有方法
	Map map[string]interface{}
	//AVLData ..
	// AVLData struct{

	// }

	//AVLNode ...
	// AVLNode struct {
	// 	Data       interface{} //数据域
	// 	hight      int32       //
	// 	mother     *AVLNode
	// 	leftChild  *AVLNode
	// 	rightChild *AVLNode
	// }
	decodeElmt struct {
		element byte
		count   uint32
		total   uint32
	}
	decodeModel struct {
		flag   FlagType
		length uint32
		data   []decodeElmt
		proto  [3]byte
	}
)

var rwMuetx sync.RWMutex

//NewGlobal 初始化一个新全局变量
func NewGlobal(v ...interface{}) *Global {
	if len(v) == 0 {
		v = append(v, nil)
	}
	return &Global{
		v:  v[0],
		rw: &sync.RWMutex{},
	}
}

func (my Global) String() string {
	return fmt.Sprintf("%v", my.v)
}

//Get get global value
func (my *Global) Get() interface{} {
	if my.rw == nil {
		my.rw = &sync.RWMutex{}
	}
	my.rw.RLock()
	defer my.rw.RUnlock()
	return my.v
}

//Set set global value
func (my *Global) Set(v interface{}) {
	if my.rw == nil {
		my.rw = &sync.RWMutex{}
	}
	my.rw.Lock()
	defer my.rw.Unlock()
	my.v = v
}

//GetInt get int value
func (my *Global) GetInt() (int, error) {
	if my.rw == nil {
		my.rw = &sync.RWMutex{}
	}
	my.rw.RLock()
	defer my.rw.RUnlock()

	if v, ok := my.v.(int); ok {
		return v, nil
	}

	return 0, errors.New("get int not ok")
}

//GetInt64 get int64 value
func (my *Global) GetInt64() (int64, error) {
	if my.rw == nil {
		my.rw = &sync.RWMutex{}
	}
	my.rw.RLock()
	defer my.rw.RUnlock()

	if v, ok := my.v.(int64); ok {
		return v, nil
	}

	return 0, errors.New("get int64 not ok")
}

//GetFloat get float64 value
func (my *Global) GetFloat() (float64, error) {
	if my.rw == nil {
		my.rw = &sync.RWMutex{}
	}
	my.rw.RLock()
	defer my.rw.RUnlock()
	if v, ok := my.v.(float64); ok {
		return v, nil
	}

	return 0, errors.New("get float64 not ok")
}

//GetBool get bool value
func (my *Global) GetBool() bool {
	my.rw.RLock()
	defer my.rw.RUnlock()
	if v, ok := my.v.(bool); ok {
		return v
	}

	return false
}

//GetMap get map value
func (my *Global) GetMap() *Map {
	if my.rw == nil {
		my.rw = &sync.RWMutex{}
		my.v = make(map[string]interface{})
	}

	my.rw.RLock()
	defer my.rw.RUnlock()
	if v, ok := my.v.(map[string]interface{}); ok {
		r := make(Map)
		r = v
		return &r
	}

	return nil
}

//Get get map value(map)
func (my *Map) Get(key string) interface{} {
	rwMuetx.RLock()
	defer rwMuetx.RUnlock()
	if my == nil {
		return nil
	}
	return (*my)[key]
}

//Set get map value (map[string]interface{})
func (my *Map) Set(key string, v interface{}) {
	if my == nil {
		my = &Map{}
		fmt.Println("nil")
	}
	rwMuetx.Lock()
	defer rwMuetx.Unlock()

	(*my)[key] = v
}

//FlagType ...
type FlagType uint8

//Code ..
type Code struct {
	Data        []byte
	Consumption time.Duration
	Src         *[]byte
}

//Encode32 ...
func Encode32(src *[]byte, flag FlagType) (*Code, error) {

	var (
		ivTmp   []byte
		ivCount = make(map[byte]uint32)
		ivStat  = make(map[byte]uint32)
		ivFlag  byte
		ivLen   [4]byte
		ivDest  = &Code{
			Src: src,
		}
		// ivSrcLen uint32
	)

	//ivFlag 编码后的第一个字节，暂时没用
	ivFlag |= byte(flag)

	//粗略统计编码时长
	ivStartTime := time.Now()

	//源数据长度
	ivSrcLen := len(*src)
	// fmt.Println(ivSrcLen)
	ivLen = i32ToByte(uint32(ivSrcLen))
	// fmt.Println(ivLen)

	//统计
	for i, b := range *src {
		ivCount[b] += uint32(i)
		ivStat[b]++
	}

	//组装flag和原始数据长度
	ivTmp = append(ivTmp, ivFlag)
	ivTmp = append(ivTmp, ivLen[:4]...)
	//组装数据
	for k, v := range ivCount {
		ivTmp = append(ivTmp, k)
		// fmt.Println(v)
		itCount := i32ToByte(v)
		ivTmp = append(ivTmp, itCount[:4]...)
		itStat := i32ToByte(ivStat[k])
		ivTmp = append(ivTmp, itStat[:4]...)
		// fmt.Println(string(k), itCount, itStat)
	}
	//组装标识
	ivTmp = append(ivTmp, []byte("nte")...)
	//获取耗时
	ivDest.Consumption = time.Now().Sub(ivStartTime)
	ivDest.Data = ivTmp

	return ivDest, nil
}

func i32ToByte(v uint32) (dest [4]byte) {
	for i := 0; i < 4; i++ {
		dest[3-i] = byte(v >> (i * 8))
	}

	return
}

func Decode32(src *[]byte) *Code {
	var (
		ivDest = &Code{
			Src: src,
		}
	)

	ivSTT := time.Now()
	ivDM := getDM(src)
	// tmp := (*decodeModel)(unsafe.Pointer(&(*src)[0]))
	// fmt.Println(tmp.data, tmp.length, tmp.proto)
	fmt.Println(ivDM)
	// itBuf := make([]byte, int(ivDM.length))

	ivDest.Consumption = time.Now().Sub(ivSTT)

	return ivDest
}

func getDM(src *[]byte) *decodeModel {
	var ivDM decodeModel

	ivDM.flag = FlagType((*src)[0])
	ivDM.length = byteToi32((*src)[1:5])
	itLen := len(*src)
	for i := 0; i < 3; i++ {
		ivDM.proto[2-i] = (*src)[itLen-i-1]
	}

	if ivDM.flag == FlagType(0) {
		if ivDM.proto != [3]byte{110, 116, 101} {
			fmt.Println("proto error")
			return nil
		}
	}

	for i := 5; i < itLen-3; i++ {
		var tmparr []byte
		tmparr = append(tmparr, (*src)[i:i+9]...)
		e := (*decodeElmt)(unsafe.Pointer(&tmparr[0]))
		i += 8
		ivDM.data = append(ivDM.data, *e)
	}
	// tmp := (*src)[5 : itLen-3]
	// tmpElmt := (*[1024]decodeElmt)(unsafe.Pointer(&tmp[0]))[:]
	// fmt.Println(tmpElmt)
	return &ivDM
}

func byteToi32(v []byte) (dest uint32) {
	for i, v0 := range v {
		dest |= (uint32(v0) << (3 - i))
	}

	return
}
