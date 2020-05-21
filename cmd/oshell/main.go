package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chzyer/readline"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	l := lua.NewState()
	runShell(l, os.Stdin, os.Stdout, os.Stderr)
}

func runShell(l *lua.LState, in io.ReadCloser, sout, serr io.WriteCloser) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
	})
	defer rl.Close()

	if err != nil {
		fmt.Fprintf(serr, "Error!: %v", err)
		fmt.Fprintln(serr)
	}
	log.SetOutput(rl.Stderr())
	for {
		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			log.Print("Goodbye!")
			return
		} else if err != nil {
			log.Printf("[Input error]>\t%v", err)
		}

		if line == "exit" {
			log.Print("Goodbye")
			return
		}
		err = l.DoString(line)
		if err != nil {
			log.Printf("[Error]> \t %v", err)
		} else {
			top := l.GetTop()
			if top != 0 {
				fmt.Fprintf(rl, "< %v", l.ToString(top))
				fmt.Fprintln(rl)
			} else {
				fmt.Fprintln(rl)
			}
		}
	}
}
