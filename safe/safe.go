package safe

import (
	"errors"
	"fmt"
	"sync"
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
	//Global 全局变量，考虑到全局变量可能存在竞态，所以加锁保护
	Global struct {
		v  interface{}
		rw *sync.RWMutex
	}
	//Map map[string]interface{}
	Map map[string]interface{}
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
	my.rw.RLock()
	defer my.rw.RUnlock()
	return my.v
}

//Set set global value
func (my *Global) Set(v interface{}) {
	my.rw.Lock()
	defer my.rw.Unlock()
	my.v = v
}

//GetInt get int value
func (my *Global) GetInt() (int, error) {
	my.rw.RLock()
	defer my.rw.RUnlock()
	if v, ok := my.v.(int); ok {
		return v, nil
	}

	return 0, errors.New("get int not ok")
}

//GetInt64 get int64 value
func (my *Global) GetInt64() (int64, error) {
	my.rw.RLock()
	defer my.rw.RUnlock()
	if v, ok := my.v.(int64); ok {
		return v, nil
	}

	return 0, errors.New("get int64 not ok")
}

//GetFloat get float64 value
func (my *Global) GetFloat() (float64, error) {
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
	rwMuetx.Lock()
	defer rwMuetx.Unlock()
	(*my)[key] = v
}
