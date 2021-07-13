package main

import (
	"fmt"
	"os/exec"
	"time"
)

func launch(task string) {
	fmt.Printf("%v: Launching: %s \n", time.Now().Format("2006-01-02 15:04:05"), task)
	cmd := exec.Command(task)
	if err := cmd.Start(); err != nil {
		fmt.Println(err.Error())
	}
}
