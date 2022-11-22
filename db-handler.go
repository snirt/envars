package main

import "os"

type DBHandler interface {
	ListDB()
	CreateDB(name string, password string, folderPath string)
	GetFile() *os.File
	DeleteDB()
}