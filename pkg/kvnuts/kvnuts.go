package kvnuts

import (
	"datastruck"
	ErrH "errorhandle"
	"fmt"
	"github.com/xujiajun/nutsdb"
	"log"
)

//connect to nutsDB
func connect() (db *nutsdb.DB) {
	var (
		err error
		c   datastruck.Config
	)

	if err = c.Config(); err != nil {
		log.Println(ErrH.ErrorLog(12012), fmt.Sprintf("%v", err))
		return
	}
	opt := nutsdb.DefaultOptions
	opt.Dir = c.NutsDB.Path

	if db, err = nutsdb.Open(opt); err != nil {
		log.Println(ErrH.ErrorLog(12161), fmt.Sprintf("%v", err))
		return
	}
	return
}

//put key values
func Put(bct string, key, val interface{}) error {
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
		log.Println(ErrH.ErrorLog(12162), fmt.Sprintf("%v", err))
		return err
	}
	return nil
}

//get key
func Get(bct string, key interface{}, valType string) (err error, myReturn string, myReturnInt int) {
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

			if e, err = tx.Get(bct, keyByte); err != nil {
				return err
			} else {
				switch valType {
				case "s":
					myReturn = string(e.Value)
				case "i":
					myReturnInt, _ = BytesToInt(e.Value, true)
				default:
					err = fmt.Errorf("my error")
					return err
				}
			}
			return nil
		})
	if err != nil {
		log.Println(ErrH.ErrorLog(12163), fmt.Sprintf("%v", err))
		return
	}
	return
}

//del key
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

//put key, value as set
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

//get all key from set
func SMem(bct string, key interface{}) (error, [][]byte) {
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
		return err, nil
	}
	return nil, items
}

//remove key, value from set
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

//judge member if exist
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

//put key, value as list
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

//get key from list
// s and e as index of list,
func LIndex(bct string, key interface{}, s, e int) (error, [][]byte) {
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
		return err, nil
	}
	return nil, item
}
