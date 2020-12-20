package middle

import "testing"

func TestReadBmp(t *testing.T) {
	bmp, err := ReadBMP("./label000.bmp")
	if nil != err {
		t.Error(err)
	}
	t.Error(bmp)
}
