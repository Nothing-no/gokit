package mgodb

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type testData1 struct {
	ID   string  `bson:"_id"`
	Name string  `bson:"name"`
	Age  float64 `bson:"age"`
}

type testData2 struct {
	ID    float64 `bson:"_id"`
	Name  string  `bson:"name"`
	Age   float64 `bson:"age"`
	Phone string  `bson:"phone"`
}

type testData3 struct {
	ID    int     `bson:"_id"`
	Name  string  `bson:"name"`
	Age   float64 `bson:"age"`
	Phone string  `bson:"phone"`
	EB    string  `bson:"education"`
}

//ok
func TestInsert(t *testing.T) {
	ts := []struct {
		col  string
		data interface{}
		err  error
	}{
		{"test", &testData1{
			ID:   "123",
			Name: "zhijian",
			Age:  24,
		}, nil},
		{"test", &testData2{
			ID:    123,
			Name:  "zhijian",
			Age:   24,
			Phone: "1234r5",
		}, nil},
		{"test", &testData1{
			ID:   "12",
			Name: "zhijian",
			Age:  24,
		}, nil},
		{"test", &testData2{
			ID:    12,
			Name:  "zhijian",
			Age:   24,
			Phone: "1234r5",
		}, nil},
		//id冲突
		// {"test", &testData3{
		// 	ID:    123,
		// 	Name:  "zhijian",
		// 	Age:   24,
		// 	Phone: "ddsadaf",
		// 	EB:    "edadaf",
		// }, nil},
		// {"test", &testData1{
		// 	ID:   "123",
		// 	Name: "nothing",
		// 	Age:  24,
		// }, nil},
	}

	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	for _, t0 := range ts {
		err := mdb.Insert(t0.col, t0.data)
		if err != t0.err {
			t.Errorf("want %v, get %v", t0.err, err)
		}
	}
}

//ok
func TestSeachOne(t *testing.T) {
	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	v, err := mdb.SearchOne("test", "123")
	if nil != err {
		t.Errorf("want nil but get:%v", err)
	}
	t.Log(v)
}

//ok
func TestSeachAll(t *testing.T) {
	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	v, err := mdb.SearchAll("test", nil)
	if nil != err {
		t.Errorf("want nil but get: %v", err)
	} else if nil != v {
		// t.Error(v)
	}
}

//ok
func TestGetPage(t *testing.T) {
	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	v, err := mdb.GetPage("test", 1, 2)
	if nil != err {
		t.Errorf("want nil but get: %v", err)
	} else if nil != v {
		// t.Error(v)
	}

}

//ok
func TestUpdate(t *testing.T) {
	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	err := mdb.Update("test", bson.M{"_id": "123"}, bson.M{"$set": bson.M{"name": "update"}})
	if nil != err {
		t.Error(err)
	}
}

//ok
func TestHowMany(t *testing.T) {
	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	l, err := mdb.HowMany("test")
	if nil != err {
		t.Error(err)
	} else if 0 != l {
		// t.Error(l)
	}
}

func TestCountWithCond(t *testing.T) {
	mdb, _ := Init("10.2.3.141:27017", "zhijian")
	l, err := mdb.CountWithCond("report", bson.M{"isvalid": true})
	if nil != err {
		t.Error(err)
	} else {
		t.Error("true", l)
	}
}
