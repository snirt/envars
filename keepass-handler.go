package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

type KeePassHandler struct {
	filename string
	db       *gokeepasslib.Database
	file     *os.File
}

var kph *KeePassHandler

// creates KeePassHandler object with unlocked db. MAKE SURE YOU CLOSE IT!
func New() *KeePassHandler {
	kph := new(KeePassHandler)

	kph.filename = ".env.kdbx"

	file, err := os.Open(kph.filename)
	if err != nil {
		if os.IsNotExist(err) {
			kph.createDatabase()
		}
	} else {
		defer file.Close()
		kph.file = file
		kph.unlockDB()
	}

	return kph
}

func GetKeyPassHandler() *KeePassHandler {
	return kph
}

func (kph *KeePassHandler) unlockDB() {
	// filename := kph.filename
	// file, err := os.Open(filename)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	kph.db = gokeepasslib.NewDatabase()

	// try to unlock the file with password input
	for {
		password := ReadPassword("Please enter your DB password")
		kph.db.Credentials = gokeepasslib.NewPasswordCredentials(password)
		err := gokeepasslib.NewDecoder(kph.file).Decode(kph.db)
		if err != nil {
			log.Fatal("could not decode the file")
		}

		err = kph.db.UnlockProtectedEntries()
		if err != nil {
			log.Print("could not open the db file")
		} else {
			log.Print("db file decrypted successfully")
			break
		}
	}
}

func (kph *KeePassHandler) lockDB() {
	db := kph.db
	db.LockProtectedEntries()

	file, err := os.Create(kph.filename)
	if err != nil {
		log.Print(err)
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
		pwd = ReadPassword("Choose a password")
		vpwd := ReadPassword("Verify password")
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

func (kph *KeePassHandler) AddRecords() {
	db := kph.db
	entries := &db.Content.Root.Groups[0].Entries

	for {
		input := ReadInput("add a new environment variable (KEY=VALUE). press enter to quit")
		// exit from loop
		if len(input) == 0 {
			break
		}

		keyVal := strings.Split(input, "=")
		if len(keyVal) != 2 {
			log.Println("invalid input!")
			continue
		}

		entry := gokeepasslib.NewEntry()
		entry.Values = append(entry.Values, mkProtectedValue(keyVal[0], keyVal[1]))
		*entries = append(*entries, entry)
	}
}
