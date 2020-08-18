//Package kvstore the leveldb
package kvstore

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/team4yf/iot-device-mqtt/pkg/utils"
)

var (
	db     *leveldb.DB
	lock   sync.Mutex
	inited bool
)

//Init  初始化数据库
func Init(dir string) {
	lock.Lock()
	defer lock.Unlock()

	if db != nil {
		return
	}
	var err error
	db, err = leveldb.OpenFile(dir, nil)

	if err != nil {
		panic(err)
	}
	inited = true
}

//Get 获取数据,如果不存在则第二个参数是 false
func Get(key string) ([]byte, bool) {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		// not found
		return nil, false
	}
	return data, true
}

//Put 插入数据
func Put(key string, data []byte) (err error) {
	err = db.Put([]byte(key), data, nil)
	return
}

//GetObject 获取到数据，并保存为对象
func GetObject(key string, data interface{}) (err error) {
	buf, err := db.Get([]byte(key), nil)
	if err != nil {
		// not found
		return
	}
	err = utils.DataToStruct(buf, data)
	return
}

//PutObject 将对象保存到数据库
func PutObject(key string, data interface{}) (err error) {
	buf := utils.Struct2Bytes(data)
	err = db.Put([]byte(key), buf, nil)
	return
}

//Close 关闭数据库
func Close() {
	if !inited {
		return
	}
	if db == nil {
		return
	}
	db.Close()
}
