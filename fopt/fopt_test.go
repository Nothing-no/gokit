package fopt

import "testing"

func TestZip(t *testing.T) {
	fp := ZipPrepare{ZipName: "test.zip", SrcFiles: []string{"./testzip", "../xls"}}
	err := fp.Zip()
	if nil != err {
		t.Error(err)
	}
}

func TestUnzip(t *testing.T) {
	fp := UnzipPrepare{SrcPath: "./test.zip", OutPath: "./test1"}
	err := fp.Unzip()
	if nil != err {
		t.Error(err)
	}
}

func TestGetFiles(t *testing.T) {
	r, err := getDirFiles("../xls")
	if nil != err {
		t.Error(err)
	}
	rs := dealName("../../test/../del")
	t.Error(r, rs)

}
