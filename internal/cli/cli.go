package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

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
	var err error
	var ids []string

	if len(cmd.Args) > 1 {
		ids = cmd.Args[1:]
	} else {
		ids = []string{}
	}

	switch cmd.Args[0] {
	case "rolls":
		err = listRecords(s.Cfg, cmd.Args[0], []string{}, printRolls)
	case "types":
		err = listRecords(s.Cfg, cmd.Args[0], ids, printTypes)
	case "custom_fields":
		err = listRecords(s.Cfg, cmd.Args[0], ids, printCustomFields)
	case "items", "instances":
		err = listRecords(s.Cfg, cmd.Args[0], ids, printItems)
	default:
		return errors.New("Invalid table provided for list.")
	}
	return err
}

func HandlerDelete(s *State, cmd Command) error {
	var ids []string

	if len(cmd.Args) < 2 {
		return errors.New("Delete command requires at least two arguments: table, one or more ids")
	}

	ids = cmd.Args[1:]

	switch cmd.Args[0] {
	case "types", "custom_fields", "items":
		for i := range ids {
			endpoint := fmt.Sprintf("%s/%s", cmd.Args[0], ids[i])
			_, err := requests.CallApi[string](s.Cfg, endpoint, "DELETE", "")
			if err != nil {
				return err
			}
			fmt.Printf("Deleted %s: %s\n", cmd.Args[0], ids[i])
		}

	default:
		return errors.New("Invalid table provided for delete.")
	}
	return nil
}

func HandlerCreate(s *State, cmd Command) error {
	switch cmd.Args[0] {
	case "types":
		if len(cmd.Args) < 2 {
			return errors.New("Create command for types requires two arguments: table, type name")
		}

		t, err := requests.CallApi[database.Type](s.Cfg, cmd.Args[0], "POST", map[string]string{
			"type": cmd.Args[1],
		})

		if err == nil {
			printTypes(s.Cfg, []database.Type{t})
		}
		return err
	case "custom_fields":
		if len(cmd.Args) < 4 {
			return errors.New("Create command for custom fields requires four arguments: table, item type, custom field name, custom field type")
		}
		cf, err := requests.CallApi[database.CustomField](s.Cfg, cmd.Args[0], "POST", map[string]string{
			"type":       cmd.Args[1],
			"field_name": cmd.Args[2],
			"field_type": cmd.Args[3],
		})

		if err == nil {
			printCustomFields(s.Cfg, []database.CustomField{cf})
		}
		return err
	default:
		return errors.New("Invalid table provided for create.")
	}
}

func HandlerUpdate(s *State, cmd Command) error {
	switch cmd.Args[0] {
	case "types":
		if len(cmd.Args) < 3 {
			return errors.New("Update command for types requires three arguments: table, id, type name")
		}
		endpoint := fmt.Sprintf("%s/%s", cmd.Args[0], cmd.Args[1])
		t, err := requests.CallApi[database.Type](s.Cfg, endpoint, "PUT", map[string]string{
			"type": cmd.Args[2],
		})

		if err == nil {
			printTypes(s.Cfg, []database.Type{t})
		}
		return err
	case "custom_fields":
		if len(cmd.Args) < 5 {
			return errors.New("Update command for custom fields requires five arguments: table, id, item type, custom field name, custom field type")
		}
		endpoint := fmt.Sprintf("%s/%s", cmd.Args[0], cmd.Args[1])
		cf, err := requests.CallApi[database.CustomField](s.Cfg, endpoint, "PUT", map[string]string{
			"type":       cmd.Args[2],
			"field_name": cmd.Args[3],
			"field_type": cmd.Args[4],
		})

		if err == nil {
			printCustomFields(s.Cfg, []database.CustomField{cf})
		}
		return err
	default:
		return errors.New("Invalid table provided for update.")
	}
}

