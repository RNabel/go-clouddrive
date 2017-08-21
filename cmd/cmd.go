package cmd

import (
	"bufio"
	"fmt"
	"go-clouddrive/clouddrive"
	"os"
	"strings"
	"sync"
)

var handlers = map[string]func(string, *clouddrive.CDrive){
	"ls": func(path string, cd *clouddrive.CDrive) {
		split := strings.Split(path, " ")
		if len(split) == 2 {
			// TODO finish this.

		} else {
			fmt.Println("wrong number of args. usage: ls <path>")
		}
	},
}

func Start(cd *clouddrive.CDrive, wg *sync.WaitGroup) {
	defer wg.Done()

	in := bufio.NewReader(os.Stdin)
	fmt.Println("Ready for terminal input")
	for {
		input, err := in.ReadString('\n')
		if err != nil {
			fmt.Println("Error ocurred when reading the string")
			continue
		}
		trimmed := strings.TrimRight(input, "\r\n")

		split := strings.Split(trimmed, " ")
		if len(split) > 0 {
			if split[0] == "exit" {
				break
			} else if f, ok := handlers[split[0]]; ok {
				f(trimmed, cd)
			} else {
				fmt.Printf("'%s' is not a valid command.\n", trimmed)
			}

		}
	}
}
