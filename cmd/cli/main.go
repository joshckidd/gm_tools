// the main package for the cli tool that will interface with the rest api

package main

import (
	"fmt"
	"os"

	"github.com/joshckidd/gm_tools/internal/cli"
	"github.com/joshckidd/gm_tools/internal/config"
)

func main() {
	// read the local configuration file to get the API URL and bearer token for the user
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	var s cli.State
	s.Cfg = &c

	cliCommands := cli.Commands{
		CommandMap: map[string]func(*cli.State, cli.Command) error{},
	}

	// register commands available fot the cli
	cliCommands.Register("roll", cli.HandlerRoll)
	cliCommands.Register("login", cli.HandlerLogin)
	cliCommands.Register("list", cli.HandlerList)
	cliCommands.Register("generate", cli.HandlerGenerate)
	cliCommands.Register("delete", cli.HandlerDelete)
	cliCommands.Register("create", cli.HandlerCreate)
	cliCommands.Register("update", cli.HandlerUpdate)
	cliCommands.Register("load", cli.HandlerLoad)
	cliCommands.Register("export", cli.HandlerExport)

	// check to make sure at least 1 argument is provided
	// individual handlers will further validate arguments
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Please provide a command and at least one argument.")
		os.Exit(1)
	}

	// use the Command struct to store the command and all arguments
	// to be passed to handler functions
	cmd := cli.Command{
		Name: args[1],
		Args: args[2:],
	}

	// execute the command
	err = cliCommands.Run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
