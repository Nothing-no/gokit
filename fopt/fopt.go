package fopt

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type (
	ZipPrepare struct {
		ZipName  string
		SrcFiles []string
		realSrc  []string
	}
	UnzipPrepare struct {
		OutPath string
		SrcPath string
	}
)

func (my *ZipPrepare) Zip() error {
	newZipFile, err := os.Create(my.ZipName)
	if nil != err {
		return err
	}
	defer newZipFile.Close()

	//新建一个zipwriter，供数据输入
	newZipWriter := zip.NewWriter(newZipFile)
	defer newZipWriter.Close()

	for _, f := range my.SrcFiles {
		r, err := getDirFiles(f)
		if nil != err {
			fmt.Println(err)
		} else {
			my.realSrc = append(my.realSrc, r...)
		}
	}

	for _, file := range my.realSrc {
		if err = addFileToZip(newZipWriter, file); nil != err {
			return err
		}
	}

	return nil
}

func getDirFiles(p string) ([]string, error) {
	var ret []string
	fi, err := os.Stat(p)
	if nil != err {
		return ret, err
	}

	if fi.IsDir() {
		fs, err := ioutil.ReadDir(p)
		if nil != err {
			return ret, err
		}
		for _, f := range fs {
			if f.IsDir() {
				tmpRet, err := getDirFiles(filepath.Join(p, f.Name()))
				if nil != err {
					return ret, err
				}
				ret = append(ret, tmpRet...)
			}
			ret = append(ret, filepath.Join(p, f.Name()))
		}
	} else {

		ret = append(ret, p)
	}

	return ret, nil
}

func addFileToZip(zw *zip.Writer, fn string) error {
	file, err := os.Open(fn)
	if nil != err {
		return err
	}
	defer file.Close()

	//获取file info
	info, err := file.Stat()
	if nil != err {
		return err
	}

	zipHeader, err := zip.FileInfoHeader(info)
	if nil != err {
		return err
	}
	// filepath.Clean(fn)
	zipHeader.Name = fn
	zipHeader.Method = zip.Deflate
	w, err := zw.CreateHeader(zipHeader)
	if nil != err {
		return err
	}
	_, err = io.Copy(w, file)

	return err
}

func dealName(name string) string {
	ret := ""
	cleanName := filepath.Clean(name)
	sep := string(os.PathSeparator)
	pelmts := strings.Split(cleanName, sep)
	fmt.Println(pelmts)
	for _, v := range pelmts {
		if v == ".." {
			continue
		} else {
			ret += v + sep
		}
	}
	return ret

}

func (my *UnzipPrepare) Unzip() error {
	newZipReader, err := zip.OpenReader(my.SrcPath)
	if nil != err {
		fmt.Println(err)
		return err
	}
	defer newZipReader.Close()

	for _, f := range newZipReader.File {
		//给定解压完整路径
		fmt.Println(my.OutPath, filepath.Clean(f.Name))
		fp := filepath.Join(my.OutPath, dealName(f.Name))
		prefix := filepath.Clean(my.OutPath) + string(os.PathSeparator)
		// fmt.Println(fp, prefix)

		//判断路径包含指定路径前缀
		if !strings.HasPrefix(fp, prefix) {
			return fmt.Errorf("%s: 不合法路径", fp)
		}

		//判断当前是否为文件夹
		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}

		err = os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if nil != err {
			return err
		}

		outFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if nil != err {
			return err
		}

		zipFile, err := f.Open()
		if nil != err {
			return err
		}

		_, err = io.Copy(outFile, zipFile)
		outFile.Close()
		zipFile.Close()

		if nil != err {
			return err
		}

	}

	return nil
}
