package ginsrv

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	newErr = errors.New("Wrong ip format: ")
	router = gin.Default()
)

type (
	Route struct {
		Method  string
		Path    string
		Handler gin.HandlerFunc
	}

	Group struct {
		Path   string
		flag   bool
		routes []Route
		Router *gin.RouterGroup
	}

	Info struct {
		Addr string `json:"domain"`
		Port string `json:"port"`
		// Rs     []Route		//routes
		Gs     []*Group //groups
		Engine *gin.Engine
	}

	Resp struct {
		ErrNum int         `json:"errcode"`
		ErrMsg string      `json:"errmsg"`
		Data   interface{} `json:"data"`
	}
)

func SetStatic(p string) {
	router.StaticFS("/static/", http.Dir(p))
}

func NewGroups(gp ...string) map[string]*Group {
	r := make(map[string]*Group)
	if len(gp) == 0 {
		r["/api/v1"] = &Group{
			Path: "/api/v1",
		}
	} else {
		for _, g := range gp {
			r[g] = &Group{
				Path: g,
			}
		}
	}
	return r
}

func (my *Group) SetRoutes(route ...Route) {
	if !my.flag {
		initGroup(my)
	}
	my.routes = route
}

func (my *Group) AddRoute() {

	for _, r := range my.routes {
		r.Method = strings.ToUpper(r.Method)
		my.Router.Handle(r.Method, r.Path, r.Handler)
	}
}

// func NewRoutes(rs ...) []*Route {
// 	l := len(rs)
// 	ret := make([]*Route, l)
// 	for i, r := range rs {
// 		ret[i] = &Route{
// 			Method:  r.p1,
// 			Path:    r.p2,
// 			Handler: r.p3,
// 		}
// 	}
// 	return ret
// }

func initGroup(g *Group) {
	g.Router = router.Group(g.Path)
	g.flag = true
}

func Init(p ...string) (*Info, error) {
	var info Info
	info.Engine = router
	switch len(p) {
	case 0:
		return nil, errors.New("please specify a port<eg:\"9999\">")
	case 1:
		info.Port = p[0]
	case 2:
		err := checkAddr(p[0])
		if nil != err {
			return nil, err
		}
		info.Port = p[1]
		info.Addr = p[0]
	default:
		return nil, errors.New("Too many params")
	}

	return &info, nil
}

func Run(gs map[string]*Group, p ...string) error {
	info, err := Init(p...)
	if nil != err {
		return err
	}
	info.Engine.Use(cors.Default())
	for _, g := range gs {
		g.AddRoute()
	}
	err = info.Engine.Run(info.Addr + ":" + info.Port)
	return err
}

func checkAddr(addr string) error {

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

func RunWithNoGroup(ip, port string, rs ...Route) {
	router.Use(cors.Default())
	for _, r := range rs {
		router.Handle(r.Method, r.Path, r.Handler)
	}

	router.Run(ip + ":" + port)
}

// func (my *Info) Run() {

// }

func TestHandle(rs []*Route) *gin.Engine {
	for _, r := range rs {
		r.Method = strings.ToUpper(r.Method)
		router.Handle(r.Method, r.Path, r.Handler)
	}

	return router
}

func ErrRespJson(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Resp{
		ErrNum: code,
		ErrMsg: msg,
	})
}

func GetLocalBoundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ip := strings.Split(localAddr.String(), ":")[0]
	return ip
}
