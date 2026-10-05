package uiToolset

import "hash/fnv"

func HashComponentIdParts(componentIdParts ...string) uint64 {
	hasher := fnv.New64a()
	for _, componentIdPart := range componentIdParts {
		hasher.Write([]byte(componentIdPart))
	}
	return hasher.Sum64()
}
