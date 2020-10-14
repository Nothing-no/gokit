package xls

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

//Border line
const (
	//无线型
	NONE = iota
	//细实线
	LightLine
	//粗实线
	HeavyLine
	//虚线
	ImaginaryLine
	//点线
	DottedLine
)

//Color 采用RGB格式
const (
	//黑色
	Black = "000000"
	//红色
	Red = "FF0000"
	//绿色
	Green = "00FF00"
	//蓝色
	Blue = "0000FF"
	//浅蓝色
	LightBlue = "BDD7EE"
	//浅黄色
	LightYellow = "FFF68F"
	//浅灰色
	LightGray = "E8E8E8"
	//西红柿色
	Tomato = "EE5C42"
)

//填充， 纯色或不填充
const (
	NoneFill = iota
	PureFill
)

//SheetStyle ....
type SheetStyle struct {
	Name string
	Each []*StyleFmt
	*excelize.File
	*sync.WaitGroup

	// PreStr interface{}
}

//StyleFmt  format for sheetcell
type StyleFmt struct {
	Axis    [][2]string
	Merge   bool
	StyleID int
}

//Workbook 工作簿信息
type Workbook struct {
	*excelize.File          //文件，考虑到要使用原来的一些方法
	SheetNames     []string //工作簿的所有sheet名字集
	Path           string   //工作簿的路径
	*sync.WaitGroup
	*sync.RWMutex
}

//SheetRead 读取数据工作页
type SheetRead struct {
	Name   string
	Data   []interface{}
	Style  []int
	Mother *Workbook
}

//SheetWrite 写数据工作页
type SheetWrite struct {
	Name          string
	GlobalData    map[string]interface{} //全局的一些配置信息
	Data          []*MapRowData          //需要写入的数据块
	Mother        *Workbook
	*sync.RWMutex //并发控制

}

//MapRowData ---
type MapRowData struct {
	RowData  []interface{}
	FirstLoc string
}

//Init 初始化
func Init(path ...string) (*Workbook, error) {
	var (
		wb  *Workbook
		f   *excelize.File
		err error
	)

	pLen := len(path)
	// fmt.Println(pLen)
	if 0 == pLen {
		f = excelize.NewFile()
		wb = &Workbook{f, f.GetSheetList(), "", &sync.WaitGroup{}, &sync.RWMutex{}}
	} else if 1 == pLen {
		f, err = excelize.OpenFile(path[0])
		if nil != err {
			return nil, err
		}
		wb = &Workbook{f, f.GetSheetList(), path[0], &sync.WaitGroup{}, &sync.RWMutex{}}
	} else if pLen > 1 {
		return nil, errors.New("Too many params for InitXls()")
	}

	return wb, nil
}

//NewOneWriteBuffer 新建一个写表的buffer
//若是提供sheet名，则返回该sheetbuffer
//若没有或提供表名错误，则返回workbook的第一张表buffer
func (my *Workbook) NewOneWriteBuffer(sheetName ...string) *SheetWrite {
	var (
		sw *SheetWrite
	)

	sLen := len(sheetName)
	if 0 == sLen {
		sw = &SheetWrite{
			Name:       my.GetSheetName(0),
			GlobalData: make(map[string]interface{}),
			Mother:     my,
			RWMutex:    &sync.RWMutex{},
		}
	} else {
		if IsExisted(my.SheetNames, sheetName[sLen-1]) {
			sw = &SheetWrite{
				Name:       sheetName[sLen-1],
				GlobalData: make(map[string]interface{}),
				Mother:     my,
				RWMutex:    &sync.RWMutex{},
			}
		} else {
			sw = &SheetWrite{
				Name:       my.GetSheetName(0),
				GlobalData: make(map[string]interface{}),
				Mother:     my,
				RWMutex:    &sync.RWMutex{},
			}
		}
	}

	return sw
}

//NewMulWriteBuffer 新建（整个工作簿）多张写表的buffers
//tip：可以将sheet名设置为索引常量，例如const Sheet1 = 0，这样也可以清楚知道自己在操作哪张表
func (my *Workbook) NewMulWriteBuffer() []*SheetWrite {
	var (
		sw []*SheetWrite
	)

	for _, v := range my.SheetNames {
		sw = append(sw, &SheetWrite{
			Name:       v,
			GlobalData: make(map[string]interface{}),
			Mother:     my,
			RWMutex:    &sync.RWMutex{},
		})
	}

	return sw
}

//IsExisted 判断elmt是否在arr中，若在，则返回true，不在，返回false
func IsExisted(arr []string, elmt string) bool {
	for _, v := range arr {
		if elmt == v {
			return true
		}
	}

	return false
}

//ToSave 保存,若没有输入参数，则保存为workbook.xlsx，若有多个参数，则用最后一个
func (my *Workbook) ToSave(name ...string) {
	nLen := len(name)
	if 0 == nLen {
		//如果没有给定路径，默认保存为worknook.xlsx
		if "" == my.Path {
			my.SaveAs("workbook.xlsx")
		} else {
			//如果有给定路径名，还是保存为原来的
			my.SaveAs(my.Path)
		}

		return
	}

	if strings.Contains(name[nLen-1], ".xls") {
		my.SaveAs(name[nLen-1])
	} else {
		my.SaveAs(name[nLen-1] + ".xlsx")
	}
}

