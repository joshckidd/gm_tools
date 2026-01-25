package main

import (
	"fmt"

	"github.com/joshckidd/gm_tools/internal/rolls"
)

func main() {
	tot, rolls := rolls.RollAll(rolls.ParseRoll("4d6+3d8"))
	fmt.Printf("total: %d\n", tot)

	for i, rs := range rolls {
		fmt.Printf("Roll set %d:\n", i)

		for j, r := range rs {
			fmt.Printf("Rol1 %d: %d\n", j, r)
		}
	}
}
