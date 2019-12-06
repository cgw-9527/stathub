package main

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestHttp(t *testing.T) {
	cmd := exec.Command("du", "-sh", "ulord_1_1_86.tgz")
	size, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	s := strings.Split(string(size), "M")
	log.Println(s[0])
}
