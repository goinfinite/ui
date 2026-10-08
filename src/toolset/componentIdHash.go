package uiToolset

import "hash/fnv"

func HashComponentIdParts(componentIdParts ...string) uint64 {
	hasher := fnv.New64a()
	partSeparator := []byte{0}
	for _, componentIdPart := range componentIdParts {
		hasher.Write([]byte(componentIdPart))
		hasher.Write(partSeparator)
	}
	return hasher.Sum64()
}
