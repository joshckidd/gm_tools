package main

import (
	"fmt"
	"os"

	"github.com/joshckidd/gm_tools/internal/cli"
	"github.com/joshckidd/gm_tools/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	var s cli.State
	s.Cfg = &c

	cliCommands := cli.Commands{
		CommandMap: map[string]func(*cli.State, cli.Command) error{},
	}

	cliCommands.Register("roll", cli.HandlerRoll)
	cliCommands.Register("login", cli.HandlerLogin)

	args := os.Args
	if len(args) < 3 {
		fmt.Println("I require an argument!")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cliCommands.Run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
