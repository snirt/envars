package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tobischo/gokeepasslib/v3"
)

type KeyPassHandler struct {
	DBPath string
	File *os.File
}

func (kph *KeyPassHandler) ListDB() {
	fmt.Println("nlanlsdnf")
}

func (kph *KeyPassHandler) CreateDB(name string, password string, folderPath string) {
	// try to get path from config
	kph.DBPath = viper.GetViper().GetString(K_KEEPASS_DB_PATH)

	//TODO check if the file from config exists

	kph.DBPath = filepath.Join(folderPath, name+".kdbx")

	if _, err := os.Stat(kph.DBPath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(folderPath, os.ModePerm)
		file, err := os.Create(kph.DBPath)
		if err != nil {
			panic(err)
		}
		kph.File = file
		defer kph.File.Close()

		// create rootGroup
		rootGroup := gokeepasslib.NewGroup()
		rootGroup.Name = "root group"
	} else {
		file, err := os.OpenFile(kph.DBPath, os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
		
		kph.File = file
		defer kph.File.Close()
	}
}

func (kph *KeyPassHandler)GetFile() *os.File {
	return kph.File
}

func (kph *KeyPassHandler) DeleteDB() {
	//TODO get DBs path from viper config file
	files, err := ioutil.ReadDir("./test")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
	//TODO continue DeleteDB implementation
}