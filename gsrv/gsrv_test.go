package gsrv

import (
	"encoding/json"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient()
	bs, err := c.Do("get", "http://127.0.0.1:12308/test")
	if nil != err {
		t.Error(err)
	}
	t.Error(string(bs))
	mp := struct {
		V1 int64       `json:"v1"`
		V2 string      `json:"v2"`
		V3 interface{} `json:"v3"`
	}{V1: 234, V2: "helo", V3: map[string]interface{}{"test": 1}}
	bs, err = json.Marshal(&mp)
	if nil != err {
		t.Error(err)
	}
	fd, _ := os.Open("./1.txt")
	bs, err = c.Header("Content-Type", "application/text").Do("post", "http://127.0.0.1:12308/hello", fd, nil)
	if nil != err {
		t.Error(err)
	}
	t.Error(string(bs))

}

func TestMarshal(t *testing.T) {
	bs, err := json.Marshal(nil)
	if nil != err {
		t.Error(err)
	}
	t.Error(string(bs))
	bs, err = json.Marshal("defaqe test")
	if nil != err {
		t.Error(err)
	}
	t.Error(string(bs))
	bs, err = json.Marshal(map[string]interface{}{"v1": 1, "v2": "test"})
	if nil != err {
		t.Error(err)
	}
	t.Error(string(bs))
	bs, err = json.Marshal([]interface{}{1, 2, 3, "123"})
	if nil != err {
		t.Error(err)
	}
	t.Error(string(bs))

}
