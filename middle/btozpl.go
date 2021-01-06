package middle

import (
	"encoding/base64"
	"fmt"
	"os"

	LZ77 "github.com/fbonhomm/LZ77/source"
)

const (
	startStr = `^XA
	^SZ2^JMA
	^MCY^PMN
	^PW822
	~JSN
	^JZY
	^LH0,0^LRN
	^XZ
	`
	bmpStr     = `~DGR:SSGFX000.GRF,%d,%d,:Z64:%s`
	startPrint = `^XA
	^FO15,208
	^XGR:SSGFX000.GRF,1,1^FS
	^PQ1,0,1,Y
	^XZ
	`
	endStr = `^XA
	^IDR:SSGFX000.GRF^XZ`

	coderChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

/*
x:width mm
y:height mm
z:dot per mm dpmm 12
8 = bits/byte

t=(x*z/8)*y*z
w=x*z/8
*/

//Bmp2Zpl ...
func Bmp2Zpl(p string, zplPath string) error {
	bmp, err := ReadBMP(p)
	if nil != err {
		fmt.Println("read bmp failed")
		return err
	}
	coder := base64.NewEncoding(coderChar)
	lz77 := LZ77.Init()
	comdata := lz77.Compression(bmp.Data)
	b64 := coder.EncodeToString(comdata)
	t := 37440
	w := 98
	bmpdone := fmt.Sprintf(bmpStr, t, w, b64)
	name := ExtName(p)
	f, err := os.Create(zplPath + "/" + name + ".zpl")
	if nil != err {
		return err
	}
	defer f.Close()
	n, err := f.Write([]byte(startStr + bmpdone + startPrint + endStr))
	if nil != err {
		fmt.Println(n, err)
		return err
	}
	// binary.Write(f, binary.LittleEndian, &bmp.Data)
	fmt.Println("total:", n)
	// f.Write([]byte(startPrint + endStr))
	return nil
}
