package safe

import (
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

//new type
type (
	Int struct {
		v int64
		l *sync.RWMutex
	}

	Float struct {
		v float64
		l *sync.RWMutex
	}

	String struct {
		v string
		l *sync.RWMutex
	}

	Bool struct {
		v bool
		l *sync.RWMutex
	}

	Void struct {
		v interface{}
		l *sync.RWMutex
	}

	Array struct {
		v []interface{}
		l *sync.RWMutex
	}

	GlobalManager struct{}
)

type Getter interface {
	Get() Void
}

type Setter interface {
	Set(Void)
}

func NewInt() *Int {
	return &Int{
		l: &sync.RWMutex{},
	}
}

func NewFloat() *Float {
	return &Float{
		l: &sync.RWMutex{},
	}
}

func NewBool() *Bool {
	return &Bool{
		l: &sync.RWMutex{},
	}
}

func NewString() *String {
	return &String{
		l: &sync.RWMutex{},
	}
}

func NewArray() *Array {
	return &Array{
		l: &sync.RWMutex{},
	}
}

func (my *Int) Get() int64 {
	my.l.RLock()
	defer my.l.RUnlock()
	return my.v
}

func (my *Bool) Get() bool {
	my.l.RLock()
	defer my.l.RUnlock()
	return my.v
}

func (my *String) Get() string {
	my.l.RLock()
	defer my.l.RUnlock()
	return my.v
}

func (my *Float) Get() float64 {
	my.l.RLock()
	defer my.l.RUnlock()
	return my.v
}

func (my *Array) Get() []interface{} {
	my.l.RLock()
	defer my.l.RUnlock()
	return my.v
}

func (my *Int) Set(v int64) {
	my.l.Lock()
	defer my.l.Unlock()
	my.v = v
}

func (my *Float) Set(v float64) {
	my.l.Lock()
	defer my.l.Unlock()
	my.v = v
}

func (my *Bool) Set(v bool) {
	my.l.Lock()
	defer my.l.Unlock()
	my.v = v
}

func (my *String) Set(v string) {
	my.l.Lock()

}
