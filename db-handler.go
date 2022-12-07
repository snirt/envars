package main

import "os"

type DBHandler interface {
	ListDB()
	CreateDB(name string, masterPassword string, folderPath string)
	GetFile() *os.File
	DeleteDB()
	AddRecord() error
}