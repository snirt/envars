package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/tobischo/gokeepasslib/v3"
)

type KeyPassHandler struct {
	conf *Config
	File *os.File
}

var kph *KeyPassHandler

func init() {
	kph = New()
}

func New() *KeyPassHandler {
	kph := new(KeyPassHandler)
	kph.conf.LoadConfigFile()
}

func GetKeyPassHandler() *KeyPassHandler {
	config := Config{}
	conf, err := config.LoadConfigFile()
	if err != nil {
		fmt.Println("could not load config file, trying to generate a new one...")
		conf, err = CreateConfigFile()
		if err != nil {
			log.Fatal(err)
		}
	}
	kph := &KeyPassHandler{
		conf:  conf,
	}
	return kph
}

func (kph *KeyPassHandler) ListDB() {
	fmt.Println("nlanlsdnf")
}

func (kph *KeyPassHandler) CreateDB(name string, masterPassword string, folderPath string) {
	// try to get path from config
	kph.conf.DBPath = viper.GetViper().GetString(K_KEEPASS_DB_PATH)

	//TODO check if the file from config exists
	if name == "" {
		name = ReadInput("Please enter the Project name")
	}
	if (masterPassword == "") {
		masterPassword = ReadInput("Create a master password, please make sure its strong enough")
	}
	if (folderPath == "") {
		folderPath = ReadInput("Enter a valid path for the DB")
	}

	kph.conf.DBPath = filepath.Join(folderPath, name+".kdbx")


	if _, err := os.Stat(kph.conf.DBPath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(folderPath, os.ModePerm)
		file, err := os.Create(kph.conf.DBPath)
		if err != nil {
			panic(err)
		}
		kph.File = file
	} else {
		file, err := os.OpenFile(kph.conf.DBPath, os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
		
		kph.File = file
	}

	// init db with root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root group"

	// create db to contain the root group
	db := &gokeepasslib.Database{
		Header: gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(masterPassword),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	kph.lockDB(db)
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


func (kph *KeyPassHandler) lockDB(db *gokeepasslib.Database) {
	db.LockProtectedEntries()

	keepassEncoder := gokeepasslib.NewEncoder(kph.File)
	if err := keepassEncoder.Encode(db); err != nil {
		panic(err)
	}

	kph.File.Close()
	log.Printf("db file encoded and closed")
}


func (kph *KeyPassHandler) unlockDB() *gokeepasslib.Database{
	kph.File, _ = os.Open(kph.conf.DBPath)

	db := gokeepasslib.NewDatabase()

	dbPwd := kph.getDbPwdFromEnv()
	if (dbPwd == "") {
		dbPwd = ReadInput("Please enter DB password for " + kph.conf.CurrentDB)
	}
    db.Credentials = gokeepasslib.NewPasswordCredentials(dbPwd)
    _ = gokeepasslib.NewDecoder(kph.File).Decode(db)

	db.UnlockProtectedEntries()
	return db
}


func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}


// write a new record to db
func (kph *KeyPassHandler) AddRecord() {
	db := kph.unlockDB()
	defer kph.lockDB(db)

	input := " = "
	entry := gokeepasslib.NewEntry()
	for len(strings.Split(input, "=")) == 2 {
		input := ReadInput("please enter 'KEY=VALUE' (empty to exit)")
		//TODO add validation
		inputSlice := strings.Split(input, "=");
		key := inputSlice[0]
		value := inputSlice[1]
	
		entry.Values = append(entry.Values, mkValue(key, value))
	}
	rootGroupEntries := db.Content.Root.Groups[0].Entries
	rootGroupEntries = append(rootGroupEntries, entry)
	db.Content.Root.Groups[0].Entries = rootGroupEntries
}


func (kph *KeyPassHandler) setDbPwdToEnv(dbPwdVal string) error {
	dbPwdKey := strings.ToUpper(kph.conf.CurrentDB) + "_" + strings.ToUpper(kph.conf.CurrentEnv)
	return  os.Setenv(dbPwdKey, dbPwdVal)
}


func (kph *KeyPassHandler) getDbPwdFromEnv() string {
	dbPwdKey := strings.ToUpper(kph.conf.CurrentDB) + "_" + strings.ToUpper(kph.conf.CurrentEnv)
	return os.Getenv(dbPwdKey)
}


func (kph *KeyPassHandler) saveToDB(subGroup []gokeepasslib.Group) error {
	db := kph.unlockDB()
	defer kph.lockDB(db)

	// get password from environment variable if exists
	dbPwd := kph.getDbPwdFromEnv()
	if (dbPwd == "") {
		dbPwd = ReadInput("Please enter DB password for " + kph.conf.CurrentDB)
	}

	db.Credentials = gokeepasslib.NewPasswordCredentials(dbPwd)
	_ = gokeepasslib.NewDecoder(kph.File).Decode(db)

	err := db.UnlockProtectedEntries()
	if err != nil {
		return err
	}

	// root entries represents environments, sub entries represents... //TODO continue here!!!
	return nil
}