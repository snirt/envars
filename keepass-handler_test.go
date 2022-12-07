package main

import (
	"os"
	"testing"
)

var Handler *KeyPassHandler = &KeyPassHandler{}
// do not change the following 2 tests order!
func TestCreateKeePassDB_ForNonExistFile(t *testing.T) {
	var Handler *DBHandler = GetKeyPassHandler()
	Handler.CreateDB("my_test_db", "password", "test")
	if (Handler.GetFile() == nil) {
		t.Error("Handler.File should be populated after creaating DB")
	}
}

func TestCreateKeePassDB_ForExistFile(t *testing.T) {
	// var Handler DBHandler = &KeyPassHandler{}
	Handler.CreateDB("my_test_db", "password", "test")
	if (Handler.GetFile() == nil) {
		t.Error("Handler.File should be populated after creaating DB")
	}

	// clean test file
	err := os.Remove("./test/my_test_db.kdbx")
	if err != nil {
		t.Fatal("couldn't remove test db file")
	}
}

func TestDeleteDB(t *testing.T) {
	Handler.CreateDB("my_test_db", "password", "test")
	Handler.DeleteDB()
}