//Flush data to excel
func (my *SheetWrite) Flush() (err error) {
	// my.Mother.SetActiveSheet(my.Mother.GetSheetIndex(my.Name))
	my.Mother.Lock()
	defer my.Mother.Unlock()
	for k, v := range my.GlobalData {
		err = my.Mother.SetCellValue(my.Name, k, v)
		if nil != err {
			return err
		}
	}

	for _, data := range my.Data {
		err = my.Mother.SetSheetRow(my.Name, data.FirstLoc, &(data.RowData))
		if err != nil {
			return err
		}
	}
	return nil
}

//OpenStyleSetting ..打开一个样式设置起
func (my *Workbook) OpenStyleSetting(name string) *SheetStyle {
	return &SheetStyle{
		name,
		make([]*StyleFmt, 0),
		my.File,
		my.WaitGroup,
	}
}

//Apply 设置表的一些样式
func (my *SheetStyle) Apply() {
	for _, v := range my.Each {
		my.setStyle(v.Axis, v.StyleID, v.Merge)
	}
}

func (my *SheetStyle) setStyle(axis [][2]string, id int, merge bool) {
	for _, v := range axis {
		if merge {
			my.MergeCell(my.Name, v[0], v[1])
		}
		my.SetCellStyle(my.Name, v[0], v[1], id)
	}
}

//AssignStyle 分配样式
func (my *SheetStyle) AssignStyle(axis *[][2]string, merge bool, p interface{}) {
	id, err := my.NewStyle(p)
	if nil != err {
		fmt.Println(err)
	}
	my.Each = append(my.Each, &StyleFmt{
		Axis:    *axis,
		Merge:   merge,
		StyleID: id,
	})
}

//SetCellDim ...
func (my *SheetStyle) SetCellDim(st, ed string, w, h float64) error {
	stCol, stRow := SplitAxis(st)
	edCol, edRow := SplitAxis(ed)
	err := my.SetColWidth(my.Name, stCol, edCol, w)
	if nil != err {
		return err
	}
	for i := stRow; i <= edRow; i++ {
		err = my.SetRowHeight(my.Name, i, h)
		if nil != err {
			return err
		}
	}
	return nil
}

//GetIncAxis 根据start坐标，获取增加step后的坐标，
//默认返回的是增加step行，将coldir选中为true，则是返回增加step列
func GetIncAxis(start string, step int, coldir ...bool) string {
	var (
		l        = len(coldir)
		col, row = SplitAxis(start)
	)

	//在列上步进
	if (0 != l) && (true == coldir[l-1]) {
		col = add26(col, step)
	} else {
		row += step
	}

	return col + strconv.Itoa(row)
}

//SplitAxis 分割行列
func SplitAxis(axis string) (string, int) {
	var (
		col    string
		rowStr string
		row    int
		l      = len(axis)
	)
	alpha := func(r rune) bool {
		return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
	}

	if 0 == strings.IndexFunc(axis, alpha) {
		lastColIndex := strings.LastIndexFunc(axis, alpha)
		if (lastColIndex >= 0) && (lastColIndex < l) {
			col, rowStr = axis[:lastColIndex+1], axis[lastColIndex+1:]
			row, _ = strconv.Atoi(rowStr)
		}
	}

	return col, row

}

//GetAxisForStyle 设置统一的坐标
func (my *SheetWrite) GetAxisForStyle(src [][2]string, target string, step int, col ...bool) [][2]string {
	tmp, _ := my.Mother.SearchSheet(my.Name, target)
	l := len(col)
	for _, v := range tmp {
		if 0 == l {
			src = append(src, [2]string{
				v,
				GetIncAxis(v, step),
			})
		} else {
			src = append(src, [2]string{
				v,
				GetIncAxis(v, step, col[l-1]),
			})
		}

	}
	return src
}

//AddNewRow 增加一个与前行格式一样的复制行
func (my *SheetWrite) AddNewRow(row int) {
	err := my.Mother.DuplicateRow(my.Name, row)
	if nil != err {
		fmt.Println(err)
		return
	}
}

//add26 26进制加法A代表1，Z为26
func add26(base string, inc int) string {
	var (
		base10 int //base所表示的十进制数
		bs     = []byte(base)
		rbuf   []byte
	)

	l := len(base)

	//转成大写及其10进制表示
	for i := 0; i < l; i++ {
		if bs[i] >= 'a' && bs[i] <= 'z' {
			bs[i] = (bs[i] - 'a') + 'A'
		}

		base10 = 26*base10 + int(bs[i]-'A'+1)
	}

	incResult := base10 + inc

	for incResult != 0 {
		c0 := (incResult - 1) % 26
		incResult = incResult / 26
		rbuf = append(rbuf, byte(c0)+'A')
	}
	lr := len(rbuf)
	for i := 0; i < lr/2; i++ {
		rbuf[i], rbuf[lr-1-i] = rbuf[lr-1-i], rbuf[i]
	}
	return string(rbuf)
}
