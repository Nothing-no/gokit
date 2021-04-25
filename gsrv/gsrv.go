package gsrv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	Client interface {
		Header(key string, value string) Client
		Do(string, string, ...interface{}) ([]byte, error)
		PostFormFile(url string, file string) ([]byte, error)
	}

	cli struct {
		header map[string]string
	}

	Srv interface {
		AddAPI(string, string, func(*gin.Context)) Srv
		G(string) Group
		GetEngine() *gin.Engine
		Run() error
		StaticFs(string) Srv
	}

	Group interface {
		AddAPI(method string, path string, handler func(*gin.Context)) Group
		Done() Srv
		Use(...gin.HandlerFunc) Group
	}

	srv struct {
		e    *gin.Engine
		gs   map[string]*group
		port string
		ip   string
	}
	group struct {
		s *srv
		g *gin.RouterGroup
	}

	Resp struct {
		ErrNum int         `json:"errcode"`
		ErrMsg string      `json:"errmsg"`
		Data   interface{} `json:"data"`
	}
)

var _ Srv = &srv{}
var _ Group = &group{}
var _ Client = &cli{}

func NewClient() Client {
	return &cli{
		header: make(map[string]string),
	}
}
func (my *cli) Header(name string, value string) Client {
	my.header[name] = value
	return my
}
func (my *cli) Do(method string, url string, data ...interface{}) ([]byte, error) {
	var (
		bs  []byte
		err error
		req *http.Request
	)
	method = strings.ToUpper(method)
	if len(data) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else if len(data) == 1 {
		bs, err = json.Marshal(data[0])
		if nil != err {
			return bs, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bs))
	} else if len(data) == 2 {
		ireader, ok := data[0].(io.Reader)
		if !ok {
			ireader, ok = data[1].(io.Reader)
		}
		req, err = http.NewRequest(method, url, ireader)
	}
	if nil != err {
		return []byte{}, err
	}

	//设置头,设置完成后删掉已有的
	for k, v := range my.header {
		req.Header.Set(k, v)
		delete(my.header, k)
	}

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (my *cli) PostFormFile(url string, file string) ([]byte, error) {
	var (
		fileWriter  io.Writer
		fileHandler *os.File
		err         error
		resp        *http.Response
		contType    string
		result      []byte
	)
	// if len(files) == 0 {
	// 	return []byte{}, errors.New("lack of file")
	// }
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// defer bodyWriter.Close()/

	fileWriter, err = bodyWriter.CreateFormFile("files", filepath.Base(file))
	if nil != err {
		goto errReturn
	}

	fileHandler, err = os.Open(file)
	if nil != err {
		goto errReturn
	}
	defer fileHandler.Close()

	_, err = io.Copy(fileWriter, fileHandler)
	if nil != err {
		goto errReturn
	}

	contType = bodyWriter.FormDataContentType()
	//如果不在这里关闭，会导致对方读取文件出错
	bodyWriter.Close()

	resp, err = http.Post(url, contType, bodyBuf)
	if nil != err {
		goto errReturn
	}
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		goto errReturn
	}

	return result, nil

errReturn:
	return []byte{}, err
}

//NewSrv 如果传入了错误的port或者ip格式，则默认为127.0.0.1:0
func NewSrv(ip, port string) Srv {
	err := checkAddr(ip)
	if nil != err {
		ip = "127.0.0.1"
	}
	_, err = strconv.Atoi(port)
	if nil != err {
		port = "0"
	}

	return &srv{
		e:    gin.Default(),
		gs:   make(map[string]*group),
		ip:   ip,
		port: port,
	}
}

func (my *srv) AddAPI(method string, path string, handler func(*gin.Context)) Srv {
	my.e.Handle(strings.ToUpper(method), path, handler)
	return my
}

func (my *srv) G(gname string) Group {
	if _, ok := my.gs[gname]; !ok {
		my.gs[gname] = &group{
			s: my,
			g: my.e.Group(gname),
		}
	}

	return my.gs[gname]
}

func (my *srv) StaticFs(p string) Srv {
	my.e.StaticFS("/static/", http.Dir(p))
	return my
}
func checkAddr(addr string) error {

	newErr := fmt.Errorf("invalid address format:%s", addr)

	elemts := strings.Split(addr, ".")
	if len(elemts) != 4 {
		return newErr
	}
	for _, e := range elemts {
		i, err := strconv.Atoi(e)
		if nil != err {
			return newErr
		}
		if i > 255 {
			return newErr
		}
	}
	return nil
}

func (my *srv) GetEngine() *gin.Engine {
	return my.e
}

func (my *srv) Run() error {
	return my.e.Run(my.ip + ":" + my.port)
}

func (my *group) AddAPI(method string, path string, handler func(*gin.Context)) Group {
	my.g.Handle(strings.ToUpper(method), path, handler)
	return my
}

func (my *group) Done() Srv {
	return my.s
}

func (my *group) Use(midware ...gin.HandlerFunc) Group {
	my.g.Use(midware...)
	return my
}

func (my Resp) RJSON(c *gin.Context) {
	c.JSON(200, gin.H{
		"errcode": my.ErrNum,
		"errmsg":  my.ErrMsg,
		"data":    my.Data,
	})
}
