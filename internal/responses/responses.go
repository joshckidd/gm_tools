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

func (cfg *ApiConfig) ApiLogin(handler func(http.ResponseWriter, *http.Request, string, *ApiConfig)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		user, err := auth.ValidateJWT(tok, cfg.TokenSecret)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		handler(w, r, user, cfg)
	}
}

func PostRoll(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Roll string `json:"roll"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
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

func PostType(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Type string `json:"type"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
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

func GetTypes(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemTypes, err := cfg.DB.GetTypes(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, itemTypes)
}

func PostCustomField(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Type      string `json:"type"`
		FieldName string `json:"field_name"`
		FieldType string `json:"field_type"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
		return
	}

	if inParams.FieldType != "roll" && inParams.FieldType != "text" {
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

func GetCustomFields(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	itemTypes, err := cfg.DB.GetCustomFields(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, itemTypes)
}

func PostItem(w http.ResponseWriter, r *http.Request, user string, cfg *ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	inParams := map[string]string{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 400, "Invalid request")
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
	}

	for k, v := range inParams {
		if k != "type" && k != "name" && k != "description" {
			_, err := cfg.DB.CreateCustomFieldValue(r.Context(), database.CreateCustomFieldValueParams{
				CustomFieldValue: v,
				CustomFieldID:    customFieldIds[k],
				TypeID:           itemType.ID,
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

func (cfg *ApiConfig) PostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	inParams := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 500, "Invalid request")
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
		respondWithError(w, 500, "Invalid request")
		return
	}

	user, err := cfg.DB.GetUserWithUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	tok, err := auth.MakeJWT(user.Username, cfg.TokenSecret, time.Hour)
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
