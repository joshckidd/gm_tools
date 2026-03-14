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
	tot, err := requests.CallApi[rolls.RollTotalResult](
		s.Cfg,
		"rolls",
		"POST",
		map[string]string{"roll": cmd.Args[0]},
	)
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
	case "items":
		err := listItems(s.Cfg)
		if err != nil {
			return err
		}
	case "instances":
		err := listInstances(s.Cfg)
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

func HandlerGenerate(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("Login command requires two arguments: item type, number")
	}

	instances, err := requests.CallApi[[]map[string]string](
		s.Cfg,
		"instances",
		"POST",
		map[string]string{
			"type":   cmd.Args[0],
			"number": cmd.Args[1],
		},
	)
	if err != nil {
		return err
	}

	printItems(s.Cfg, instances)

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
	rolls, err := requests.CallApi[[]rolls.RollTotalResult](cfg, "rolls", "GET", "")
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
	types, err := requests.CallApi[[]database.Type](cfg, "types", "GET", "")
	if err != nil {
		return err
	}

	for i := range types {
		fmt.Printf("%d. ID: %s\n   Name: %s\n", i+1, types[i].ID, types[i].TypeName)
	}

	return nil
}

func listCustomFields(cfg *config.CliConfig) error {
	fields, err := requests.CallApi[[]database.CustomField](cfg, "custom_fields", "GET", "")
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
			fields[i].ID,
			fields[i].CustomFieldName,
			typeMap[fields[i].TypeID],
		)
	}

	return nil
}

func listItems(cfg *config.CliConfig) error {
	items, err := requests.CallApi[[]map[string]string](cfg, "items", "GET", "")
	if err != nil {
		return err
	}

	printItems(cfg, items)

	return nil
}

func listInstances(cfg *config.CliConfig) error {
	instances, err := requests.CallApi[[]map[string]string](cfg, "instances", "GET", "")
	if err != nil {
		return err
	}

	printItems(cfg, instances)

	return nil
}

func printItems(cfg *config.CliConfig, items []map[string]string) {
	typeMap, _ := getTypeMap(cfg)

	for i := range items {
		typeId, _ := uuid.Parse(items[i]["type"])
		fmt.Printf("%d. ID: %s\n   Type: %s\n   Name: %s\n   Description: %s\n",
			i+1,
			items[i]["id"],
			typeMap[typeId],
			items[i]["name"],
			items[i]["description"],
		)
		for k, v := range items[i] {
			if k != "id" && k != "type" && k != "name" && k != "description" && k != "created_at" && k != "updated_at" && k != "username" {
				fmt.Printf("   %s: %s\n", k, v)
			}
		}
	}
}

func getTypeMap(cfg *config.CliConfig) (map[uuid.UUID]string, error) {
	types, err := requests.CallApi[[]database.Type](cfg, "types", "GET", "")
	if err != nil {
		return map[uuid.UUID]string{}, err
	}

	typeMap := make(map[uuid.UUID]string, len(types))
	for i := range types {
		typeMap[types[i].ID] = types[i].TypeName
	}

	return typeMap, nil
}
