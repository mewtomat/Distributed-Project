package kyklos

import (
	"bytes"
	"math/big"
)

func between(id1_space, id2_space, key_space KeySpace) bool {
	id1 := id1_space.Data
	id2 := id2_space.Data
	key := key_space.Data
	// Check for ring wrap around
	if bytes.Compare(id1, id2) == 1 {
		return bytes.Compare(id1, key) == -1 ||
			bytes.Compare(id2, key) == 1
	}

	if bytes.Compare(id1, id2) == 0 {
		return true
	}

	// Handle the normal case
	return bytes.Compare(id1, key) == -1 &&
		bytes.Compare(id2, key) == 1
}

// Checks if a key is between two ID's, right inclusive
func betweenRightIncl(id1_space, id2_space, key_space KeySpace) bool {
	id1 := id1_space.Data
	id2 := id2_space.Data
	key := key_space.Data
	// Check for ring wrap around
	if bytes.Compare(id1, id2) == 1 {
		return bytes.Compare(id1, key) == -1 ||
			bytes.Compare(id2, key) >= 0
	}

	if bytes.Compare(id1, id2) == 0 {
		return true
	}

	return bytes.Compare(id1, key) == -1 &&
		bytes.Compare(id2, key) >= 0
}

// Checks if a key is between two ID's, left inclusive
func betweenLeftIncl(id1_space, id2_space, key_space KeySpace) bool {
	id1 := id1_space.Data
	id2 := id2_space.Data
	key := key_space.Data
	// Check for ring wrap around
	if bytes.Compare(id1, id2) == 1 {
		return bytes.Compare(id1, key) <=0 ||
			bytes.Compare(id2, key) == 1
	}

	if bytes.Compare(id1, id2) == 0 {
		return true
	}

	return bytes.Compare(id1, key) <=0 &&
		bytes.Compare(id2, key) == 1
}

// Computes the offset by (n + 2^exp) % (2^mod)
func powerOffset(id_space KeySpace, exp int, mod int) KeySpace {
	id:=id_space.Data
	// Copy the existing slice
	off := make([]byte, len(id))
	copy(off, id)

	// Convert the ID to a bigint
	idInt := big.Int{}
	idInt.SetBytes(id)

	// Get the offset
	two := big.NewInt(2)
	offset := big.Int{}
	offset.Exp(two, big.NewInt(int64(exp)), nil)

	// Sum
	sum := big.Int{}
	sum.Add(&idInt, &offset)

	// Get the ceiling
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(mod)), nil)

	// Apply the mod
	idInt.Mod(&sum, &ceil)

	// Add together
	return KeySpace{Data:idInt.Bytes()}
}

func negativePowerOffset(id_space KeySpace, exp int, mod int) KeySpace {
	id:=id_space.Data
	// Copy the existing slice
	off := make([]byte, len(id))
	copy(off, id)

	// Convert the ID to a bigint
	idInt := big.Int{}
	idInt.SetBytes(id)

	// Get the offset
	two := big.NewInt(2)
	offset := big.Int{}
	offset.Exp(two, big.NewInt(int64(exp)), nil)

	// Diff
	diff := big.Int{}
	diff.Sub(&idInt, &offset)

	// Get the ceiling
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(mod)), nil)

	// Apply the mod
	idInt.Mod(&diff, &ceil)

	// Add together
	return KeySpace{Data:idInt.Bytes()}
}