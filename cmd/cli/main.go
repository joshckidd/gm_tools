package main

import (
	"fmt"

	"github.com/joshckidd/gm_tools/internal/rolls"
)

func main() {
	tot, rolls := rolls.RollAll(rolls.ParseRoll("min4d6-max3d8e"))
	fmt.Printf("total: %d\n", tot)

	for i, rs := range rolls {
		fmt.Printf("Roll set %d: %d\n", i, rs.Result)

		for j, r := range rs.IndividualRolls {
			fmt.Printf("Rol1 %d: %d\n", j, r)
		}
	}
}
