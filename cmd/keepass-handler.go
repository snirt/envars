package cmd

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

type KeePassHandler struct {
	filename string
	db       *gokeepasslib.Database
	file     *os.File
	password string
}

var kph *KeePassHandler

// creates KeePassHandler object with unlocked db. MAKE SURE YOU CLOSE IT!
func New() *KeePassHandler {
	kph := new(KeePassHandler)

	kph.filename = DB_FILE_NAME

	file, err := os.Open(kph.filename)
	if err != nil {
		if os.IsNotExist(err) {
			kph.createDatabase()
		}
	} else {
		defer file.Close()
		kph.file = file
		err := kph.unlockDB()
		if err != nil {
			Print(err.Error(), COLOR_RED)
			os.Exit(1)
		}
	}
	return kph
}

func GetKeyPassHandler() *KeePassHandler {
	return kph
}

func (kph *KeePassHandler) unlockDB() error {
	newDB := gokeepasslib.NewDatabase()

	// try to unlock the file with password input
	for {
		if os.Getenv(ENVARS_PWD) != "" {
			kph.password = os.Getenv(ENVARS_PWD)
		} else if passwordArg != "" {
			kph.password = passwordArg
		} else {
			return fmt.Errorf("%s", ERR_PWD_ERR_MSG)
		}
		newDB.Credentials = gokeepasslib.NewPasswordCredentials(kph.password)
		err := gokeepasslib.NewDecoder(kph.file).Decode(newDB)
		if err != nil {
			return fmt.Errorf(ERR_DECODE)
		}

		err = newDB.UnlockProtectedEntries()
		if err != nil {
			return fmt.Errorf(ERR_FILE_OPEN)
		} else {
			break
		}
	}
	kph.db = newDB
	return nil
}

func (kph *KeePassHandler) lockDB() {
	db := kph.db
	err := db.LockProtectedEntries()
	if err != nil {
		Print(ERR_LOCK_DB, COLOR_RED)
	}

	file, err := os.Create(kph.filename)
	if err != nil {
		Print(err.Error(), COLOR_RED)
	}
	defer file.Close()

	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(db); err != nil {
		Print(err.Error(), COLOR_RED)
	}
}

func (kph *KeePassHandler) createDatabase() {
	fmt.Println("Creating a new database...")
	var pwd string
	for {
		pwd = ReadPassword("Choose a password")
		vpwd := ReadPassword("Verify password")
		if pwd == vpwd {
			break
		}
	}

	file, err := os.Create(kph.filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

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

func isValidEnvVar(s string) bool {
	pattern := regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*=[0-9A-Za-z_-]*$")
	return pattern.MatchString(s)
}

func (kph *KeePassHandler) AddVariables() {
	db := kph.db
	entries := db.Content.Root.Groups[0].Entries

	// init the variables map
	variablesMap := make(map[string]gokeepasslib.Entry)
	for _, e := range entries {
		variablesMap[e.Values[0].Key] = e
	}

	for {
		input := ReadInput("add a new environment variable (KEY=VALUE). press enter to quit")
		// exit from loop
		if len(input) == 0 {
			break
		}

		if !isValidEnvVar(input) {
			Print("invalid input!", COLOR_RED)
			continue
		}

		keyVal := strings.Split(input, "=")

		// check if the variable key already exists
		val, exists := variablesMap[keyVal[0]]
		if exists {
			// ask if to update or not
			input := ReadInput(keyVal[0] + " already exists. do you want to update this variable value? (y/N)")
			if strings.ToLower(input) == "y" {
				val.Values[0].Value.Content = keyVal[1]
			}
		} else {
			// add the variable to the map
			entry := gokeepasslib.NewEntry()
			entry.Values = append(entry.Values, mkValue(keyVal[0], keyVal[1]))
			variablesMap[keyVal[0]] = entry
		}
	}

	// write the variablesMap to db
	entries = []gokeepasslib.Entry{}
	for _, val := range variablesMap {
		entries = append(entries, val)
	}
	db.Content.Root.Groups[0].Entries = entries

}

func (kph *KeePassHandler) KeepPwdInSession() {
	fmt.Println(kph.password)
}

func (kph *KeePassHandler) ListVariables(exportable bool) {
	var prefix string
	if exportable {
		if runtime.GOOS == "windows" {
			prefix = "set"
		} else {
			prefix = "export"
		}
	}

	db := kph.db
	entries := &db.Content.Root.Groups[0].Entries
	for _, enrty := range *entries {
		val := enrty.Values[0]
		fmt.Printf("%s %s='%s'\n", prefix, val.Key, val.Value.Content)
	}
}

func (kph *KeePassHandler) RemoveVariables(variables []string) {
	db := kph.db
	entries := db.Content.Root.Groups[0].Entries

	for _, v := range variables {
		found := false
		for i, e := range entries {
			if e.Values[0].Key == v {
				if len(entries) > 1 {
					entries = append(entries[:i], entries[i+1:]...)
				} else {
					entries = []gokeepasslib.Entry{}
				}
				found = true
			}
		}
		if found {
			Print(v+" removed", COLOR_RED)
		} else {
			Print(v+" not found", COLOR_RED)
		}
	}
	db.Content.Root.Groups[0].Entries = entries
}
