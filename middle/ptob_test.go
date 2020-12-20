package middle

import "testing"

func TestPdf2Bmp(t *testing.T) {
	p := "D:/work-doc/Ireport-test/test1.pdf"
	// p := "D:/myown/myprinter/test.pdf"
	b := "./"
	err := Pdf2Bmp(p, b)
	if nil != err {
		t.Error(err)
	}
	t.Error("test")
}
