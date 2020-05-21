package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chzyer/readline"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
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
	defaultPrompt := "> "
	expectingInput := "  "
	code := bytes.Buffer{}
	for {
		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if code.Len() > 0 {
				rl.SetPrompt(expectingInput)
				code.Reset()
				continue
			}
			log.Print("Goodbye!")
			return
		} else if err != nil {
			log.Printf("[Input error]>\t%v", err)
		}

		if line == "exit" {
			log.Print("Goodbye")
			return
		}
		code.WriteString(line)

		_, err = parse.Parse(bytes.NewBuffer(code.Bytes()), "_shell")
		if err != nil {
			// invalid code
			// just move to the next line
			rl.SetPrompt(expectingInput)
			continue
		} else {
			rl.SetPrompt(defaultPrompt)
		}
		err = l.DoString(code.String())
		if err != nil {
			fmt.Fprintf(rl, "[Error]\t%v", err)
			fmt.Fprintln(rl)
		}
		code.Reset()
	}
}
