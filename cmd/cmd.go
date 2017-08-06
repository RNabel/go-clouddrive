package cmd

import (
	"bufio"
	"os"
	"fmt"
	"strings"
)

var handlers = map[string]func(string){
	"ls": func(path string) {
		split := strings.Split(path, " ")
		if len(split) == 2 {
			//fpath := split[1]

		} else {
			fmt.Println("wrong number of args. usage: ls <path>")
		}
	},
}

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ready for terminal input")
	for {
		in := scanner.Text()
		split := strings.Split(in, " ")
		if len(split) > 0 {
			if split[0] == "exit" {
				break
			} else if f, ok := handlers[split[0]]; ok {
				f(in)
			} else {
				fmt.Printf("'%s' is not a valid command.", )
			}

		}
		fmt.Println(in)
	}
}
