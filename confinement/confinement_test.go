package confinement

import (
	"testing"

	"go.uber.org/goleak"
)

func TestAdhocConfinement(t *testing.T) {
	AdhocConfinement()
}

func TestLexicalConfinement(t *testing.T) {
	LexicalConfinement()
}

func TestLexicalConfinementII(t *testing.T) {
	LexicalConfinementII()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
