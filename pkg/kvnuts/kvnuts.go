package kvnuts

import (
	"bytes"
	"datastruck"
	"encoding/binary"
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

//put key values to nutsDB
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
			case byte:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case byte:
				keyByte = key.([]byte)
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

//get key to nutsDB
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
			case byte:
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
			case byte:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case byte:
				keyByte = key.([]byte)
			}
			return tx.SAdd(bct, keyByte, valByte)
		})
	if err != nil {
		return err
	}
	return nil
}

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
			case byte:
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
			case byte:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case byte:
				keyByte = key.([]byte)
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
			case byte:
				keyByte = key.([]byte)
			}

			switch val.(type) {
			case string:
				valByte = []byte(val.(string))
			case int:
				valByte, _ = IntToBytes(val.(int), 3)
			case byte:
				keyByte = key.([]byte)
			}
			if err := tx.RPush(bct, keyByte, valByte); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

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
	err = db.Update(
		func(tx *nutsdb.Tx) error {
			switch key.(type) {
			case string:
				keyByte = []byte(key.(string))
			case int:
				keyByte, _ = IntToBytes(key.(int), 3)
			case byte:
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

//turn byte to int
func BytesToInt(b []byte, isSymbol bool) (int, error) {
	if isSymbol {
		return bytesToIntS(b)
	}
	return bytesToIntU(b)
}

//字节数(大端)组转成int(无符号的)
func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

//字节数(大端)组转成int(有符号)
func bytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

//整形转换成字节
func IntToBytes(n int, b byte) ([]byte, error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 3, 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytes b param is invaild")
}
