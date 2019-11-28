package main

import (
	"log"
	"os/exec"
	"testing"
)

func TestHttp(t *testing.T) {
	cmd := exec.Command("ps", "aux|", "grep", "ulordd", "|grep", "-v", "grep")
	out1, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	if out1 != nil {
		log.Println(out1)
	} else {
		log.Println("---------")
	}
}
