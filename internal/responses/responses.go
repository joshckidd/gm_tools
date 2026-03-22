// internal package that provides the functions that generate api responses

package responses

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joshckidd/gm_tools/internal/auth"
	"github.com/joshckidd/gm_tools/internal/database"
	"github.com/joshckidd/gm_tools/internal/rolls"

	"net/http"
)

type ApiConfig struct {
	DB          *database.Queries
	TokenSecret string
}

// a wrapper for all response functions that handles logging the user in before executing the rest of the function
func (cfg *ApiConfig) ApiLogin(handler func(http.ResponseWriter, *http.Request, string, *ApiConfig)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}

		user, err := auth.ValidateJWT(tok, cfg.TokenSecret)
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}

		handler(w, r, user, cfg)
	}
}

// handles generating an individual roll
func PostRoll(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Roll string `json:"roll"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	tot := rolls.RollAll(rolls.ParseRoll(inParams.Roll))
	tot.RollString = inParams.Roll

	aggRoll, err := cfg.DB.CreateAggregateRoll(r.Context(), database.CreateAggregateRollParams{
		Result:   int32(tot.TotalResult),
		String:   tot.RollString,
		Username: user,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	for i := range tot.IndividualResults {
		_, err = cfg.DB.CreateRoll(r.Context(), database.CreateRollParams{
			String:          tot.IndividualResults[i].RollString,
			Result:          int32(tot.IndividualResults[i].Result),
			AggregateRollID: aggRoll.ID,
			Username:        user,
			IndividualRolls: fmt.Sprint(tot.IndividualResults[i].IndividualRolls),
		})
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	}

	respondWithJSON(w, 200, tot)
}

// handles getting all recent rolls
func GetRolls(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	userAggregateRolls, err := cfg.DB.GetAggregateRolls(r.Context(), user)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	res := make([]rolls.RollTotalResult, len(userAggregateRolls))

	for i := range userAggregateRolls {
		userRolls, err := cfg.DB.GetRolls(r.Context(), userAggregateRolls[i].ID)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		resParts := make([]rolls.RollResult, len(userRolls))

		for j := range userRolls {
			resParts[j].Type = rolls.ParseRoll(userRolls[j].String)[0]
			resParts[j].RollString = userRolls[j].String
			resParts[j].Result = int(userRolls[j].Result)
			ss := strings.Split(strings.Trim(userRolls[j].IndividualRolls, "[]"), " ")
			is := make([]int, len(ss))
			for k, s := range ss {
				is[k], _ = strconv.Atoi(s)
			}
			resParts[j].IndividualRolls = is
		}

		res[i].TotalResult = int(userAggregateRolls[i].Result)
		res[i].RollString = userAggregateRolls[i].String
		res[i].IndividualResults = resParts
	}

	respondWithJSON(w, 200, res)
}

// handles creating a new type
func PostType(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Type string `json:"type"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	itemType, err := cfg.DB.CreateType(r.Context(), database.CreateTypeParams{
		TypeName: inParams.Type,
		Username: user,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, itemType)
}

// handles getting all types
func GetTypes(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemTypes, err := cfg.DB.GetTypes(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, itemTypes)
}

// handles creating a new custom_field
func PostCustomField(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Type      string `json:"type"`
		FieldName string `json:"field_name"`
		FieldType string `json:"field_type"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	if inParams.FieldType != "roll" && inParams.FieldType != "picklist" {
		respondWithError(w, 422, "Bad value passed for field_type. Expecting 'roll' or 'text'.")
		return
	}

	itemType, err := cfg.DB.GetTypeByName(r.Context(), inParams.Type)
	if err != nil {
		respondWithError(w, 422, "Bad value passed for 'type'")
		return
	}

	customField, err := cfg.DB.CreateCustomFields(r.Context(), database.CreateCustomFieldsParams{
		TypeID:          itemType.ID,
		Username:        user,
		CustomFieldName: inParams.FieldName,
		CustomFieldType: inParams.FieldType,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, customField)
}

// handles getting all custom_fields
func GetCustomFields(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemTypes, err := cfg.DB.GetCustomFields(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, itemTypes)
}

// handles creating a new item
func PostItem(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := map[string]string{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	itemType, err := cfg.DB.GetTypeByName(r.Context(), inParams["type"])
	if err != nil {
		respondWithError(w, 422, "Bad value passed for 'type'")
		return
	}

	customFields, err := cfg.DB.GetCustomFieldForType(r.Context(), itemType.ID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if inParams["name"] == "" || inParams["description"] == "" {
		respondWithError(w, 422, "Item must have a name and description.")
		return

	}

	customFieldIds := map[string]uuid.UUID{}

	for k := range inParams {
		if k != "type" && k != "name" && k != "description" {
			found := false
			for i := range customFields {
				//to do: validate custom field types
				if k == customFields[i].CustomFieldName {
					found = true
					customFieldIds[k] = customFields[i].ID
				}
			}
			if !found {
				respondWithError(w, 422, fmt.Sprintf("No field %s for item type %s.", k, itemType.TypeName))
				return
			}
		}
	}

	item, err := cfg.DB.CreateItem(r.Context(), database.CreateItemParams{
		TypeID:          itemType.ID,
		Username:        user,
		ItemName:        inParams["name"],
		ItemDescription: inParams["description"],
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	itemMap := map[string]string{
		"name":        item.ItemName,
		"description": item.ItemDescription,
		"type":        item.TypeID.String(),
		"created_at":  item.CreatedAt.String(),
		"updated_at":  item.UpdatedAt.String(),
		"username":    item.Username,
	}

	for k, v := range inParams {
		if k != "type" && k != "name" && k != "description" {
			_, err := cfg.DB.CreateCustomFieldValue(r.Context(), database.CreateCustomFieldValueParams{
				CustomFieldValue: v,
				CustomFieldID:    customFieldIds[k],
				ItemID:           item.ID,
				Username:         user,
			})
			if err != nil {
				respondWithError(w, 500, err.Error())
				return
			}
			itemMap[k] = v
		}
	}

	respondWithJSON(w, 200, itemMap)
}

// handles listing all items
func GetItems(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	t := r.URL.Query().Get("type")

	var baseItems []database.Item
	var err error

	if t == "" {
		baseItems, err = cfg.DB.GetItems(r.Context())
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	} else {
		itemType, err := cfg.DB.GetTypeByName(r.Context(), t)
		if err != nil {
			respondWithError(w, 422, "Bad value passed for 'type'")
			return
		}
		baseItems, err = cfg.DB.GetItemIdsByType(r.Context(), itemType.ID)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	}

	items := make([]map[string]string, len(baseItems))

	for i := range baseItems {
		item, err := fillOutItemFields(baseItems[i], r, cfg)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		items[i] = item
	}

	respondWithJSON(w, 200, items)
}

// handles creating a new user
func (cfg *ApiConfig) PostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(inParams.Password)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       inParams.Username,
		HashedPassword: hashedPassword,
	})
	if err.Error() == "pq: duplicate key value violates unique constraint \"users_pkey\"" {
		respondWithError(w, 409, fmt.Sprintf("%s is already in use as a username. Please select another.", inParams.Username))
		return
	} else if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 201, database.User{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
	})
}

// handles logging in a user
func (cfg *ApiConfig) UserLogin(w http.ResponseWriter, r *http.Request) {
	type loginParam struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	type returnUserRow struct {
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Token     string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := loginParam{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	user, err := cfg.DB.GetUserWithUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	// set JWT expiration here
	tok, err := auth.MakeJWT(user.Username, cfg.TokenSecret, 4*time.Hour)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	val, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if val == true {
		userResp := returnUserRow{
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     tok,
		}
		respondWithJSON(w, 200, userResp)
		return
	}
	respondWithError(w, 401, "Incorrect email or password")
}

// handles deleting a type by id
func DeleteType(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	typeId, err := uuid.Parse(r.PathValue("typeId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	err = cfg.DB.DeleteType(r.Context(), typeId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	respondWithJSON(w, 200, typeId)
}

// handles deleting a custom_field by id
func DeleteCustomField(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	fieldId, err := uuid.Parse(r.PathValue("fieldId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	err = cfg.DB.DeleteCustomField(r.Context(), fieldId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	respondWithJSON(w, 200, fieldId)
}

// handles deleting an item by id
func DeleteItem(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemId, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	err = cfg.DB.DeleteItem(r.Context(), itemId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	respondWithJSON(w, 200, itemId)
}

// handles updating a type by id
func PutType(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	typeId, err := uuid.Parse(r.PathValue("typeId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Type string `json:"type"`
	}{}

	err = decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
		return
	}

	itemType, err := cfg.DB.UpdateType(r.Context(), database.UpdateTypeParams{
		TypeName: inParams.Type,
		ID:       typeId,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, itemType)
}

// handles getting a single type by id
func GetType(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	typeId, err := uuid.Parse(r.PathValue("typeId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	itemType, err := cfg.DB.GetTypeById(r.Context(), typeId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	respondWithJSON(w, 200, itemType)
}

// handles getting a single custom_field by id
func GetCustomField(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	fieldId, err := uuid.Parse(r.PathValue("fieldId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	customField, err := cfg.DB.GetCustomFieldById(r.Context(), fieldId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	respondWithJSON(w, 200, customField)
}

// handles getting a single item by id
func GetItem(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemId, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	baseItem, err := cfg.DB.GetItemById(r.Context(), itemId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	item, err := fillOutItemFields(baseItem, r, cfg)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, item)
}

// handles updating a custom_field by id
func PutCustomField(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	fieldId, err := uuid.Parse(r.PathValue("fieldId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	customField, err := cfg.DB.GetCustomFieldById(r.Context(), fieldId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Type      string `json:"type"`
		FieldName string `json:"field_name"`
		FieldType string `json:"field_type"`
	}{}

	err = decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
		return
	}

	if inParams.FieldType != "" && (inParams.FieldType == "picklist" || inParams.FieldType == "roll") {
		customField.CustomFieldType = inParams.FieldType
	} else {
		respondWithError(w, 422, "Bad value passed for field_type. Expecting 'roll' or 'text'.")
		return
	}

	if inParams.Type != "" {
		itemType, err := cfg.DB.GetTypeByName(r.Context(), inParams.Type)
		if err != nil {
			respondWithError(w, 422, "Bad value passed for 'type'")
			return
		}
		customField.TypeID = itemType.ID
	}

	if inParams.FieldName != "" {
		customField.CustomFieldName = inParams.FieldName
	}

	updatedCustomField, err := cfg.DB.UpdateCustomField(r.Context(), database.UpdateCustomFieldParams{
		TypeID:          customField.TypeID,
		CustomFieldName: customField.CustomFieldName,
		CustomFieldType: customField.CustomFieldType,
		ID:              fieldId,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, updatedCustomField)
}

// handles updating an item by id
func PutItem(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemId, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	baseItem, err := cfg.DB.GetItemById(r.Context(), itemId)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	item, err := fillOutItemFields(baseItem, r, cfg)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	inParams := map[string]string{}

	err = decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	var itemType database.Type

	if inParams["type"] != "" {
		itemType, err = cfg.DB.GetTypeByName(r.Context(), inParams["type"])
		if err != nil {
			respondWithError(w, 422, "Bad value passed for 'type'")
			return
		}
		item["type"] = itemType.ID.String()
	} else {
		typeId, err := uuid.Parse(item["type"])
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		itemType, err = cfg.DB.GetTypeById(r.Context(), typeId)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	}

	if inParams["name"] != "" {
		item["name"] = inParams["name"]
	}

	if inParams["description"] != "" {
		item["description"] = inParams["description"]
	}

	customFields, err := cfg.DB.GetCustomFieldForType(r.Context(), itemType.ID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	customFieldIds := map[string]uuid.UUID{}

	for k := range inParams {
		if k != "type" && k != "name" && k != "description" {
			found := false
			for i := range customFields {
				//to do: validate custom field types
				if k == customFields[i].CustomFieldName {
					found = true
					customFieldIds[k] = customFields[i].ID
				}
			}
			if !found {
				respondWithError(w, 422, fmt.Sprintf("No field %s for item type %s.", k, inParams["type"]))
				return
			}
		}
	}

	customFieldValues, err := cfg.DB.GetCustomFieldValues(r.Context(), itemId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	for k := range inParams {
		if k != "type" && k != "name" && k != "description" && inParams[k] != "" {
			item[k] = inParams[k]
			found := false
			for i := range customFieldValues {
				//to do: validate custom field types
				if k == customFieldValues[i].CustomFieldName {
					found = true
					_, err := cfg.DB.UpdateCustomFieldValue(r.Context(), database.UpdateCustomFieldValueParams{
						ID:               customFieldValues[i].ID,
						CustomFieldValue: inParams[k],
					})
					if err != nil {
						respondWithError(w, 500, err.Error())
						return
					}
				}
			}
			if !found {
				_, err := cfg.DB.CreateCustomFieldValue(r.Context(), database.CreateCustomFieldValueParams{
					CustomFieldID:    customFieldIds[k],
					CustomFieldValue: inParams[k],
					ItemID:           itemId,
					Username:         user,
				})
				if err != nil {
					respondWithError(w, 500, err.Error())
					return
				}
			}
		}
	}

	_, err = cfg.DB.UpdateItem(r.Context(), database.UpdateItemParams{
		ID:              itemId,
		ItemName:        item["name"],
		ItemDescription: item["description"],
		TypeID:          itemType.ID,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, item)
}

// handles generating new instances of items
func PostInstances(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Number string `json:"number"`
		Type   string `json:"type"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	itemType, err := cfg.DB.GetTypeByName(r.Context(), inParams.Type)
	if err != nil {
		respondWithError(w, 422, "Bad value passed for 'type'")
		return
	}

	itemIds, err := cfg.DB.GetItemIdsByType(r.Context(), itemType.ID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	items := rolls.RandomFromSliceN(itemIds, rolls.RollAll(rolls.ParseRoll(inParams.Number)).TotalResult)

	instances := make([]map[string]string, len(items))

	for i := range items {
		instMap := map[string]string{
			"name":        items[i].ItemName,
			"description": items[i].ItemDescription,
		}

		inst, err := cfg.DB.CreateInstance(r.Context(), database.CreateInstanceParams{
			ItemID:   items[i].ID,
			Username: user,
		})
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		instMap["id"] = inst.ID.String()

		customFields, err := cfg.DB.GetCustomFieldValues(r.Context(), items[i].ID)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		for j := range customFields {
			switch customFields[j].CustomFieldType {
			case "roll":
				value, err := cfg.DB.CreateCustomFieldInstanceValue(r.Context(), database.CreateCustomFieldInstanceValueParams{
					CustomFieldValue: strconv.Itoa(rolls.RollAll(rolls.ParseRoll(customFields[j].CustomFieldValue)).TotalResult),
					Username:         user,
					CustomFieldID:    customFields[j].CustomFieldID,
					InstanceID:       inst.ID,
				})
				if err != nil {
					respondWithError(w, 500, err.Error())
					return
				}
				instMap[customFields[j].CustomFieldName] = value.CustomFieldValue
			case "picklist":
				value, err := cfg.DB.CreateCustomFieldInstanceValue(r.Context(), database.CreateCustomFieldInstanceValueParams{
					CustomFieldValue: rolls.RandomPicklistValue(customFields[j].CustomFieldValue),
					Username:         user,
					CustomFieldID:    customFields[j].CustomFieldID,
					InstanceID:       inst.ID,
				})
				if err != nil {
					respondWithError(w, 500, err.Error())
					return
				}
				instMap[customFields[j].CustomFieldName] = value.CustomFieldValue
			}
		}

		instances[i] = instMap
	}

	respondWithJSON(w, 200, instances)
}

// handles getting all recent instances
func GetInstances(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	baseInstances, err := cfg.DB.GetInstances(r.Context(), user)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	instances := make([]map[string]string, len(baseInstances))

	for i := range baseInstances {
		instance, err := fillOutInstanceFields(baseInstances[i], r, cfg)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		instances[i] = instance
	}

	respondWithJSON(w, 200, instances)
}

// Used for testing that the service is up
func (cfg *ApiConfig) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("GM Tools is up and running."))
}

