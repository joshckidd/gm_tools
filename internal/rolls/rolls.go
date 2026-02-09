package rolls

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
)

type RollAggregate int

const (
	Sum RollAggregate = iota
	Max
	Min
)

type RollSignum int

const (
	Positive RollSignum = iota
	Negative
)

type RollType struct {
	Number    int           `json:"number"`
	Dice      int           `json:"dice"`
	Aggregate RollAggregate `json:"aggregate"`
	Signum    RollSignum    `json:"signum"`
	Exploding bool          `json:"exploding"`
}

type RollResult struct {
	Type            RollType `json:"type"`
	RollString      string   `json:"roll_string"`
	Result          int      `json:"result"`
	IndividualRolls []int    `json:"individual_rolls"`
}

type RollTotalResult struct {
	TotalResult       int          `json:"total_result"`
	IndividualResults []RollResult `json:"individual_results"`
}

func (r RollType) Roll() RollResult {
	var finalValue int
	individualRolls := make([]int, r.Number)

	for i := range r.Number {
		if r.Exploding {
			for individualRolls[i] = 0; individualRolls[i]%r.Dice == 0; {
				individualRolls[i] += rand.IntN(r.Dice) + 1
			}
		} else {
			individualRolls[i] = rand.IntN(r.Dice) + 1
		}
	}

	switch r.Aggregate {
	case Min:
		finalValue = r.Dice
	default:
		finalValue = 0
	}

	for _, ir := range individualRolls {
		switch r.Aggregate {
		case Sum:
			finalValue += ir
		case Min:
			if ir < finalValue {
				finalValue = ir
			}
		case Max:
			if ir > finalValue {
				finalValue = ir
			}
		}
	}

	if r.Signum == Negative {
		finalValue *= -1
	}

	return RollResult{
		Type:            r,
		Result:          finalValue,
		IndividualRolls: individualRolls,
		RollString:      r.GetRollString(),
	}
}

func (r RollType) GetRollString() string {
	var s, a, e string
	switch r.Signum {
	case Negative:
		s = "-"
	default:
		s = ""
	}
	switch r.Aggregate {
	case Min:
		a = "min"
	case Max:
		a = "max"
	default:
		a = ""
	}
	switch r.Exploding {
	case true:
		e = "e"
	default:
		e = ""
	}
	return fmt.Sprintf("%s%s%dd%d%s", s, a, r.Number, r.Dice, e)
}

func RollAll(rs []RollType) RollTotalResult {
	allRolls := make([]RollResult, len(rs))
	res := 0

	for i, r := range rs {
		rolls := r.Roll()
		res += rolls.Result
		allRolls[i] = rolls
	}

	return RollTotalResult{
		TotalResult:       res,
		IndividualResults: allRolls,
	}
}

func ParseRoll(str string) []RollType {
	r := regexp.MustCompile("([+-])?(min|max|sum)?([0-9]+)d?([0-9]+)?(e)?")

	parsedRolls := r.FindAllStringSubmatch(str, -1)

	res := make([]RollType, len(parsedRolls))

	for i, pr := range parsedRolls {
		if pr[1] == "-" {
			res[i].Signum = Negative
		} else {
			res[i].Signum = Positive
		}
		switch pr[2] {
		case "min":
			res[i].Aggregate = Min
		case "max":
			res[i].Aggregate = Max
		default:
			res[i].Aggregate = Sum
		}
		res[i].Number, _ = strconv.Atoi(pr[3])
		if pr[4] != "" {
			res[i].Dice, _ = strconv.Atoi(pr[4])
		} else {
			res[i].Dice = 1
		}
		if pr[5] == "e" {
			res[i].Exploding = true
		} else {
			res[i].Exploding = false
		}
	}
	return res
}
