package cli

import (
	"errors"
	"fmt"

	"github.com/joshckidd/gm_tools/internal/config"
	"github.com/joshckidd/gm_tools/internal/requests"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	commandFunction, ok := c.CommandMap[cmd.Name]
	if !ok {
		return errors.New("Command not found.")
	}

	return commandFunction(s, cmd)
}

func HandlerRoll(s *State, cmd Command) error {
	tot, err := requests.GenerateRoll(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("total: %d\n", tot.TotalResult)

	for i, rs := range tot.IndividualResults {
		fmt.Printf(" - Roll set %d: %d\n", i, rs.Result)

		for j, r := range rs.IndividualRolls {
			fmt.Printf(" --- Roll %d: %d\n", j, r)
		}
	}
	return nil
}
