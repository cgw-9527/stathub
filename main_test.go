package main

import (
	"log"
	"os/exec"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	cmd := exec.Command("ulord-cli", "stop")
	cmd.CombinedOutput()
	time.Sleep(60 * time.Second)
check:
	str := "ps aux|grep ulordd|grep -v grep"
	cmd = exec.Command("sh", "-c", str)
	out1, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	if string(out1) != "" {
		time.Sleep(60 * time.Second)
		goto check
	}
	s := "nohup ulordd &"
	cmd = exec.Command("sh", "-c", s)
	cmd.Run()
}
