//Package kvstore the leveldb
package kvstore

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/team4yf/iot-device-mqtt/pkg/utils"
)

var db *leveldb.DB

var lock sync.Mutex

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
}

func Get(key string) (data []byte, err error) {
	data, err = db.Get([]byte(key), nil)
	return
}

func Put(key string, data []byte) (err error) {
	err = db.Put([]byte(key), data, nil)
	return
}

func GetObject(key string) (data interface{}, err error) {
	buf, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	err = utils.DataToStruct(buf, data)
	return
}

func PutObject(key string, data interface{}) (err error) {
	buf := utils.Struct2Bytes(data)
	err = db.Put([]byte(key), buf, nil)
	return
}

func Close() {
	db.Close()
}