func HandlerLoad(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("Delete command requires at least two arguments: operation, csv file name")
	}

	var method string

	switch cmd.Args[0] {
	case "insert":
		method = "POST"
	case "update":
		method = "PUT"
	case "delete":
		method = "DELETE"
	default:
		return fmt.Errorf("%s is not a valid operation.", cmd.Args[0])
	}

	records, err := parseCSV(cmd.Args[1])
	if err != nil {
		return err
	}

	for _, record := range records {
		item, err := requests.CallApi[map[string]string](s.Cfg, "items", method, record)
		if err != nil {
			return err
		}
		fmt.Printf("%sed %s: %s\n", cmd.Args[0], item["type"], item["name"])
	}

	return nil
}

func HandlerExport(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("Delete command requires at least two arguments: item type, csv file name")
	}

	endpoint := fmt.Sprintf("items?type=%s", cmd.Args[0])

	records, err := requests.CallApi[[]map[string]string](s.Cfg, endpoint, "GET", "")
	if err != nil {
		return err
	}

	return createCSV(s.Cfg, cmd.Args[1], records)
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

	return err
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

func printRolls(_ *config.CliConfig, rolls []rolls.RollTotalResult) {
	for i := range rolls {
		fmt.Printf("Roll %d - %s:\n", i+1, rolls[i].RollString)
		printRoll(rolls[i])
	}
}

func listRecords[T any](cfg *config.CliConfig, endpoint string, ids []string, printRecords func(*config.CliConfig, []T)) error {
	records, err := getOneOrMore[T](cfg, endpoint, "GET", ids)
	if err != nil {
		return err
	}

	printRecords(cfg, records)

	return nil
}

func getOneOrMore[T any](cfg *config.CliConfig, endpoint, method string, ids []string) ([]T, error) {
	var records []T
	var err error

	if len(ids) == 0 {
		records, err = requests.CallApi[[]T](cfg, endpoint, method, "")
		if err != nil {
			return records, err
		}
	} else {
		for i := range ids {
			endpointId := fmt.Sprintf("%s/%s", endpoint, ids[i])
			r, err := requests.CallApi[T](cfg, endpointId, "GET", "")
			if err != nil {
				return []T{}, err
			}
			records = append(records, r)
		}
	}

	return records, nil
}

func printTypes(_ *config.CliConfig, types []database.Type) {
	for i := range types {
		fmt.Printf("%d. ID: %s\n   Name: %s\n", i+1, types[i].ID, types[i].TypeName)
	}
}

func printCustomFields(cfg *config.CliConfig, fields []database.CustomField) {
	typeMap, _ := getTypeMap(cfg)

	for i := range fields {
		fmt.Printf("%d. ID: %s\n   Name: %s\n   Item Type: %s\n   Custom Field Type: %s\n",
			i+1,
			fields[i].ID,
			fields[i].CustomFieldName,
			typeMap[fields[i].TypeID],
			fields[i].CustomFieldType,
		)
	}
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

func parseCSV(filename string) ([]map[string]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return []map[string]string{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	rows, err := reader.ReadAll()
	if err != nil {
		return []map[string]string{}, err
	}

	records := make([]map[string]string, len(rows)-1)
	var headerRow []string
	for i, row := range rows {
		if i == 0 {
			headerRow = row
		} else {
			rowMap := make(map[string]string, len(row))
			for j, v := range row {
				rowMap[headerRow[j]] = v
			}
			records[i-1] = rowMap
		}
	}

	return records, nil
}

func createCSV(cfg *config.CliConfig, filename string, records []map[string]string) error {
	if len(records) == 0 {
		return errors.New("No records to export.")
	}
	csvFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	headerMap := make(map[string]int, len(records[0]))
	typeMap, err := getTypeMap(cfg)
	if err != nil {
		return err
	}

	for i, record := range records {
		if i == 0 {
			j := 0
			outRow := make([]string, len(record))
			for key := range record {
				headerMap[key] = j
				outRow[j] = key
				j += 1
			}
			err = writer.Write(outRow)
			if err != nil {
				return err
			}
		}

		outRow := make([]string, len(record))
		for key, val := range record {
			if key == "type" {
				typeId, err := uuid.Parse(val)
				if err != nil {
					return err
				}
				outRow[headerMap[key]] = typeMap[typeId]
			} else {
				outRow[headerMap[key]] = val
			}
		}
		err = writer.Write(outRow)
		if err != nil {
			return err
		}
	}

	return nil
}
