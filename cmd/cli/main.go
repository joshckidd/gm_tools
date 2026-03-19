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
	cliCommands.Register("list", cli.HandlerList)
	cliCommands.Register("generate", cli.HandlerGenerate)
	cliCommands.Register("delete", cli.HandlerDelete)
	cliCommands.Register("create", cli.HandlerCreate)
	cliCommands.Register("update", cli.HandlerUpdate)
	cliCommands.Register("load", cli.HandlerLoad)
	cliCommands.Register("export", cli.HandlerExport)

	args := os.Args
	if len(args) < 3 {
		fmt.Println("Please provide a command and at least one argument.")
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
