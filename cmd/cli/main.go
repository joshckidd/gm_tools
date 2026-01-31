package main

import (
	"fmt"

	"github.com/joshckidd/gm_tools/internal/rolls"
)

func main() {
	tot := rolls.RollAll(rolls.ParseRoll("min4d6-max3d8e"))
	fmt.Printf("total: %d\n", tot.TotalResult)

	for i, rs := range tot.IndividualResults {
		fmt.Printf("Roll set %d: %d\n", i, rs.Result)

		for j, r := range rs.IndividualRolls {
			fmt.Printf("Rol1 %d: %d\n", j, r)
		}
	}
}
