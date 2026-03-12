package cli

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshckidd/gm_tools/internal/config"
	"github.com/joshckidd/gm_tools/internal/database"
	"github.com/joshckidd/gm_tools/internal/requests"
	"github.com/joshckidd/gm_tools/internal/rolls"
)

type State struct {
	Cfg *config.CliConfig
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
	tot, err := requests.GenerateRoll(s.Cfg, cmd.Args[0])
	if err != nil {
		return err
	}

	printRoll(tot)

	return nil
}

func HandlerList(s *State, cmd Command) error {
	switch cmd.Args[0] {
	case "rolls":
		err := listRolls(s.Cfg)
		if err != nil {
			return err
		}
	case "types":
		err := listTypes(s.Cfg)
		if err != nil {
			return err
		}
	case "custom_fields":
		err := listCustomFields(s.Cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("Login command requires two arguments: username, password")
	}

	err := requests.LoginUser(s.Cfg, cmd.Args[0], cmd.Args[1])
	if err != nil {
		return err
	}

	err = s.Cfg.SetToken()
	if err != nil {
		return err
	}

	return nil
}

func printRoll(tot rolls.RollTotalResult) {
	fmt.Printf("total: %d\n", tot.TotalResult)

	for i, rs := range tot.IndividualResults {
		fmt.Printf(" - Roll set %d: %d\n", i, rs.Result)

		for j, r := range rs.IndividualRolls {
			fmt.Printf(" --- Roll %d: %d\n", j, r)
		}
	}
}

func listRolls(cfg *config.CliConfig) error {
	rolls, err := requests.GetRecords[rolls.RollTotalResult](cfg, "rolls")
	if err != nil {
		return err
	}

	for i := range rolls {
		fmt.Printf("Roll %d - %s:\n", i+1, rolls[i].RollString)
		printRoll(rolls[i])
	}

	return nil
}

func listTypes(cfg *config.CliConfig) error {
	types, err := requests.GetRecords[database.Type](cfg, "types")
	if err != nil {
		return err
	}

	for i := range types {
		fmt.Printf("%d. ID: %s\n   Name: %s\n", i+1, types[i].ID, types[i].TypeName)
	}

	return nil
}

func listCustomFields(cfg *config.CliConfig) error {
	fields, err := requests.GetRecords[database.CustomField](cfg, "custom_fields")
	if err != nil {
		return err
	}

	typeMap, err := getTypeMap(cfg)
	if err != nil {
		return err
	}

	for i := range fields {
		fmt.Printf("%d. ID: %s\n   Name: %s\n   Type: %s\n",
			i+1,
			fields[i].ID, fields[i].CustomFieldName,
			typeMap[fields[i].TypeID])
	}

	return nil
}

func getTypeMap(cfg *config.CliConfig) (map[uuid.UUID]string, error) {
	types, err := requests.GetRecords[database.Type](cfg, "types")
	if err != nil {
		return map[uuid.UUID]string{}, err
	}

	typeMap := make(map[uuid.UUID]string, len(types))
	for i := range types {
		typeMap[types[i].ID] = types[i].TypeName
	}

	return typeMap, nil
}
