package kvnuts

import (
	"fmt"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	MyERR "github.com/ranzhendong/irishman/pkg/errorhandle"
	"github.com/xujiajun/nutsdb"
	"log"
	"time"
)

//connect to nutsDB
func connect() (db *nutsdb.DB) {
	var (
		c   datastruck.Config
		err error
	)

	if err = c.Config(); err != nil {
		log.Println(MyERR.ErrorLog(12012), fmt.Sprintf("%v", err))
		return
	}
	opt := nutsdb.DefaultOptions
	opt.Dir = c.NutsDB.Path

	for {
		time.Sleep(5 * time.Millisecond)
		if db, err = nutsdb.Open(opt); err == nil {
			goto GETDB
		} else {
			log.Println(MyERR.ErrorLog(12161), fmt.Sprintf("%v", err))
		}
	}
GETDB:
	return
}

//Put key values
func Put(bct string, key, val interface{}) error {
	var (
		keyByte, valByte []byte
		err              error
		db               *nutsdb.DB
	)

	db = connect()
	defer func() {
		_ = db.Close()
	}()
	err = db.Update(
		func(tx *nutsdb.Tx) error {
			//judge key type
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(val.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case []uint8:
				valByte = val.([]byte)
			}
			if err = tx.Put(bct, keyByte, valByte, 0); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Println(MyERR.ErrorLog(12162), fmt.Sprintf("%v", err))
		return err
	}
	return nil
}

//Get key
func Get(bct string, key interface{}, valType string) (myReturn string, myReturnInt int, err error) {
	var (
		keyByte []byte
		e       *nutsdb.Entry
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()
	err = db.View(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			if e, err = tx.Get(bct, keyByte); err == nil {
				switch valType {
				case "s":
					myReturn = string(e.Value)
				case "i":
					myReturnInt = BytesToInt(e.Value, true)
				default:
					err = fmt.Errorf("my error")
					return err
				}
				return nil
			}
			return err
		})
	if err != nil {
		//log.Println(MyERR.ErrorLog(12163), fmt.Sprintf("%v", err))
		return
	}
	return
}

//Del key
func Del(bct string, key interface{}) error {
	var (
		keyByte []byte
		err     error
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.Update(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}
			if err = tx.Delete(bct, keyByte); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

//SAdd Put key, but value as set
func SAdd(bct string, key, val interface{}) error {
	var (
		keyByte, valByte []byte
		err              error
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.Update(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case []uint8:
				valByte = val.([]byte)
			}
			return tx.SAdd(bct, keyByte, valByte)
		})
	if err != nil {
		return err
	}
	return nil
}

//SMem Get all key from set
func SMem(bct string, key interface{}) ([][]byte, error) {
	var (
		keyByte []byte
		err     error
		items   [][]byte
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.View(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			if items, err = tx.SMembers(bct, keyByte); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
	}
	return items, nil
}

//SRem Remove key, value from set
func SRem(bct string, key, val interface{}) error {
	var (
		keyByte, valByte []byte
		err              error
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.Update(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case []uint8:
				valByte = val.([]byte)
			}
			if err := tx.SRem(bct, keyByte, valByte); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

//SIsMem Judge member if exist
func SIsMem(bct string, key, val interface{}) bool {
	var (
		keyByte, valByte []byte
		err              error
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.View(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case []uint8:
				valByte = val.([]byte)
			}
			if _, err := tx.SIsMember(bct, keyByte, valByte); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return false
	}
	return true
}

//LAdd is put key, value as list
func LAdd(bct string, key, val interface{}) error {
	var (
		keyByte, valByte []byte
		err              error
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.Update(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 4)
			case []uint8:
				valByte = val.([]byte)
			}
			return tx.LPush(bct, keyByte, valByte)
		})
	if err != nil {
		return err
	}
	return nil
}

/*
LIndex is Get key from list
s and e as index of list,
*/
func LIndex(bct string, key interface{}, s, e int) ([][]byte, error) {
	var (
		keyByte []byte
		item    [][]byte
		err     error
	)
	db := connect()
	defer func() {
		_ = db.Close()
	}()

	err = db.View(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case []uint8:
				keyByte = key.([]byte)
			}

			if item, err = tx.LRange(bct, keyByte, s, e); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return item, nil
}
