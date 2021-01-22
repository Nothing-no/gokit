# gokit
### ec 错误信息收集，生成错误log，无缓冲
>```go
> //最好在main函数添加以下if语句
> if ec.GetECStatus() {
>   ec.Close()    
> }
> ec.Debug("debug messge")
> ec.Errorf("info message")
>```
### xls 操作excel,
>```go
>//Init后，一些方法可参考 github.com/360EntSecGroup-Skylar/excelize/v2
> wb := xls.Init()  //带参数则为打开文件，不带参数则为新建文件
> wb.ToSave()   //带参数存为参数名，不带参数默认存为workbook.xlsx
> ```
