package main

import (
	"log"
	"os/exec"
	"testing"
)

func TestHttp(t *testing.T) {
	cmd := exec.Command("pidof", "ulordd")
	out1, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("checkstatus ulordd 2", err)
	}
	log.Println(string(out1) + "")
}