// fill out custom field values for base items
func fillOutItemFields(baseItem database.Item, r *http.Request, cfg *ApiConfig) (map[string]string, error) {
	item := map[string]string{
		"id":          baseItem.ID.String(),
		"name":        baseItem.ItemName,
		"description": baseItem.ItemDescription,
		"type":        baseItem.TypeID.String(),
		"created_at":  baseItem.CreatedAt.String(),
		"updated_at":  baseItem.UpdatedAt.String(),
		"username":    baseItem.Username,
	}

	customFields, err := cfg.DB.GetCustomFieldValues(r.Context(), baseItem.ID)
	if err != nil {
		return map[string]string{}, err
	}

	for j := range customFields {
		item[customFields[j].CustomFieldName] = customFields[j].CustomFieldValue
	}

	return item, nil
}

// fill out custom field values for base instances
func fillOutInstanceFields(baseInstance database.GetInstancesRow, r *http.Request, cfg *ApiConfig) (map[string]string, error) {
	instance := map[string]string{
		"id":          baseInstance.ID.String(),
		"name":        baseInstance.ItemName,
		"description": baseInstance.ItemDescription,
		"type":        baseInstance.TypeID.String(),
		"created_at":  baseInstance.CreatedAt.String(),
		"updated_at":  baseInstance.UpdatedAt.String(),
		"username":    baseInstance.Username,
	}

	customFields, err := cfg.DB.GetCustomFieldInstanceValues(r.Context(), baseInstance.ID)
	if err != nil {
		return map[string]string{}, err
	}

	for j := range customFields {
		instance[customFields[j].CustomFieldName] = customFields[j].CustomFieldValue
	}

	return instance, nil
}

// provide a json http response
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	val, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(val)
}

// provide an error http response
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}

	respError := returnError{
		Error: msg,
	}

	dat, err := json.Marshal(respError)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
