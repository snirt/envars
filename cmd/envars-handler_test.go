package cmd

import (
	// "os"
	"os"
	"testing"
)

type MockPasswordReader struct{}

func (m *MockPasswordReader) ReadPassword() string {
	return "mocked_password"
}

type MockReader struct{}

func (m *MockReader) ReadPassword(s string) string {
	return "mocked_password"
}

var i int = 0
var inputArr = []string{"key=value", ""}

func (m *MockReader) ReadInput(s string) string {
	returnVal := inputArr[i]
	i++
	return returnVal
}

func TestNewDatabase(t *testing.T) {
	defer cleanup()
	os.Setenv("ENVARS_PWD", "mocked_password")
	var reader Reader = &MockReader{}
	var envars EnvarsInterface = New(reader)
	envars.AddVariables(reader)
	envars.lockDB()
}

func cleanup() {
	os.Unsetenv("ENVARS_PWD")
	os.Remove(".env.kdbx")
}
