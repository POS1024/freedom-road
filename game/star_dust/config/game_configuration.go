package config

import (
	"encoding/json"
	"os"
	"reflect"
	"star_dust/inter"
	"sync"
)

type GameConfiguration struct {
	fileName string
	fileType string
	filePath string
	locker   sync.Mutex
	data     struct {
		Volume   int64  `json:"volume"`
		Language string `json:"language"`
	}
}

func (g *GameConfiguration) Store() error {
	g.locker.Lock()
	defer g.locker.Unlock()

	configFile, err := os.Create(g.filePath + "/" + g.fileName + "." + g.fileType)
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(&(g.data)); err != nil {
		return err
	}
	return nil
}

func (g *GameConfiguration) Read() error {
	g.locker.Lock()
	defer g.locker.Unlock()

	configFile, err := os.Open(g.filePath + "/" + g.fileName + "." + g.fileType)
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&(g.data)); err != nil {
		return err
	}
	return nil
}

func (g *GameConfiguration) Get(field string) (interface{}, bool) {
	g.locker.Lock()
	defer g.locker.Unlock()

	v := reflect.ValueOf(g.data)
	f := reflect.Indirect(v).FieldByName(field)

	if !f.IsValid() {
		return nil, false
	}
	return f.Interface(), true
}

func (g *GameConfiguration) Set(field string, value interface{}) bool {
	g.locker.Lock()
	defer g.locker.Unlock()

	v := reflect.ValueOf(&g.data).Elem()
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

func NewGameConfiguration(fileName string, fileType string, filePath string) inter.Configurator {
	g := &GameConfiguration{
		fileName: fileName,
		fileType: fileType,
		filePath: filePath,
	}
	g.Read()
	return g
}

var GameConfigurationConfigurator = NewGameConfiguration("game_configuration", "json", "/Users/admin/Desktop/star_dust/data")
