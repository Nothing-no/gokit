#### Usage

```go
func main(){
    //zip example
    fp := fopt.ZipPrepare{
        ZipName:"test.zip",
        SrcFiles: []string["f1","f2"],
    }
    if err := fp.Zip(); nil != err {
        fmt.Println(err)
        return
    }

    //unzip example
    fp1 := fopt.UnzipPrepare{
        OutPath:"./testunzip",
        SrcPath:"./test.zip",
    }

    if err := fp1.Unzip(); nil != err {
        fmt.Println(err)
        return
    }

}
```