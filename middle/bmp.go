package middle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

//BMPHead 14B
type BMPHead struct {
	Type   [2]byte //2B type:BM(default)
	Size   uint32  //4B bmp file size
	R1     uint16  //2B reserve
	R2     uint16  //2B reserve
	Offset uint32  //4B start address
}

//BMPInfo 40B
type BMPInfo struct {
	Size            uint32 //4B struct size
	Width           int32  //4B horizontal width of bitmap in pixels
	Height          int32  //4B vertical height of bitmap in pixels
	Planes          uint16 //2B number of planes
	BitPerPixel     uint16 //2B 1,4,8,16,24----
	Compression     uint32 //4B 0 no compression, 1 8bit RLE encoding, 2 4bit RLE encoding
	ImageSize       uint32 //4B siez of image if Compression =0 it is valid to set 0
	XpixelsPerM     uint32 //4B horizontal resolution:pixels/meter
	YpixelsPerM     uint32 //4B vertical resolution
	ColorUsed       uint32 //4B number of actually used colors
	ImportantColors uint32 //4B
}

//BMPColorTable 4B
type BMPColorTable struct {
	Red      uint8
	Green    uint8
	Blue     uint8
	Reserved uint8
}

//BMP ...
type BMP struct {
	Head       BMPHead
	Info       BMPInfo
	ColorTable []BMPColorTable
	Data       []byte
}

//ReadBMP 读取bmp图片
func ReadBMP(p string) (*BMP, error) {
	var (
		bmp BMP
		err error
	)

	// binary.Open(p)
	f, err := os.Open(p)
	if nil != err {
		return nil, err
	}
	bs, err := ioutil.ReadAll(f)
	if nil != err {
		return nil, err
	}

	// binary.Read(bytes.NewReader(bs[:2]), binary.LittleEndian, &bmp.Head.Type)
	// binary.Read(bytes.NewReader(bs[2:6]), binary.LittleEndian, &bmp.Head.Size)
	// binary.Read(bytes.NewReader(bs[6:8]), binary.LittleEndian, &bmp.Head.R1)
	// binary.Read(bytes.NewReader(bs[8:10]), binary.LittleEndian, &bmp.Head.R2)
	// binary.Read(bytes.NewReader(bs[10:14]), binary.LittleEndian, &bmp.Head.Offset)
	//提取文件头
	binary.Read(bytes.NewReader(bs[0:14]), binary.LittleEndian, &bmp.Head)

	//提取图片信息
	binary.Read(bytes.NewReader(bs[14:54]), binary.LittleEndian, &bmp.Info)

	//获取图片数据
	bmp.Data = make([]byte, bmp.Info.ImageSize)
	binary.Read(bytes.NewReader(bs[bmp.Head.Offset:]), binary.LittleEndian, &bmp.Data)

	//处理配色表
	if bmp.Head.Offset != 54 {
		bmp.ColorTable = make([]BMPColorTable, (bmp.Head.Offset-54)/4)
	}

	return &bmp, err
}

func (my BMP) String() string {
	return fmt.Sprintf(`
	bitmap type: %s,
	file size: %d,
	the start location of bitmap data: %d,
	image width: %d,
	image height: %d,
	Compression: %d,
	Bit per Pixel:%d,
	Image size: %d,
	Extracted size: %d
	`,
		my.Head.Type, my.Head.Size, my.Head.Offset,
		my.Info.Width, my.Info.Height, my.Info.Compression,
		my.Info.BitPerPixel, my.Info.ImageSize, len(my.Data))
}

func (my *BMP) ToGray() error {
	// var gray image.Gray
	return nil
}
