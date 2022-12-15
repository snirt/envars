package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ReadInput(msg string) string {
	fmt.Println(msg)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}

func ReadPassword(msg string) string {
	fmt.Println(msg)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Something is wrong with password")
		return ""
	}
	return string(bytePassword)
}

type Color string

func Print(s string, c Color) {
	fmt.Fprintln(os.Stdout, c, s, COLOR_NONE)
}
