package arch

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//Zip zip files
/* opts:
*	t: 将所有路径都去掉，将文件压缩在第一层目录下
*	e: 将保留压缩的目录，去掉目录前的路径
*	n: 去掉输入路径，压缩路径下的文件
*
 */
func Zip(dest string, src []string, opts ...string) error {

	newZipFile, err := os.Create(dest)
	if nil != err {
		return err
	}
	defer newZipFile.Close()

	newZipWriter := zip.NewWriter(newZipFile)
	defer newZipWriter.Close()

	// var (
	// 	iTmpSrc = []string{}
	// )

	for _, f := range src {
		// r, err := getDirFiles(f)
		// if nil != err {
		// 	return err
		// }
		// iTmpSrc = append(iTmpSrc, r...)
		filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
			zipHeader, err := zip.FileInfoHeader(info)
			if nil != err {
				return err
			}
			// filepath.Clean(fn)
			if len(opts) == 1 {

				switch opts[0] {
				//去掉所有路径，只保存当前文件
				case "t":
					zipHeader.Name = info.Name()
				//若带有目录，则去掉输入路径
				case "e":
					zipHeader.Name = strings.TrimPrefix(path, filepath.Dir(f)+string(os.PathSeparator))
				case "n":
					zipHeader.Name = strings.TrimPrefix(path, f+string(os.PathSeparator))
				default:
					zipHeader.Name = path
				}
			}
			if !info.IsDir() {
				fmt.Println(zipHeader.Name)
				zipHeader.Method = zip.Deflate
				w, err := newZipWriter.CreateHeader(zipHeader)
				if nil != err {
					return err
				}
				file, err := os.Open(path)
				if nil != err {
					return err
				}
				_, err = io.Copy(w, file)
			}

			return err
		})
	}

	return err
	// for _, file := range iTmpSrc {
	// 	if err = addFileToZip(newZipWriter, file, opts...); nil != err {
	// 		return err
	// 	}
	// }

	// return nil

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

func addFileToZip(zw *zip.Writer, fn string, opt ...string) error {
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
	if len(opt) == 1 {

		switch opt[0] {
		//去掉所有路径，只保存当前文件
		case "t":
			zipHeader.Name = info.Name()
		//若带有目录，则保留文件外面的一层目录
		case "s":
			flist := strings.Split(fn, string(os.PathSeparator))
			l := len(flist)
			if l >= 2 {
				zipHeader.Name = filepath.Join(flist[l-2], flist[l-1])
			} else {
				zipHeader.Name = fn
			}
		default:
			zipHeader.Name = fn
		}
	}
	fmt.Println(zipHeader.Name)
	zipHeader.Method = zip.Deflate
	w, err := zw.CreateHeader(zipHeader)
	if nil != err {
		return err
	}
	_, err = io.Copy(w, file)

	return err
}

func Unzip(dest string, src string, opts ...string) error {
	newZipReader, err := zip.OpenReader(src)
	if nil != err {
		return err
	}
	defer newZipReader.Close()

	for _, f := range newZipReader.File {
		fp := filepath.Join(dest, f.Name)
		prefix := filepath.Clean(dest) + string(os.PathSeparator)
		if !strings.HasPrefix(fp, prefix) {
			return fmt.Errorf("%s: 不合法路径", fp)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}

		err = os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if nil != err {
			return err
		}

		outFile, err := os.Create(fp)
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
