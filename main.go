package main

import (
	"fmt"
	"github.com/cubxxw/sealer-runtime/apply"
)

func main() {
	processor := &apply.Processor{}
	fmt.Sprintln("ls -al")
	fmt.Printf("%#v\n", processor)

	if err := processor.ApplyClusterFile(); err != nil {
		fmt.Println(" handler a error")
	}
}
