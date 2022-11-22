package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// "github.com/tobischo/gokeepasslib/v3"
)


func CreateLocalDB(path string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Sprintf("Creating a new database: %s", path)
	fmt.Println("Enter a password")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input.")
	}
	input = strings.TrimSuffix(input, "\n")
	if len(input) == 0 {
		fmt.Println("Choose a better password")
	}

}