package mgodb

import (
	"logsys/ec"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	MDB_URL string = ":27017"
)

//MgoDB 数据库结构体
type MgoDB struct {
	Url    string
	DbName string
	*mgo.Session
	*mgo.Database
}

//Init 初始化数据库
func Init(url, dbName string) (*MgoDB, error) {
	var (
		err error
		s   *mgo.Session
	)
	if "" == url && "" == dbName {
		//当传入参数均为对应的空值，采用默认的连接到本地的方式
		s, err = mgo.Dial(MDB_URL)
		if nil != err {
			ec.Info(err.Error())
			return nil, err
		}
		url = MDB_URL
		dbName = "test"
	} else if "" == url {
		//当传入的url为空值时，默认连接本地端口号为27017的mangodb
		s, err = mgo.Dial(MDB_URL)
		if nil != err {
			ec.Info(err.Error())
			return nil, err
		}
		url = MDB_URL
	} else if "" == dbName {
		//传入dbName为对应空值，连接url所在会话，默认打开test数据库
		s, err = mgo.Dial(url)
		if nil != err {
			ec.Info(err.Error())
			return nil, err
		}
		dbName = "test"
	} else {
		s, err = mgo.Dial(url)
		if nil != err {
			ec.Info(err.Error())
			return nil, err
		}
	}
	db := s.DB(dbName)

	mdb := MgoDB{
		Url:      url,
		DbName:   dbName,
		Session:  s,
		Database: db,
	}
	return &mdb, nil
}

//Insert 插入文档到指定的集合中，若指定集合不存在，则新建，成功返回nil，失败返回error
func (my *MgoDB) Insert(cName string, data ...interface{}) error {
	err := my.C(cName).Insert(data...)
	if nil != err {
		ec.Info(err.Error())
		return err
	}

	return nil
}

//SearchOne 通过id来查找数据
func (my *MgoDB) SearchOne(cName, id string) (interface{}, error) {
	var (
		result interface{} //查找结果
		err    error       //错误值
	)

	err = my.C(cName).FindId(id).One(&result)
	if nil != err {
		ec.Info(err.Error())
		return nil, err
	}

	return result, nil
}

//SearchAll 在集合cName中查找所有满足q的文档（数据）
func (my *MgoDB) SearchAll(cName string, q interface{}) ([]interface{}, error) {
	var (
		result []interface{}
		err    error
	)

	err = my.C(cName).Find(q).All(&result)
	if nil != err {
		ec.Info(err.Error())
		return nil, err
	}

	return result, err
}

//GetPage 返回一页数据，page指定页数(最小是0页)，size指定一页的大小
func (my *MgoDB) GetPage(cName string, page int, size int) ([]interface{}, error) {
	var (
		result []interface{}
		err    error
	)
	err = my.C(cName).Find(nil).Skip(page * size).Limit(size).All(&result)
	if nil != err {
		ec.Info(err.Error())
		return nil, err
	}

	return result, nil
}

//Search 查找集合中满足q的前limit条文档（数据—）
func (my *MgoDB) Search(cName string, q interface{}, page int, size int) ([]interface{}, error) {
	var (
		result []interface{}
		err    error
	)

	err = my.C(cName).Find(q).Skip(page * size).Limit(size).All(&result)
	if nil != err {
		ec.Info(err.Error())
		return nil, err
	}
	return result, nil
}

//Update 在指定的集合中修改数据
func (my *MgoDB) Update(cName string, where, modify interface{}) error {
	err := my.C(cName).Update(where, modify)
	if nil != err {
		ec.Info(err.Error())
		return err
	}

	return nil
}

//HowMany 返回一个集合中有多少个文档
func (my *MgoDB) HowMany(cName string) (int, error) {
	return my.C(cName).Find(nil).Count()
}

//CountWithCond ...
func (my *MgoDB) CountWithCond(cName string, cond interface{}) (int, error) {
	return my.C(cName).Find(cond).Count()
}

//FakeDel 假删除，若status为1，则说明是已删除数据
func (my *MgoDB) FakeDel(cName string, data interface{}) bool {
	err := my.C(cName).Update(data, bson.M{"$set": bson.M{"status": 1}})
	if nil != err {
		ec.Info(err.Error())
		return false
	}

	return true
}
