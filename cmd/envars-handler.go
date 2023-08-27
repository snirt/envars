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

type EnvarsInterface interface {
	AddVariables(Reader)
	RemoveVariables([]string)
	ListVariables(bool)

	createDatabase(Reader)
	lockDB() error
	unlockDB() error
}
type Envars struct {
	filename string
	db       *gokeepasslib.Database
	file     *os.File
	password string
}

// var envars *Envars

// creates Envars object with unlocked db. MAKE SURE YOU CLOSE IT!
func New(reader Reader) EnvarsInterface {
	envars := new(Envars)
	envars.filename = DB_FILE_NAME

	file, err := os.Open(envars.filename)
	if err != nil {
		if os.IsNotExist(err) {
			envars.createDatabase(reader)
		} else {
			Print(err.Error(), COLOR_RED)
		}
	} else {
		defer file.Close()
		envars.file = file
		if err := envars.unlockDB(); err != nil {
			Print(err.Error(), COLOR_RED)
		}
	}
	return envars
}

// unlockDB is a function that attempts to unlock the database file using the provided password.
//
// It takes no parameters.
// It returns an error.
func (envars *Envars) unlockDB() error {
	newDB := gokeepasslib.NewDatabase()

	// try to unlock the file with password input
	for {
		if os.Getenv(ENVARS_PWD) != "" {
			envars.password = os.Getenv(ENVARS_PWD)
		} else {
			return fmt.Errorf("%s", ERR_PWD_ERR_MSG)
		}
		newDB.Credentials = gokeepasslib.NewPasswordCredentials(envars.password)
		if err := gokeepasslib.NewDecoder(envars.file).Decode(newDB); err != nil {
			return fmt.Errorf(ERR_DECODE)
		}
		if err := newDB.UnlockProtectedEntries(); err != nil {
			return fmt.Errorf(ERR_FILE_OPEN)
		} else {
			break
		}
	}
	envars.db = newDB
	return nil
}

// lockDB locks the database, creates a file, and encodes the database into the file.
//
// No parameters.
// No return values.
func (envars *Envars) lockDB() error {
	db := envars.db
	err := db.LockProtectedEntries()
	if err != nil {
		Print(ERR_LOCK_DB, COLOR_RED)
		return fmt.Errorf("could not lock db")
	}

	file, err := os.Create(envars.filename)
	if err != nil {
		Print(err.Error(), COLOR_RED)
		return fmt.Errorf("could not open the db file to save the data")
	}
	defer file.Close()

	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(db); err != nil {
		Print(err.Error(), COLOR_RED)
		return fmt.Errorf("could not encode the db")
	}
	return nil
}

func (envars *Envars) createDatabase(reader Reader) {
	fmt.Println("Creating a new database...")
	var pwd string
	for {
		pwd = reader.ReadPassword("Choose a password")
		vpwd := reader.ReadPassword("Verify password")
		if pwd == vpwd {
			break
		}
	}

	file, err := os.Create(envars.filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// init db with root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root group"

	// create db to contain the root group
	envars.db = &gokeepasslib.Database{
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

func (envars *Envars) AddVariables(reader Reader) {
	db := envars.db
	entries := db.Content.Root.Groups[0].Entries

	// init the variables map
	variablesMap := make(map[string]gokeepasslib.Entry)
	for _, e := range entries {
		variablesMap[e.Values[0].Key] = e
	}

	for {
		input := reader.ReadInput("add a new environment variable (KEY=VALUE). press enter to quit")
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
			input := reader.ReadInput(keyVal[0] + " already exists. do you want to update this variable value? (y/N)")
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

func (envars *Envars) KeepPwdInSession() {
	fmt.Println(envars.password)
}

func (envars *Envars) ListVariables(exportable bool) {
	var prefix string
	if exportable {
		if runtime.GOOS == "windows" {
			prefix = "set"
		} else {
			prefix = "export"
		}
	}

	db := envars.db
	entries := &db.Content.Root.Groups[0].Entries
	for _, enrty := range *entries {
		val := enrty.Values[0]
		fmt.Printf("%s %s='%s'\n", prefix, val.Key, val.Value.Content)
	}
}

func (envars *Envars) RemoveVariables(variables []string) {
	db := envars.db
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
