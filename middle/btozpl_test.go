package middle

import "testing"

func TestBmp2Zpl(t *testing.T) {
	err := Bmp2Zpl("./test1000.bmp", ".")
	if nil != err {
		t.Error(err)
	}

	t.Error("done")
}
