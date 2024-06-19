package config

import (
	"encoding/json"
	"os"
	"reflect"
	"star_dust/inter"
	"sync"
)

type UserData struct {
	fileName string
	fileType string
	filePath string
	locker   sync.Mutex
	data     struct {
		UserName   string `json:"user_name"`
		Health     int64  `json:"health"`
		BlueAmount int64  `json:"blue_amount"`
		Grade      int64  `json:"grade"`
		Money      int64  `json:"money"`
	}
}

func (u *UserData) Store() error {
	u.locker.Lock()
	defer u.locker.Unlock()

	configFile, err := os.Create(u.filePath + "/" + u.fileName + "." + u.fileType)
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(&(u.data)); err != nil {
		return err
	}
	return nil
}

func (u *UserData) Read() error {
	u.locker.Lock()
	defer u.locker.Unlock()

	configFile, err := os.Create(u.filePath + "/" + u.fileName + "." + u.fileType)
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&(u.data)); err != nil {
		return err
	}
	return nil
}

func (u *UserData) Get(field string) (interface{}, bool) {
	u.locker.Lock()
	defer u.locker.Unlock()

	v := reflect.ValueOf(u.data)
	f := reflect.Indirect(v).FieldByName(field)

	if !f.IsValid() {
		return nil, false
	}
	return f.Interface(), true
}

func (u *UserData) Set(field string, value interface{}) bool {
	u.locker.Lock()
	defer u.locker.Unlock()

	v := reflect.ValueOf(&u.data).Elem()
	f := v.FieldByName(field)

	if !f.IsValid() {
		return false
	}

	if !f.CanSet() {
		return false
	}

	val := reflect.ValueOf(value)
	if f.Type() != val.Type() {
		return false
	}

	f.Set(val)
	return true
}

func NewUserData(fileName string, fileType string, filePath string) inter.Configurator {
	u := &UserData{
		fileName: fileName,
		fileType: fileType,
		filePath: filePath,
	}
	u.Read()
	return u
}

var UserDataConfigurator = NewUserData("user_data", "json", "/Users/admin/Desktop/star_dust/data")
