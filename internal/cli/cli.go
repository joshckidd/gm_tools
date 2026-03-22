// this internal package holds all of the data structures and handler finctions that power the cli

package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

// handler for the roll command
// expects a single argument that is a roll string
// prints the result of a provided roll string
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

// handler for the list command
// expects an argument for the type of record to print and an optional number of additional ids for records
// types of records include rolls, types, custom_fields, items, or instances
// if no ids are supplied, print all records
// if ids are provided print just the records that correspond to those ids
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

// handler for the delete command
// expects an argument for the type of record to delete and one or more ids indicating the record(s) to delete
// types of records include types, custom_fields, items
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

// handler for the create command
// expects a first argument that indicates creating either types or custom_fields
// if a type is being created, a second argument is needed indicating the type name
// if a custom_field is being created, three additional arguments are needed
// indicating the item type, the custom field name, and the custom field type
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

// handler for the update command
// expects a first argument that indicates updating either types or custom_fields
// expects a second argument that is the id of the record to be updated
// if a type is being udated, a third argument is needed indicating the type name
// if a custom_field is being updated, three additional arguments are needed
// indicating the item type, the custom field name, and the custom field type
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

// handler for the load command
// expects a first command indicating insert, update, or delete
// expects a second argument indicating a csv file to be used as source data
// uses the csv file to insert, update, or delete items
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
		var endpoint string
		if cmd.Args[0] == "insert" {
			endpoint = "items"
		} else {
			endpoint = fmt.Sprintf("items/%s", record["id"])
			delete(record, "id")
		}
		if cmd.Args[0] == "delete" {
			id, err := requests.CallApi[string](s.Cfg, endpoint, method, record)
			if err != nil {
				return err
			}
			fmt.Printf("%sed %s\n", cmd.Args[0], id)
		} else {
			item, err := requests.CallApi[map[string]string](s.Cfg, endpoint, method, record)
			if err != nil {
				return err
			}
			fmt.Printf("%sed %s: %s\n", cmd.Args[0], item["type"], item["name"])
		}
	}

	return nil
}

// handler for the export command
// expects a first command indicating the item type to export
// expects a second argument indicating the name of the csv file to be written
// exports a csv file of all of the items of the provided type
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

// handler for the register command
// expects the username and password as arguments
func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("Register command requires two arguments: username, password")
	}

	user, err := requests.CallApi[database.User](
		s.Cfg,
		"users",
		"POST",
		map[string]string{
			"username": cmd.Args[0],
			"password": cmd.Args[1],
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("Created user %s\n", user.Username)

	return nil
}

// handler for the login command
// expects the username and password as arguments
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("Login command requires two arguments: username, password")
	}

	err := requests.LoginUser(s.Cfg, cmd.Args[0], cmd.Args[1])
	if err != nil {
		return err
	}

	err = s.Cfg.SetToken()

	if err == nil {
		fmt.Printf("%s logged in.\n", cmd.Args[0])
	}

	return err
}

// handler for the generate command
// expects a first argument that is an item type
// expects a second arguent that is a roll string
// creates instances of items of the type specified in the number of the executed roll string
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

// sets the api URL
func HandlerConnect(s *State, cmd Command) error {
	s.Cfg.APIUrl = cmd.Args[0]
	newAPIUrl, err := url.JoinPath(s.Cfg.APIUrl, "gm_tools")
	if err != nil {
		return err
	}

	_, err = http.Get(newAPIUrl)
	if err != nil {
		return err
	}

	s.Cfg.SetToken()
	fmt.Printf("Connected to %s\n", s.Cfg.APIUrl)
	return nil
}

// used to print a roll
func printRoll(tot rolls.RollTotalResult) {
	fmt.Printf("total: %d\n", tot.TotalResult)

	for i, rs := range tot.IndividualResults {
		fmt.Printf(" - Roll set %d: %d\n", i, rs.Result)

		for j, r := range rs.IndividualRolls {
			fmt.Printf(" --- Roll %d: %d\n", j, r)
		}
	}
}

// used to print multiple rolls
func printRolls(_ *config.CliConfig, rolls []rolls.RollTotalResult) {
	for i := range rolls {
		fmt.Printf("Roll %d - %s:\n", i+1, rolls[i].RollString)
		printRoll(rolls[i])
	}
}

// used to print one or more records of any type
func listRecords[T any](cfg *config.CliConfig, endpoint string, ids []string, printRecords func(*config.CliConfig, []T)) error {
	records, err := getOneOrMore[T](cfg, endpoint, "GET", ids)
	if err != nil {
		return err
	}

	printRecords(cfg, records)

	return nil
}

// used to get one or more records of any type from the database
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

// used to print type records
func printTypes(_ *config.CliConfig, types []database.Type) {
	for i := range types {
		fmt.Printf("%d. ID: %s\n   Name: %s\n", i+1, types[i].ID, types[i].TypeName)
	}
}

// used to print custom_field records
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

// used to print one or more item/instance records
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

// create a map of type UUIDs to their names
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

// parse the csv into data structures that can be marshalled into json for the api
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

// use a slice of items to create a csv file
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
			headRow := make([]string, len(record))
			for key := range record {
				headerMap[key] = j
				headRow[j] = key
				j += 1
			}
			err = writer.Write(headRow)
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

	writer.Flush()

	return nil
}

// handler for the help command
func HandlerHelp(_ *State, _ Command) error {
	fmt.Println(`Usage:
  gm_tools_cli [command] [arguments]

Available Commands:
  connect	Sets the API URL for GM Tools
			Expects one argument that is the API URL
  create	Creates a new record    
			Expects a first argument that indicates creating either types or custom_fields
			If a type is being created, a second argument is needed indicating the type name
			If a custom_field is being created, three additional arguments are needed
				indicating the item type, the custom field name, and the custom field type
  delete    Deletes a record by id
  			Expects an argument for the type of record to delete and one or more ids indicating the record(s) to delete
			Types of records include types, custom_fields, items
  export	Exports items of specified type as a csv file    
  			Expects a first command indicating the item type to export
			Expects a second argument indicating the name of the csv file to be written
  generate	Creates instances of items of the type specified in the number of the executed roll string
  			Expects a first argument that is an item type
			Expects a second arguent that is a roll string
  help      List this help information
  list      Lists one or more records of the specified type  
  			Expects an argument for the type of record to print and an optional number of additional ids for records
			Types of records include rolls, types, custom_fields, items, or instances
			If no ids are supplied, print all records
			If ids are provided print just the records that correspond to those ids
  load		Loads inserts, updates, or deletes of items from a csv file      
			Expects a first command indicating insert, update, or delete
			Expects a second argument indicating a csv file to be used as source data
  login     Logs a user in
			Expects the username and password as arguments
  register	Creates a new user
			Expects the username and password as arguments
  roll		Prints the result of a provided roll string
        	Expects a single argument that is a roll string
  update	Updates an existing record by id
			Expects a first argument that indicates updating either types or custom_fields
			Expects a second argument that is the id of the record to be updated
			If a type is being udated, a third argument is needed indicating the type name
			If a custom_field is being updated, three additional arguments are needed
				indicating the item type, the custom field name, and the custom field type`)

	return nil
}
