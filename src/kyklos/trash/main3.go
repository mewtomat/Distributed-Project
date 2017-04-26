package main

import (
	"math/big"
	"fmt"
)

func main(){
	// Convert the ID to a bigint
	idInt := *big.NewInt(5)
	// idInt.SetBytes(id)

	// Get the offset
	two := big.NewInt(2)
	offset := big.Int{}
	offset.Exp(two, big.NewInt(3), nil)

	// Diff
	diff := big.Int{}
	diff.Sub(&idInt, &offset)

	// Get the ceiling
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(3)), nil)

	// Apply the mod
	idInt.Mod(&diff, &ceil)

	fmt.Println(idInt.Int64())

	// Add together
	// return KeySpace{data:idInt.Bytes()}
}