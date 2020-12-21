### use
```go
pdfFile := "./pdf/test.pdf"
bmpOutPath := "./bmp"
err := middle.Pdf2Bmp(pdfFile, bmpOutPath)
if nil != err {
    fmt.Println(err)
    return
}
```