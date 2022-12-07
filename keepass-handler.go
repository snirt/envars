package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

type KeePassHandler struct {
	filename string
	db       *gokeepasslib.Database
	file     *os.File
}

var kph *KeePassHandler

// func init() {
// 	kph = New()
// }

// creates KeePassHandler object with unlocked db. MAKE SURE YOU CLOSE IT!
func New() *KeePassHandler {
	kph := new(KeePassHandler)

	kph.filename = ".env.kdbmx"

	file, err := os.Open(kph.filename)
	if err != nil {
		if os.IsNotExist(err) {
			kph.createDatabase()
		}
	} else {
		kph.unlockDB()
	}
	defer file.Close()
	kph.file = file

	return kph
}

func GetKeyPassHandler() *KeePassHandler {
	return kph
}

func (kph *KeePassHandler) unlockDB() {
	filename := kph.filename
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db := kph.db

	// try to unlock the file with password input
	for db.UnlockProtectedEntries() == nil {
		password := ReadInput("Please enter your DB password")
		db.Credentials = gokeepasslib.NewPasswordCredentials(password)
		_ = gokeepasslib.NewDecoder(file).Decode(db)
	}
}

func (kph *KeePassHandler) lockDB() {
	db := kph.db
	db.LockProtectedEntries()

	file, err := os.Open(kph.filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		file.Close()
		log.Printf("db file encoded and closed")
	}()

	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(db); err != nil {
		panic(err)
	}

}

func (kph *KeePassHandler) createDatabase() {
	fmt.Println("Creating a new database...")
	file, err := os.Create(kph.filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var pwd string
	for {
		pwd = ReadInput("Choose a password")
		vpwd := ReadInput("Verify password")
		if pwd == vpwd {
			break
		}
	}

	// init db with root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root group"

	// create db to contain the root group
	kph.db = &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(pwd),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}
}

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}

func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value, Protected: w.NewBoolWrapper(true)},
	}
}
