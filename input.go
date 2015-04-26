package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func confirmDefault(question string, defaultAnswer bool) bool {
	var prompt string
	if defaultAnswer {
		prompt = "[Y/n]"
	} else {
		prompt = "[y/N]"
	}
	fmt.Print(question, " ", prompt, " ")

	r := bufio.NewReader(os.Stdin)
	answer, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	answer = strings.TrimSpace(answer)

	switch string(answer) {
	case "y", "Y", "yes", "Yes", "YES":
		return true
	case "n", "N", "no", "No", "NO":
		return false
	default:
		return defaultAnswer
	}
}
