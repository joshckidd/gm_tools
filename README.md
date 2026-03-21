# GM Tools

## Overview

## How To Use

## API Documentation

## CLI Documentation

### Usage

- gm_tools_cli [command] [arguments]

### Available Commands

- create
  - Creates a new record
  - Expects a first argument that indicates creating either types or custom_fields
  - If a type is being created, a second argument is needed indicating the type name
  - If a custom_field is being created, three additional arguments are needed indicating the item type, the custom field name, and the custom field type
- delete
  - Deletes a record by id
  - Expects an argument for the type of record to delete and one or more ids indicating the record(s) to delete
  - Types of records include types, custom_fields, items
- export
  - Exports items of specified type as a csv file
  - Expects a first command indicating the item type to export
  - Expects a second argument indicating the name of the csv file to be written
- generate
  - Creates instances of items of the type specified in the number of the executed roll string
  - Expects a first argument that is an item type
  - Expects a second arguent that is a roll string
- help
  - List this help information
- list
  - Lists one or more records of the specified type  
  - Expects an argument for the type of record to print and an optional number of additional ids for records
  - Types of records include rolls, types, custom_fields, items, or instances
  - If no ids are supplied, print all records
  - If ids are provided print just the records that correspond to those ids
- load
  - Loads inserts, updates, or deletes of items from a csv file
  - Expects a first command indicating insert, update, or delete
  - Expects a second argument indicating a csv file to be used as source data
- login
  - Logs a user in
  - Expects the username and password as arguments
- roll
  - Prints the result of a provided roll string
  - Expects a single argument that is a roll string
- update
  - Updates an existing record by id
  - Expects a first argument that indicates updating either types or custom_fields
  - Expects a second argument that is the id of the record to be updated
  - If a type is being udated, a third argument is needed indicating the type name
  - If a custom_field is being updated, three additional arguments are needed indicating the item type, the custom field name, and the custom field type

## Roadmap
