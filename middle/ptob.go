package middle

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
	"golang.org/x/image/bmp"
)

//Pdf2Bmp 将pdf转成位图，pdfPath：pdf所在路径，bmpPath：生成的bmp图像所在路径
//eg: Pdf2Bmp("./test.pdf", "./bmp")
func Pdf2Bmp(pdfPath string, bmpPath string) error {

	//打开pdf
	pdf, err := fitz.New(pdfPath)
	if nil != err {
		fmt.Println(err)
		return err
	}
	name := ExtName(pdfPath)
	for i := 0; i < pdf.NumPage(); i++ {
		//将pdf的每一页转成图像
		img, err := pdf.Image(i)
		if nil != err {
			fmt.Println(err)
			return err
		}

		//创建用来存除图像的
		tmpFile, err := os.Create(filepath.Join("", fmt.Sprintf(name+"%03d.bmp", i)))
		if nil != err {
			fmt.Println(err)
			return err
		}
		defer tmpFile.Close()
		// img = image.NewGray(img.Bounds())
		//encode bmp
		err = bmp.Encode(tmpFile, img)
		if nil != err {
			fmt.Println(err)
			return err
		}

	}

	return nil
}

//GrayImag 将原图数据转成灰度数据
func GrayImag(src image.Image) image.Image {
	bound := src.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	fmt.Println(dx, dy)
	gray := image.NewGray(bound)
	var gcolor color.Gray
	for r := 0; r < dx; r++ {
		for c := 0; c < dy; c += 4 {
			srcc := src.At(r, c)
			red, green, blue, _ := srcc.RGBA()
			gcolor = color.Gray{Y: uint8(float64(red)*0.299 + float64(green)*0.587 + float64(blue)*0.114)}
			gray.SetGray(r, c, gcolor)
			gray.SetGray(r, c+1, gcolor)
			gray.SetGray(r, c+2, gcolor)
			gray.SetGray(r, c+3, gcolor)
		}
	}

	return gray
}

//ExtName 提取文件的名字
func ExtName(p string) string {
	base := path.Base(p)
	suffix := path.Ext(base)
	return strings.TrimRight(base, suffix)
}
