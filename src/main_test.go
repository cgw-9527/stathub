package main

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestHttp(t *testing.T) {
	str := "ps aux|grep ulordd|grep -v grep"
	cmd := exec.Command("sh", "-c", str)
	out1, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("checkstatus ulordd 2", err)
	}
	s := strings.Split(string(out1), "\n")
	s = strings.Split(s[1], "   ")
	s = strings.Split(s[1], " ")
	log.Println(s[2])
}
