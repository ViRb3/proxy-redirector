package main

import (
	"strings"
	"fmt"
)

func splitLines(content string) []string {
	return strings.Split(strings.Replace(content, "\r\n", "\n", -1), "\n")
}

func printError(err error) {
	fmt.Print("ERROR: ")
	fmt.Println(err)
}
