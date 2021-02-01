package safe

import (
	"strconv"
	"testing"
	"time"
)

var glob = NewGlobal()

func TestGetInt(t *testing.T) {
	st := time.Now()
	glob.Set(map[string]interface{}{"heloo": 1234})
	t.Error(time.Now().Sub(st).Nanoseconds())

	v := glob.GetMap()
	t.Error(v)

	go func() {
		for i := 0; i < 10; i++ {
			v0 := glob.GetMap()

			v0.Set(strconv.Itoa(i), i)

			// glob.Set(v)
		}
		t.Error("done")
	}()

	go func() {
		var i int
		for {
			if i >= 10 {
				i = 0
			}
			v := glob.GetMap().Get(strconv.Itoa(i))
			i++
			t.Error(v)
		}
	}()
	time.Sleep(2 * time.Second)
	t.Error(v)
}
