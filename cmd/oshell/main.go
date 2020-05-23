package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/andrebq/organic/cell"
	"github.com/andrebq/organic/lua/binding"
	"github.com/andrebq/organic/medium"
	"github.com/chzyer/readline"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

var (
	agarAddr = flag.String("agar", "localhost:6739", "Address of our Agar instances (aka Redis)")
	cellname = flag.String("cell", "", "Name of our cell, if empty a random value is used")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agar, err := medium.NewAgar(ctx, "")
	if err != nil {
		panic(err)
	}
	c := cell.Grow(computeName(*cellname), agar)

	l := lua.NewState()
	l.PreloadModule("organic", binding.Loader(c))
	runShell(l, os.Stdin, os.Stdout, os.Stderr)
}

func computeName(input string) string {
	if input != "" {
		return input
	}
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return "oshell-" + base64.URLEncoding.EncodeToString(buf)
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
			rl.Clean()
			if code.Len() > 0 {
				rl.SetPrompt(expectingInput)
				fmt.Fprintln(rl, "...discarded...")
				fmt.Fprintln(rl, "...type ctrl+c again to exit...")
				code.Reset()
				rl.Refresh()
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
