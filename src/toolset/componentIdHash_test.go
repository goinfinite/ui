package uiToolset

import "testing"

func TestHashComponentIdParts(t *testing.T) {
	firstHash := HashComponentIdParts("/records", "name")
	repeatedHash := HashComponentIdParts("/records", "name")
	if repeatedHash != firstHash {
		t.Errorf("HashNotStable: %d != %d", repeatedHash, firstHash)
	}

	otherHash := HashComponentIdParts("/records", "status")
	if otherHash == firstHash {
		t.Errorf("HashCollidedForDifferentParts: %d", otherHash)
	}

	if HashComponentIdParts() == firstHash {
		t.Errorf("HashCollidedForEmptyParts: %d", firstHash)
	}
}
