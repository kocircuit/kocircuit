package hash

import (
	"testing"
)

func TestHash(t *testing.T) {
	if MixScalar(1) == MixScalar(-2) {
		t.Errorf("MixScalar collision at %q = %q", MixScalar(1), MixScalar(-2))
	}
}
