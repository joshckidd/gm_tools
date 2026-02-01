package main

import (
	"fmt"
	"os"

	"github.com/joshckidd/gm_tools/internal/requests"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("I require an argument!")
		os.Exit(0)
	}

	tot, _ := requests.GenerateRoll(args[1])
	fmt.Printf("total: %d\n", tot.TotalResult)

	for i, rs := range tot.IndividualResults {
		fmt.Printf(" - Roll set %d: %d\n", i, rs.Result)

		for j, r := range rs.IndividualRolls {
			fmt.Printf(" --- Roll %d: %d\n", j, r)
		}
	}
}
