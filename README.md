# GM Tools

## Overview

This is a tool for Game Masters running Table Top Role Playing Games. It will generate rolls, store items of any kind (e.g. monsters, magic items) in a database, and select random items based on criteria provided.

I created this tool specifically to help me in running games using the [Cypher System](https://cypher-system.com) created by [Monte Cook Games](https://www.montecookgames.com). Most of the examples that I will use below reference the Cypher System, but it could be used with any system where you would want to generate random items.

There are two main parts of GM Tools in its current state:

- Rolls
- Items and Instances

### Rolls

For all rolls generated in GM tools, a _roll string_ is used that should look familiar to anyone who plays TTRPGs regularly. For example, 2d6 is a roll string that indicated that 2 6-sided dice should be rolled and their totals summed. I've added some additional items to standard roll strings to allow for more functionality. Each roll string has five parts.

- Number - The number of dice to be rolled. This is the only element of the roll string that is _required_. What this means is that 10 is a valid roll string all on its own. It just always has a value of 10.
- Die - The number of sides on the dice you are rolling. This is indicated with a 'd' after the number and before the die. 4d6, 2d8, and 1d20 are all roll strings with a number and die. Again this should all be familiar.
- Exploding - This will be a familiar concept to many TTRPG players. It indicates that when the highest value is rolled on a die, that die should be rolled again and that result should be added to the original result. If the highest value is rolled a second time, the die should be rolled a third time, etc. I note this by adding an 'e' after the die value. So 2d6e indicates rolling 2 exploding 6-sided dice. In the absence of 'e', dice are assumed _not_ to be exploding.
- Aggregate - The possibilities here are 'min', 'max', and 'sum' and they are added to the front of the roll string. This indicates how the resulting rolls should be aggregated. 'min' indicates that the minimum value should be selected. 'max' indicates that the maximum value should be selected. 'sum' indicates that the value of all of the rolls should be added together for the final result. This is the default if no aggregate is specified. max2d20 is how you would indicate rolling with advantage in D&D. min2d20 is how you would represent rolling with disadvantage.
- Signum - This indicates whether the final result should be positive or negative and is indicated by a '+' or '-' added to the front of the roll string. Positive is the default value. This is usefull because __roll strings can include multiple rolls__. For example, 1d6+1 is a valid roll string, as are: 1d8+4d6+4, min2d20-2, etc.

### Items and Instances

In GM Tools, an _item_ is anything that you might want to generate randomly. This could be magic items, monsters, etc. For my [Numenera](https://numenera.com) game, I needed a way to randomly generate cyphers, artifacts, and oddities.

GM Tools allows you to create any number of _types_ for items. I created the types cypher, artifact, and oddity. You might create types for potion, monster, etc.

To start, every item has a name and a description, but you may need to store additional information for a specific type. For example, my cypher type requires additional information for level (how strong is the cypher), form (what it looks like), and cypher type (does it take up one slot in your inventory or two.) My artifact type requires information for level, form, and depletion. And for oddity, just name and description are enough. To store this information, GM Tools uses _custom fields_. For each custom field, GM Tools requires an item type that it applies to, a name, and a custom field type.

There are two valid cutom field types, _roll_ and _picklist_. A roll custom field contains a roll string and indicates that, when an item is randomly generated, the value for that field is generated based on the roll string provided. For example, my level field on cypher is a roll and 1d6 might be the value stored in that field. What that means is that when that item is generated, a value from 1 to 6 is generated as the level. Rolls can also store static integer values. For example, if I want the level for a specific cypher to always be 10, I can just use 10 as the roll string in the level field. A picklist custom field contains a series of values separated by ';' and indicates that, when an item is randomly generated, one of the values will be randomly selected for that field. My form field on cypher is a picklist and 'pill;injectable' is a value that might be stored there. What that means is that when that item is generated, one of either 'pill' or 'injectable' is selected as the form. Picklists can also store string values. For example, if I want the form for a specific cypher to always be 'headband', I set the value to 'headband'. Since there is only one option, that option is always selected.

For my example above, I create 6 custom fields:

- Item Type: cypher, Name: level, Custom Field Type: roll
- Item Type: cypher, Name: form, Custom Field Type: picklist
- Item Type: cypher, Name: cypher_type, Custom Field Type: picklist
- Item Type: artifact, Name: level, Custom Field Type: roll
- Item Type: artifact, Name: form, Custom Field Type: picklist
- Item Type: articact, Name: depletion, Custom Field Type: picklist

Note that level and form need to be created separately for both cypher and artifact.

When I have my types and custom fields set up, I can start populating my database with items. Again, an item is anything that I might want to generate randomly for my game. To aid in populating items into the database, the GM Tools CLI allows you to load or export csv files. When you load a csv file, you can insert, update, or delete records. If you are inserting or updating, your csv file must have a header row indicating the columns in the file and those columns must include 'name', 'description', and 'type'. If you are updating, you must also have an 'id' column where the id of the record to be updated is stored. Additionally, you may have columns with the name of any custom fields for that item type. When you export a csv file, you export all values for all items of a specific type.

When you have items in your database, you can generate _instances_ of items. When you generate instances, you provide a type of item you would like to generate and a roll string indicating the number of instances you would like to generate. If the roll string is 1d6, 1 to 6 instances will be generated. If the roll string is 5, exactly 5 instances will be generated. Then that number of items will be randomly selected from the database using the given type and instances will be created from those items. Custom fields will be filled in with static values at that time. Rolls will become numbers and picklists will become simple strings.

All generated rolls and instances will be stored for a limited amount of time so that they can be referenced later.

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

- Create a true multiuser api. Right now, the api requires a user to log in, but it is intended to be used by a single user only. Addition of full multiuser functionality including better security, refresh tokens, user permissions, etc. is the next big feature for me and one that will enable all other roadmap features.
- Add an inventory system. Right now instances of items are only meant to last a short while. My plan is to create inventories for users where items can be stored long term until they are deleted. Each user would have their own inventory and GMs can move instances of items into their inventory or the inventory of one of their players.
- Add a player party system. This would allow players to be grouped together into a party. This could enable some additional group permissions. And it would also enable group inventories. So the party itself could have an inventory rather than just personal inventories.
- Add item tagging. It would be nice to generate items not just by type, but by other characteristics of the item. For example, a GM might want or generate fire related magic items, if they are found in a volcano. They might want to generate water related monsters, if the players are travelling on the ocean. A robust tagging system where one or more tags can be attached to an item could accomplish this.
- A web front end. Let's be real. My players are not going to use a command line tool. But they would use a web front end. This is the whole reason why I created a rest API rather than just a command line tool that talks to a postgres database.
- Item search. This would enable finding specific items in the database based on name, description, tags, etc. A great to have feature, but probably the hardest one in this list to implement well.
